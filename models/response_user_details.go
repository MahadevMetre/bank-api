package models

import (
	"bankapi/constants"
	"bankapi/requests"
	"bankapi/responses"
	jsonEncoder "encoding/json"
	"strings"

	"bitbucket.org/paydoh/paydoh-commons/pkg/json"
)

type ResponseUserDetails struct {
	First_name              string                 `json:"first_name"`
	Last_name               string                 `json:"last_name"`
	Email                   string                 `json:"email"`
	Gender                  string                 `json:"gender"`
	Account_Holder          string                 `json:"acc_holder"`
	Customer_ID             string                 `json:"customer_id"`
	Account_number          string                 `json:"acc_number"`
	IFSC                    string                 `json:"ifsc"`
	UPI_ID                  string                 `json:"upi_id"`
	Bank_Name               string                 `json:"bank_name"`
	Branch_Name             string                 `json:"branch_name"`
	Branch_Add              string                 `json:"Branch_Add"`
	Date_of_birth           string                 `json:"dob"`
	CommunicationAddress    string                 `json:"communication_address"`
	AddressAsPerAadhaar     string                 `json:"Add_asper_Aadhar"`
	Is_JointHolder          string                 `json:"IsJointHolder"`
	Ac_balance              string                 `json:"account_balance"`
	Nominee_Name            string                 `json:"Nominee_Name"`
	IsEmailVerified         bool                   `json:"is_email_verified"`
	IsAccountDetailMailSent bool                   `json:"is_account_detail_email_sent"`
	ProfilePic              jsonEncoder.RawMessage `json:"profile_pic"`
}

func (r *ResponseUserDetails) Bind(personalInformation *UserPersonalInformationAndAccountData, accountDetail *responses.AccountDetailResponse, nomineeName string) error {

	permanentAddress := TrimLine(accountDetail.AddressOne, accountDetail.AddressTwo, accountDetail.AddressThree, accountDetail.City, accountDetail.State, "India", accountDetail.Pincode)

	communicationAddr := requests.CommunicationAddress{}
	if personalInformation.CommunicationAddress.Valid {
		if err := json.Unmarshal([]byte(personalInformation.CommunicationAddress.String), &communicationAddr); err != nil {
			return err
		}
	}

	if !personalInformation.IsAddrSameAsAdhaar && communicationAddr.City != "" {
		r.CommunicationAddress = TrimLine(
			communicationAddr.HouseNo,
			communicationAddr.StreetName,
			communicationAddr.Locality,
			communicationAddr.Landmark,
			communicationAddr.City,
			communicationAddr.State,
			"India",
			communicationAddr.PinCode)
	} else {
		r.CommunicationAddress = permanentAddress
	}

	r.First_name = personalInformation.FirstName
	r.Last_name = personalInformation.LastName
	r.Email = personalInformation.Email
	r.Gender = personalInformation.Gender
	r.Account_Holder = accountDetail.AccountTitle
	r.Customer_ID = accountDetail.CustomerId
	r.Account_number = accountDetail.AccountNumber
	r.IFSC = constants.IfscCode
	r.UPI_ID = personalInformation.UpiId.String
	r.Bank_Name = constants.BankName
	r.Branch_Name = constants.BranchName
	r.Branch_Add = constants.BranchAddress
	r.Date_of_birth = accountDetail.Dob
	r.AddressAsPerAadhaar = permanentAddress
	r.Is_JointHolder = "1"
	r.Ac_balance = accountDetail.AccountBalance
	r.Nominee_Name = nomineeName
	r.IsEmailVerified = personalInformation.IsEmailVerified
	r.IsAccountDetailMailSent = personalInformation.IsAccountDetailMailSent
	r.IsAccountDetailMailSent = personalInformation.IsAccountDetailMailSent
	r.ProfilePic = personalInformation.ProfilePic

	return nil
}

func TrimLine(fullAddress ...string) string {
	var nameParts []string
	for _, n := range fullAddress {
		trimmed := strings.TrimSpace(n)
		if trimmed != "" {
			nameParts = append(nameParts, trimmed)
		}
	}

	name := strings.Join(nameParts, ", ")

	return name
}
