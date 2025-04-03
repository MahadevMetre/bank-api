package authorization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	"bitbucket.org/paydoh/paydoh-commons/jsonwebtoken"
	"bitbucket.org/paydoh/paydoh-commons/pkg/task"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"
)

type Authorization interface {
	GetAuthorizationToken(request *requests.AuthorizationRequest) (interface{}, error)
}

// AuthorizationStore @impl Authorization static check
type AuthorizationStore struct {
	db              *sql.DB
	context         context.Context
	LoggerService   *commonSrv.LoggerService
	Mongo           *database.Document
	Redis           *database.InMemory
	taskEnqueuer    task.TaskEnqueuer
	auditLogService services.AuditLogService
}

func NewAuthorizationStore(log *commonSrv.LoggerService, db *sql.DB, ctx context.Context, mongo *database.Document, redis *database.InMemory, taskEnqueuer task.TaskEnqueuer, auditLogService services.AuditLogService) *AuthorizationStore {
	return &AuthorizationStore{
		db:              db,
		context:         ctx,
		LoggerService:   log,
		Mongo:           mongo,
		Redis:           redis,
		taskEnqueuer:    taskEnqueuer,
		auditLogService: auditLogService,
	}
}

func (store *AuthorizationStore) GetAuthorizationToken(ctx context.Context, request *requests.AuthorizationRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:    constants.AUTHORIZATION,
		Message:   "GetAuthorizationToken log",
		UserID:    utils.GetUserIDFromContext(ctx),
		RequestID: utils.GetRequestIDFromContext(ctx),
	}

	existingUser, err := models.GetUserDataByMobileNumber(store.db, request.MobileNumber)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			logData.Message = "GetAuthorizationToken: User not found"
			store.LoggerService.LogError(logData)
			insertedRow := models.NewUserData()

			passphrase := security.GenerateRandomPassphrase()

			encryptedPassphrase, err := security.Encrypt([]byte(passphrase), []byte(constants.AesPassPhrase))

			if err != nil {
				logData.Message = "GetAuthorizationToken: Error while encrypting passphrase"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			userid, err := security.GenerateRandomUUID(15)
			if err != nil {
				logData.Message = "GetAuthorizationToken: Error while generating userid"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			rand, err := security.GenerateRandomUUID(6)
			if err != nil {
				logData.Message = "GetAuthorizationToken: Error while generating random string"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			applicantId := fmt.Sprintf("PAYDOH%s", rand)
			deviceId := utils.GetDeviceIdFromContext(ctx)

			if deviceId == "" {
				return nil, errors.New("device id should not be empty")
			}

			err = store.db.QueryRow(
				`
					INSERT into user_data (user_id, mobile_number, applicant_id, signing_key, device_id) VALUES($1, $2, $3, $4, $5)
					RETURNING id, user_id, mobile_number, signing_key, created_at, updated_at;
				`,
				userid,
				request.MobileNumber,
				applicantId,
				encryptedPassphrase,
				deviceId,
			).Scan(
				&insertedRow.Id,
				&insertedRow.UserId,
				&insertedRow.MobileNumber,
				&insertedRow.SigningKey,
				&insertedRow.CreatedAt,
				&insertedRow.UpdatedAt,
			)

			if err != nil {
				var pqErr *pq.Error
				if errors.As(err, &pqErr) {
					if pqErr.Code == "23505" && pqErr.Constraint == "user_data_device_id_key" {
						logData.Message = "GetAuthorizationToken: Duplicate key error: " + pqErr.Constraint
						store.LoggerService.LogError(logData)
						return nil, errors.New("this device is already registered with another mobile number")
					}
				}

				logData.Message = "GetAuthorizationToken: Error while inserting user data"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			//creation of user onboarding status data
			err = models.CreateUserOnboardingStatus(constants.AUTHORIZATION_STEP, userid)
			if err != nil {
				return nil, err
			}

			newUser := &models.User{
				UserId:       userid,
				MobileNumber: request.MobileNumber,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			store.InsertUserInMongoDB(newUser)
			store.InsertUserNotificationPreferencesInMongoDB(newUser)

			userIP := utils.GetSourceIPFromContext(ctx)
			if userIP == "" {
				return nil, errors.New("device ip should not be empty")
			}

			token, err := store.generateToken(userid, insertedRow.SigningKey, utils.GetUserAgentFromContext(ctx), userIP)
			if err != nil {
				logData.Message = "GetAuthorizationToken: Error while generating jwt token"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			response := map[string]interface{}{
				"token": token,
			}

			return response, nil
		}

		logData.Message = "GetAuthorizationToken: Error while fetching user data"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	userIP := utils.GetSourceIPFromContext(ctx)
	if userIP == "" {
		return nil, errors.New("device ip should not be empty")
	}

	token, err := store.generateToken(existingUser.UserId, existingUser.SigningKey, utils.GetUserAgentFromContext(ctx), userIP)
	if err != nil {
		logData.Message = "GetAuthorizationToken: Error while generating jwt token"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	response := map[string]interface{}{
		"token": token,
	}

	requestJson, err := request.ToJSON()
	if err != nil {
		logData.Message = "GetAuthorizationToken: Error while marshaling request to json"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	decoded, err := security.Decrypt(existingUser.SigningKey, []byte(constants.AesPassPhrase))
	if err != nil {
		return "", fmt.Errorf("error decrypting signing key: %v", err)
	}

	encrytedData, err := security.Encrypt(requestJson, []byte(decoded))
	if err != nil {
		logData.Message = "GetAuthorizationToken: Error while encrypting data" + err.Error()
		store.LoggerService.LogError(logData)
		return nil, err
	}

	// saving audit-log
	if err := store.auditLogService.Save(ctx, &services.AuditLog{
		UserID:         existingUser.UserId,
		ApplicantID:    existingUser.ApplicantId,
		RequestURL:     "/api/authorization",
		HTTPMethod:     "POST",
		RequestBody:    encrytedData,
		ResponseStatus: 200,
		Action:         constants.AUTH_ACTION,
	}); err != nil {
		logData.Message = "GetAuthorizationToken: Error while saving audit log"
		store.LoggerService.LogError(logData)
	}

	logData.Message = "GetAuthorizationToken: Token generated successfully"
	logData.EndTime = time.Now()
	store.LoggerService.LogInfo(logData)
	return response, nil
}

func (store *AuthorizationStore) generateToken(userId string, signingKey string, userAgent string, userIP string) (string, error) {
	existingToken, err := store.Redis.Get(fmt.Sprintf(constants.TokenKeyFormat, userId))
	if err == nil && existingToken != "" {
		return existingToken, nil
	}

	decoded, err := security.Decrypt(signingKey, []byte(constants.AesPassPhrase))
	if err != nil {
		return "", fmt.Errorf("error decrypting signing key: %v", err)
	}

	userId3 := userId[:3]
	signKey3 := decoded[:3]

	encryptingText := fmt.Sprintf("%s%s|%s%s", signKey3, userId[3:], userId3, decoded[3:])

	token, err := jsonwebtoken.GenerateJWT(encryptingText, userAgent, userIP, constants.JwtKey, constants.JwtExpTime, "paydoh-bank")
	if err != nil {
		return "", fmt.Errorf("error generating jwt token: %v", err)
	}

	err = store.Redis.Set(fmt.Sprintf(constants.TokenKeyFormat, userId), token, constants.JwtExpTime)
	if err != nil {
		return "", fmt.Errorf("error setting token in redis: %v", err)
	}

	return token, nil
}

func (store *AuthorizationStore) GetAuthorizationTokenByUserId(ctx context.Context, request *requests.AuthorizationRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.AUTHORIZATION,
		RequestURI: "/api/authorization",
		Message:    "GetAuthorizationToken By UserId log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingUser := models.NewUserData()

	err := store.db.QueryRow(
		`
			SELECT id, user_id, mobile_number, signing_key, created_at, updated_at
			FROM user_data
			WHERE user_id = $1
		`,
		request.UserId,
	).Scan(
		&existingUser.Id,
		&existingUser.UserId,
		&existingUser.MobileNumber,
		&existingUser.SigningKey,
		&existingUser.CreatedAt,
		&existingUser.UpdatedAt,
	)

	if err != nil {
		logData.Message = "GetAuthorizationTokenByUserId: Error while fetching user data"
		store.LoggerService.LogError(logData)
		return nil, errors.New("error while fetching user data")
	}

	decoded, err := security.Decrypt(existingUser.SigningKey, []byte(constants.AesPassPhrase))
	if err != nil {
		logData.Message = "GetAuthorizationTokenByUserId: Error while decrypting signing key"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	encryptingText := fmt.Sprintf("%s|%s", existingUser.UserId, decoded)

	token, err := jsonwebtoken.GenerateJWT(encryptingText, utils.GetUserAgentFromContext(ctx), utils.GetSourceIPFromContext(ctx), constants.JwtKey, constants.JwtExpTime, "paydoh-bank")

	if err != nil {
		logData.Message = "GetAuthorizationTokenByUserId: Error while generating jwt token"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	response := map[string]interface{}{
		"token": token,
	}

	logData.Message = "GetAuthorizationTokenByUserId: Token generated successfully"
	logData.ResponseSize = len(token)
	logData.ResponseBody = string(token)
	logData.EndTime = time.Now()
	store.LoggerService.LogInfo(logData)

	return response, nil
}

func (store *AuthorizationStore) SmsVerification(mobileNumber, msg, simOperator string) error {
	userMobileNumber := utils.RemoveCountryCode(mobileNumber)

	user, err := models.GetUserDataByMobileNumber(store.db, userMobileNumber)
	if err != nil {
		store.LoggerService.Logger.Error("SmsVerification: error getting user" + err.Error())
		return err
	}

	msgCode, err := store.Redis.Get(user.UserId)
	if err != nil {
		store.LoggerService.LogError(&commonSrv.LogEntry{
			Message: "error getting msg code from redis",
		})
	}

	logData := &commonSrv.LogEntry{
		Action: "ROUTE_MOBILE_CALLBACK",
	}

	if msgCode != msg {
		logData.Message = "SmsVerification: msg code does not match"
		store.LoggerService.LogError(logData)

		if err := models.UpdateIsActiveAndSimVerifiedStatus(true, false, user.UserId); err != nil {
			logData.Message = "SmsVerification: Error updating device:- " + err.Error()
			store.LoggerService.LogError(logData)
		}

		if err := models.UpdateUserOnboardingStatusV2(user.UserId, false); err != nil {
			logData.Message = "SmsVerification: Error updating user onboarding status:- " + err.Error()
			store.LoggerService.LogError(logData)
			return err
		}

		return errors.New("SmsVerification: verification code does not match")
	}

	if err := models.UpdateDevice(store.db, &models.DeviceData{
		IsSimVerified: true,
		SimOperator:   types.FromString(simOperator),
	}, user.UserId); err != nil {
		store.LoggerService.Logger.Error("SmsVerification: Error updating device:- " + err.Error())
		return err
	}

	if err := models.UpdateUserOnboardingStatusV2(user.UserId, true); err != nil {
		logData.Message = "SmsVerification: Error updating user onboarding status:- " + err.Error()
		store.LoggerService.LogError(logData)
		return err
	}

	// update onboarding step only for first time user
	if err := store.UpdateAuthorizationCompleteStep(user.UserId); err != nil {
		logData.Message = "SmsVerification: Error updating user onboarding status:- " + err.Error()
		store.LoggerService.LogError(logData)
	}

	store.LoggerService.Logger.Info("SmsVerification successfully updated")

	return nil
}

func (store *AuthorizationStore) UpdateAuthorizationCompleteStep(userId string) error {
	res, err := models.GetDataByStepName(constants.PERSONAL_DETAILS_STEP)
	if err != nil {
		return err
	}

	_, err = store.db.Exec(
		`
			UPDATE user_onboarding_status
			SET current_step_id = $1, updated_at = NOW()
			WHERE user_id = $2 AND is_sim_verification_complete = true AND current_stage_id = $3 
		`,
		res.ID,
		userId,
		res.StageId,
	)
	if err != nil {
		return err
	}

	return nil
}

// insert user in mogngo db
func (store *AuthorizationStore) InsertUserInMongoDB(user *models.User) error {
	if _, err := store.Mongo.InsertOne(
		constants.MONGO_USER_DB,
		constants.MONGO_USER_COLLECTION,
		user,
		true,
	); err != nil {
		return err
	}
	return nil
}

// insert prefrence in Mongo
func (store *AuthorizationStore) InsertUserNotificationPreferencesInMongoDB(user *models.User) error {
	masterData, err := store.GetMasterNotificationType()
	if err != nil {
		return err
	}

	request := models.NewNotificationPreferences()
	request.UserID = user.UserId
	if err := request.Bind1(masterData); err != nil {
		return err
	}

	_, err = store.Mongo.InsertOne(constants.MONGO_USER_DB, constants.NotificationPref_Collection, request, false)
	if err != nil {

		return err
	}
	return nil
}

func (s *AuthorizationStore) GetMasterNotificationType() ([]models.AddNotificationPreferences, error) {

	masterNotification, err := s.Mongo.Find(constants.MONGO_USER_DB, constants.NotificationPrefrenceMaster, bson.M{}, nil, nil, 0, 0, false)
	if err != nil {

		return nil, err
	}

	notificationPrefs := []models.AddNotificationPreferences{}

	if err := masterNotification.All(s.Mongo.GetContext(), &notificationPrefs); err != nil {

		return nil, err
	}

	return notificationPrefs, nil
}

// update user in mogngo db
func (store *AuthorizationStore) UpdateUserInMongoDB(userId string, updateData interface{}) error {
	if _, err := store.Mongo.UpdateOne(
		constants.MONGO_USER_DB,
		constants.MONGO_USER_COLLECTION,
		bson.M{
			"user_id": userId,
		},
		updateData,
		false,
		false,
	); err != nil {
		return err
	}
	return nil
}
