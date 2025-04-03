package models

import (
	"bankapi/constants"
	"bankapi/requests"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type KycConsent struct {
	Id                        uuid.UUID `json:"id"`
	UserId                    string    `json:"user_id"`
	IndianResident            bool      `json:"indian_resident"`
	PoliticallyExposedPerson  bool      `json:"politically_exposed_person"`
	AadharConsent             bool      `json:"aadhar_consent"`
	VirtualDebitCardConsent   bool      `json:"virtual_card_consent"`
	PhysicalDebitCardConsent  bool      `json:"physical_card_consent"`
	Aadhar2Consent            bool      `json:"aadhar2_consent"`
	AddressChangeConsent      bool      `json:"address_change_consent"`
	NominationConsent         bool      `json:"nomination_consent"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
	LocationConsent           bool      `json:"location_consent"`
	PrivacyPolicyConsent      bool      `json:"privacy_policy_consent"`
	TermsAndCondition         bool      `json:"terms_and_condition"`
	BankSmsVerificationStatus bool      `json:"bank_sms_verification_status"`
}

func NewKycConsent() *KycConsent {
	return &KycConsent{}
}

func (k *KycConsent) Bind(request *requests.KycConsentRequest, userId string) {
	k.UserId = userId
	k.IndianResident = request.IndianResident
	k.PoliticallyExposedPerson = request.NotPoliticallyExposed
	k.AadharConsent = request.AadharConsent
	k.VirtualDebitCardConsent = request.VirtualDebitCardConsent
	k.PhysicalDebitCardConsent = request.PhysicalDebitCardConsent
	k.Aadhar2Consent = request.Aadhar2Consent
	k.AddressChangeConsent = request.AddressChangeConsent
	k.NominationConsent = request.NominationConsent
	k.LocationConsent = request.LocationConsent
	k.PrivacyPolicyConsent = request.PrivacyPolicyConsent
	k.TermsAndCondition = request.TermsAndCondition
}

func (k *KycConsent) Marshal() ([]byte, error) {
	return json.Marshal(k)
}

func (k *KycConsent) UnMarshal(data []byte) error {
	return json.Unmarshal(data, k)
}

func FindKycConsentByUserId(db *sql.DB, userId string) (*KycConsent, error) {
	consent := NewKycConsent()
	row := db.QueryRow("SELECT id, user_id, indian_resident, politically_exposed_person, aadhar_consent, virtual_card_consent, physical_card_consent, aadhar2_consent,address_change_consent,nomination_consent,location_consent,privacy_policy_consent,terms_and_condition,created_at, updated_at FROM kyc_consent WHERE user_id = $1", userId)

	if err := row.Scan(
		&consent.Id,
		&consent.UserId,
		&consent.IndianResident,
		&consent.PoliticallyExposedPerson,
		&consent.AadharConsent,
		&consent.VirtualDebitCardConsent,
		&consent.PhysicalDebitCardConsent,
		&consent.Aadhar2Consent,
		&consent.AddressChangeConsent,
		&consent.NominationConsent,
		&consent.LocationConsent,
		&consent.PrivacyPolicyConsent,
		&consent.TermsAndCondition,
		&consent.CreatedAt,
		&consent.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrKycConsentNotProvided
		}
		return nil, err
	}

	return consent, nil
}

func FindKycConsentByUserIdV2(db *sql.DB, userId string) (*KycConsent, error) {
	consent := NewKycConsent()

	query := `SELECT
        kc.id,
        kc.user_id,
        kc.indian_resident,
        kc.politically_exposed_person,
        kc.aadhar_consent,
		kc.virtual_card_consent,
        kc.physical_card_consent,
        kc.created_at,
        kc.updated_at,
        kc.bank_sms_verification_status
    FROM
		kyc_consent kc
    WHERE
        kc.user_id= $1`

	row := db.QueryRow(query, userId)

	if err := row.Scan(
		&consent.Id,
		&consent.UserId,
		&consent.IndianResident,
		&consent.PoliticallyExposedPerson,
		&consent.AadharConsent,
		&consent.VirtualDebitCardConsent,
		&consent.PhysicalDebitCardConsent,
		&consent.CreatedAt,
		&consent.UpdatedAt,
		&consent.BankSmsVerificationStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrKycConsentNotProvided
		}
		return nil, err
	}

	return consent, nil
}

func InsertKycConsent(db *sql.DB, model *KycConsent) error {
	if _, err := db.Exec(
		"INSERT INTO kyc_consent (user_id, indian_resident, politically_exposed_person, aadhar_consent, virtual_card_consent, physical_card_consent) VALUES ($1, $2, $3, $4, $5, $6)",
		model.UserId,
		model.IndianResident,
		model.PoliticallyExposedPerson,
		model.AadharConsent,
		model.VirtualDebitCardConsent,
		model.PhysicalDebitCardConsent,
	); err != nil {
		return err
	}

	return nil
}

func UpdateKycConsent(db *sql.DB, updateModel *KycConsent, userId string) error {

	var clause []string
	var params []interface{}

	if updateModel.IndianResident {
		clause = append(clause, "indian_resident = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.IndianResident)
	}

	if updateModel.PoliticallyExposedPerson {
		clause = append(clause, "politically_exposed_person = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.PoliticallyExposedPerson)
	}

	if updateModel.AadharConsent {
		clause = append(clause, "aadhar_consent = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.AadharConsent)
	}

	if updateModel.VirtualDebitCardConsent {
		clause = append(clause, "virtual_card_consent = $"+strconv.Itoa(len(params)+1))
		params = append(params, updateModel.VirtualDebitCardConsent)
	}

	if updateModel.PhysicalDebitCardConsent {
		clause = append(clause, "physical_card_consent = $"+strconv.Itoa(len(params)+1))
		params = append(params, updateModel.PhysicalDebitCardConsent)
	}

	if updateModel.Aadhar2Consent {
		clause = append(clause, "aadhar2_consent = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.Aadhar2Consent)
	}

	if updateModel.AddressChangeConsent {
		clause = append(clause, "address_change_consent = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.AddressChangeConsent)
	}

	if updateModel.NominationConsent {
		clause = append(clause, "nomination_consent = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.NominationConsent)
	}

	if updateModel.LocationConsent {
		clause = append(clause, "location_consent = $"+strconv.Itoa(len(params)+1))
		params = append(params, updateModel.LocationConsent)
	}

	if updateModel.PrivacyPolicyConsent {
		clause = append(clause, "privacy_policy_consent = $"+strconv.Itoa(len(params)+1))
		params = append(params, updateModel.PrivacyPolicyConsent)
	}

	if updateModel.TermsAndCondition {
		clause = append(clause, "terms_and_condition = $"+strconv.Itoa(len(params)+1))
		params = append(params, updateModel.TermsAndCondition)
	}

	if updateModel.BankSmsVerificationStatus {
		clause = append(clause, "bank_sms_verification_status = $"+strconv.Itoa(len(params)+1))
		params = append(params, updateModel.BankSmsVerificationStatus)
	}

	if len(clause) > 0 {
		clause = append(clause, "updated_at=NOW()")
		clauseStr := strings.Join(clause, ", ")
		params = append(params, userId)
		query := "UPDATE kyc_consent SET " + clauseStr + " WHERE user_id = $" + strconv.Itoa(len(params))
		if _, err := db.Exec(query, params...); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func UpsertKycConsent(db *sql.DB, userId string, consent *KycConsent) error {
	_, err := FindKycConsentByUserId(db, userId)
	if err != nil {
		if err == constants.ErrKycConsentNotProvided {
			if !consent.AadharConsent {
				return errors.New("aadhar consent is required")
			}
			if !consent.IndianResident {
				return errors.New("indian resident consent is required")
			}
			if !consent.PoliticallyExposedPerson {
				return errors.New("politically exposed person consent is required")
			}
			return InsertKycConsent(db, consent)
		}
		return fmt.Errorf("error finding kyc consent: %v", err)
	}
	return UpdateKycConsent(db, consent, userId)
}
