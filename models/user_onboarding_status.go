package models

import (
	"bankapi/config"
	"bankapi/constants"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserOnboardingStatus struct {
	ID                            uuid.UUID `json:"id" db:"id"`
	UserID                        string    `json:"user_id" db:"user_id"`
	CurrentStageID                uuid.UUID `json:"-" db:"current_stage_id"`
	CurrentStageName              string    `json:"current_stage_name"`
	CurrentStepName               string    `json:"current_step_name"`
	CurrentStepID                 uuid.UUID `json:"-" db:"current_step_id,omitempty"`
	IsSimVerificationComplete     bool      `json:"is_sim_verification_complete"`
	IsKycConsentComplete          bool      `json:"is_kyc_consent_complete"`
	IsDemographicFetchComplete    bool      `json:"is_demographic_fetch_complete"`
	IsAccountCreationComplete     bool      `json:"is_account_creation_complete"`
	IsAccountAuditComplete        bool      `json:"is_account_audit_complete"`
	IsDebitCardConsentComplete    bool      `json:"is_debit_card_consent_complete"`
	IsDebitCardPaymentComplete    bool      `json:"is_debit_card_payment_complete"`
	IsDebitCardGenerationComplete bool      `json:"is_debit_card_generation_complete"`
	IsUpiGenerationComplete       bool      `json:"is_upi_generation_complete"`
	IsUpiPinSetupComplete         bool      `json:"is_upi_pin_setup_complete"`
	IsMPinSetupComplete           bool      `json:"is_m_pin_setup_complete"`
	CreatedAt                     time.Time `json:"-" db:"created_at"`
	UpdatedAt                     time.Time `json:"-" db:"updated_at"`
}

type OnboardingStage struct {
	ID         uuid.UUID `json:"id"`
	StageOrder int       `json:"stage_order"`
	StageName  string    `json:"stage_name"`
}

type UserOnboardingStep struct {
	ID        uuid.UUID `json:"id"`
	StepName  string    `json:"step_name"`
	StepOrder int       `json:"step_order"`
	StageId   uuid.UUID `json:"stage_id"`
}

type StageStepDetails struct {
	StageID    uuid.UUID `json:"stage_id"`
	StepID     uuid.UUID `json:"id"`
	StageOrder int       `json:"stage_order"`
	StepOrder  int       `json:"step_order"`
	StageName  string    `json:"stage_name"`
}

func GetDataByStepName(stepName string) (*UserOnboardingStep, error) {
	db := config.GetDB()

	query := `
			SELECT oss.id, oss.step_name, oss.stage_id, oss.step_order
			FROM onboarding_stage_steps oss
			WHERE oss.step_name = $1;
	`

	var step UserOnboardingStep
	row := db.QueryRow(query, stepName)

	err := row.Scan(&step.ID, &step.StepName, &step.StageId, &step.StepOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to map data: %w", err)
	}

	return &step, nil
}

func GetDataByStageName(stageName string) (*OnboardingStage, error) {
	db := config.GetDB()

	query := `
		SELECT oss.id, oss.stage_name, oss.stage_order
		FROM onboarding_stages oss
		WHERE oss.stage_name = $1;
	`

	var step OnboardingStage
	row := db.QueryRow(query, stageName)

	err := row.Scan(&step.ID, &step.StageName, &step.StageOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to map data: %w", err)
	}

	return &step, nil
}

func InsertUserOnboardingStatus(status UserOnboardingStatus) error {
	db := config.GetDB()
	query := `INSERT INTO user_onboarding_status (
		user_id, current_stage_id, current_step_id, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, status.UserID, status.CurrentStageID, status.CurrentStepID,
		status.CreatedAt, status.UpdatedAt)
	return err
}

func GetUserOnboardingStatus(userId string) (*UserOnboardingStatus, error) {
	db := config.GetDB()
	query := `
	SELECT 
		uos.id, 
		uos.user_id, 
		os.stage_name,
		COALESCE(oss.step_name, '') AS step_name,  -- Handle NULL step_name
		uos.is_sim_verification_complete, 
		uos.is_kyc_consent_complete, 
		uos.is_demographic_fetch_complete, 
		uos.is_account_creation_complete, 
		uos.is_account_audit_complete, 
		uos.is_debit_card_consent_complete, 
		uos.is_debit_card_payment_complete,
		uos.is_debit_card_generation_complete, 
		uos.is_upi_generation_complete, 
		uos.is_upi_pin_setup_complete, 
		uos.is_m_pin_setup_complete,
		uos.created_at,
		uos.updated_at
	FROM 
		user_onboarding_status uos
	JOIN 
		onboarding_stages os ON uos.current_stage_id = os.id
	LEFT JOIN 
		onboarding_stage_steps oss ON uos.current_step_id = oss.id  -- LEFT JOIN to allow NULL steps
	WHERE 
		uos.user_id = $1;
	`

	row := db.QueryRow(query, userId)

	var status UserOnboardingStatus
	err := row.Scan(
		&status.ID,
		&status.UserID,
		&status.CurrentStageName,
		&status.CurrentStepName,
		&status.IsSimVerificationComplete,
		&status.IsKycConsentComplete,
		&status.IsDemographicFetchComplete,
		&status.IsAccountCreationComplete,
		&status.IsAccountAuditComplete,
		&status.IsDebitCardConsentComplete,
		&status.IsDebitCardPaymentComplete,
		&status.IsDebitCardGenerationComplete,
		&status.IsUpiGenerationComplete,
		&status.IsUpiPinSetupComplete,
		&status.IsMPinSetupComplete,
		&status.CreatedAt,
		&status.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no onboarding status found for user %s", userId)
		}
		return nil, fmt.Errorf("failed to map data: %w", err)
	}

	return &status, nil
}

func UserOnboardingStatusAlreadyExists(userID string) (bool, error) {
	db := config.GetDB()
	query := `SELECT EXISTS(SELECT 1 FROM user_onboarding_status WHERE user_id = $1)`

	var exists bool
	err := db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking user onboarding status: %w", err)
	}

	return exists, nil
}

func CreateUserOnboardingStatus(stepName, userId string) error {
	data, err := GetDataByStepName(stepName)
	if err != nil {
		return err
	}

	status := UserOnboardingStatus{
		UserID:         userId,
		CurrentStageID: data.StageId,
		CurrentStepID:  data.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = InsertUserOnboardingStatus(status)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserOnboardingStatusV2(userId string, isSimVerificationComplete bool) error {
	db := config.GetDB()
	_, err := db.Exec("UPDATE user_onboarding_status SET is_sim_verification_complete = $1, updated_at = NOW() WHERE user_id = $2", isSimVerificationComplete, userId)
	return err
}

func UpdateUserOnboardingStatus(stepOrStageName, userId string) error {
	var details StageStepDetails
	db := config.GetDB()

	query := `
	SELECT os.id AS stage_id, COALESCE(oss.id, NULL) AS step_id, os.stage_order, COALESCE(oss.step_order, NULL) AS step_order, os.stage_name
	FROM onboarding_stage_steps oss
	RIGHT JOIN onboarding_stages os ON oss.stage_id = os.id
	WHERE oss.step_name = $1
	`

	err := db.QueryRow(query, stepOrStageName).Scan(&details.StageID, &details.StepID, &details.StageOrder, &details.StepOrder, &details.StageName)
	if err != nil {
		stageQuery := `
		SELECT id AS stage_id, stage_order, stage_name 
		FROM onboarding_stages 
		WHERE stage_name = $1
		`
		err = db.QueryRow(stageQuery, stepOrStageName).Scan(&details.StageID, &details.StageOrder, &details.StageName)
		if err != nil {
			return fmt.Errorf("error retrieving stage by stepName or stageName: %v", err)
		}
		details.StepID = uuid.Nil
	}

	updateQuery := `
	UPDATE user_onboarding_status
	SET current_stage_id = $1,
		current_step_id = $2,
		updated_at = now()
	WHERE user_id = $3
	`
	_, err = db.Exec(updateQuery, details.StageID, details.StepID, userId)
	if err != nil {
		return fmt.Errorf("error updating user onboarding status: %v", err)
	}

	if details.StepID == uuid.Nil {
		updateCompletionQuery := `
			UPDATE user_onboarding_status
			SET is_sim_verification_complete = CASE WHEN $1 = 1 THEN TRUE ELSE is_sim_verification_complete END,
				is_kyc_consent_complete = CASE WHEN $1 = 2 THEN TRUE ELSE is_kyc_consent_complete END,
				is_demographic_fetch_complete = CASE WHEN $1 = 3 THEN TRUE ELSE is_demographic_fetch_complete END,
				is_account_creation_complete = CASE WHEN $1 = 4 THEN TRUE ELSE is_account_creation_complete END,
				is_debit_card_consent_complete = CASE WHEN $1 = 5 THEN TRUE ELSE is_debit_card_consent_complete END,
				is_debit_card_payment_complete = CASE WHEN $1 = 6 THEN TRUE ELSE is_debit_card_payment_complete END,
				is_debit_card_generation_complete = CASE WHEN $1 = 7 THEN TRUE ELSE is_debit_card_generation_complete END,
				is_upi_generation_complete = CASE WHEN $1 = 8 THEN TRUE ELSE is_upi_generation_complete END,
				is_upi_pin_setup_complete = CASE WHEN $1 = 9 THEN TRUE ELSE is_upi_pin_setup_complete END,
				is_m_pin_setup_complete = CASE WHEN $1 = 10 THEN TRUE ELSE is_m_pin_setup_complete END,
				updated_at = now()
			WHERE user_id = $2
		`
		_, err = db.Exec(updateCompletionQuery, details.StageOrder, userId)
		if err != nil {
			return fmt.Errorf("error updating completion flag: %v", err)
		}
		return nil
	}

	// Check if the user has completed all steps in the current stage (if steps exist)
	checkCompletionQuery := `
	SELECT COUNT(*)
	FROM onboarding_stage_steps oss
	JOIN user_onboarding_status uos ON oss.stage_id = uos.current_stage_id
	WHERE uos.user_id = $1 AND oss.step_order > (
			SELECT step_order
			FROM onboarding_stage_steps
			WHERE id = $2
	)
	`

	var remainingSteps int
	err = db.QueryRow(checkCompletionQuery, userId, details.StepID).Scan(&remainingSteps)
	if err != nil {
		return fmt.Errorf("error checking remaining steps: %v", err)
	}

	if remainingSteps == 0 {
		updateCompletionQuery := `
			UPDATE user_onboarding_status
			SET is_sim_verification_complete = CASE WHEN $1 = 1 THEN TRUE ELSE is_sim_verification_complete END,
				is_kyc_consent_complete = CASE WHEN $1 = 2 THEN TRUE ELSE is_kyc_consent_complete END,
				is_demographic_fetch_complete = CASE WHEN $1 = 3 THEN TRUE ELSE is_demographic_fetch_complete END,
				is_account_creation_complete = CASE WHEN $1 = 4 THEN TRUE ELSE is_account_creation_complete END,
				is_debit_card_consent_complete = CASE WHEN $1 = 5 THEN TRUE ELSE is_debit_card_consent_complete END,
				is_debit_card_payment_complete = CASE WHEN $1 = 6 THEN TRUE ELSE is_debit_card_payment_complete END,
				is_debit_card_generation_complete = CASE WHEN $1 = 7 THEN TRUE ELSE is_debit_card_generation_complete END,
				is_upi_generation_complete = CASE WHEN $1 = 8 THEN TRUE ELSE is_upi_generation_complete END,
				is_upi_pin_setup_complete = CASE WHEN $1 = 9 THEN TRUE ELSE is_upi_pin_setup_complete END,
				is_m_pin_setup_complete = CASE WHEN $1 = 10 THEN TRUE ELSE is_m_pin_setup_complete END,
				updated_at = now()
			WHERE user_id = $2
		`
		_, err = db.Exec(updateCompletionQuery, details.StageOrder, userId)
		if err != nil {
			return fmt.Errorf("error updating completion flag: %v", err)
		}
	}

	if stepOrStageName == constants.AUDIT_CALLBACK_STEP {
		updateQuery := `
			UPDATE user_onboarding_status
			SET is_account_audit_complete = true, updated_at = now()
			WHERE user_id = $1
		`
		_, err = db.Exec(updateQuery, userId)
		if err != nil {
			return fmt.Errorf("error updating  %v", err)
		}
	}

	return nil
}
