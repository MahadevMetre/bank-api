package statement

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/security"
	"bankapi/services"
	"bankapi/stores/transaction"
	"bankapi/utils"
)

type Statement interface {
	GetAccountState(authValue *models.AuthValues)
}

type Store struct {
	db                 *sql.DB
	m                  *database.Document
	TransactionHistory *transaction.TransactionStore
	service            *services.BankApiService
	LoggerService      *commonSrv.LoggerService
	bankService        *services.BankApiService
}

type Transaction struct {
	Date    string
	Details string
	Credit  string
	Debit   string
	Balance string
}

type PDFData struct {
	UserName             string
	UserAddress          string
	UserPhone            string
	UserEmail            string
	AccountNumber        string
	AccountType          string
	IFSC                 string
	BranchName           string
	BranchAddress        string
	StatementPeriodStart string
	StatementPeriodEnd   string
	OpeningBalance       string
	MoneyIn              string
	MoneyOut             string
	ClosingBalance       string
	Transactions         []Transaction
}

func NewStore(
	log *commonSrv.LoggerService,
	db *sql.DB,
	m *database.Document,
	memory *database.InMemory,
	txnHistory *transaction.TransactionStore,
) *Store {
	bankService := services.NewBankApiService(log, memory)
	return &Store{
		db:                 db,
		m:                  m,
		service:            services.NewBankApiService(log, memory),
		LoggerService:      log,
		TransactionHistory: txnHistory,
		bankService:        bankService,
	}
}

func (s *Store) GetAccountStatement(ctx context.Context, authValue *models.AuthValues, request *requests.StatementRequest) (interface{}, error) {
	if request.Type == "pdf" {
		return s.GeneratePDF(ctx, request, authValue.UserId)
	} else if request.Type == "csv" {
		return s.GenerateCSV(ctx, request, authValue.UserId)
	}
	return nil, errors.New("Unsupported statement type")
}

func (s *Store) GetAccountState(ctx context.Context, authValue *models.AuthValues, request *requests.StatementRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.STATEMENT,
		RequestURI: "/api/statement/get-statement",
		Message:    "GetAccountState log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	accountData, err := models.GetAccountDataByUserId(s.db, authValue.UserId)

	if err != nil {
		logData.Message = "GetAccountState: Error getting account data by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	outgoingStatementRequest := requests.NewOutgoingStatementRequest()

	if err := outgoingStatementRequest.BindAndValidate(authValue.UserId, accountData.AccountNumber, request.FromDate, request.ToDate); err != nil {
		logData.Message = "GetAccountState: Error binding and validating outgoing statement request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	statementResponse, err := s.service.GetBankStatement(ctx, outgoingStatementRequest)

	if err != nil {
		logData.Message = "GetAccountState: Error getting bank statement"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if statementResponse.ErrorCode != "0" && statementResponse.ErrorCode != "00" {
		logData.Message = "GetAccountState: Error in bank statement response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(statementResponse.ErrorMessage)
	}

	byteData, err := statementResponse.Marshal()

	if err != nil {
		logData.Message = "GetAccountState: Error marshaling statement response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteData, []byte(authValue.Key))

	if err != nil {
		logData.Message = "GetAccountState: Error encrypting data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetAccountState: Account state retrieved successfully"
	logData.ResponseSize = len(byteData)
	logData.ResponseBody = string(byteData)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) GeneratePDF(ctx context.Context, request *requests.StatementRequest, userId string) ([]byte, error) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("Unable to parse template: %v", err)
	}

	accountDetails, err := models.GetUserAndAccountDetailByUserID(s.db, userId)
	if err != nil {
		return nil, err
	}

	req := requests.TransactionRequest{
		UserId:   userId,
		FromDate: request.FromDate,
		ToDate:   request.ToDate,
	}
	txnData, err := s.TransactionHistory.GetTransactionData(ctx, req)
	if err != nil {
		return nil, err
	}

	transactions := []Transaction{}
	var moneyIn float64
	var moneyOut float64

	if len(txnData) == 0 {
		return nil, errors.New("no transaction data found for the given date range")
	}

	for _, v := range txnData {
		transaction := Transaction{
			Date:    parseDateData(v.TransactionDate),
			Details: v.TransactionDescription,
			Balance: formatBalance(v.RunningTotal),
		}

		if v.CodeDRCR == "C" {
			transaction.Credit = v.TransactionAmount
			f64, err := strconv.ParseFloat(v.TransactionAmount, 64)
			if err == nil {
				fmt.Printf("Converted to float64: %f\n", f64)
			}
			moneyIn += f64
		} else if v.CodeDRCR == "D" {
			transaction.Debit = v.TransactionAmount
			f64, err := strconv.ParseFloat(v.TransactionAmount, 64)
			if err == nil {
				fmt.Printf("Converted to float64: %f\n", f64)
			}
			moneyOut += f64
		}

		transactions = append(transactions, transaction)
	}

	var communicationAddr requests.CommunicationAddress
	var mailAdr string
	if accountDetails.CommunicationAddress.String != "" {
		if err := json.Unmarshal([]byte(accountDetails.CommunicationAddress.String), &communicationAddr); err != nil {
			return nil, err
		}

		mailAdr = requests.MapAddress(
			communicationAddr.HouseNo,
			communicationAddr.StreetName,
			communicationAddr.Locality,
			communicationAddr.Landmark,
			communicationAddr.City,
			communicationAddr.State,
			communicationAddr.PinCode,
		)
	}

	data := PDFData{
		UserName:             accountDetails.FirstName + " " + accountDetails.MiddleName + " " + accountDetails.LastName,
		UserAddress:          mailAdr,
		UserPhone:            accountDetails.MobileNumber,
		UserEmail:            accountDetails.Email,
		AccountNumber:        accountDetails.AccountNumber,
		AccountType:          "Savings",
		IFSC:                 constants.IfscCode,
		BranchName:           constants.BranchName,
		BranchAddress:        constants.BranchAddress,
		StatementPeriodStart: formatDate(req.FromDate),
		StatementPeriodEnd:   formatDate(req.ToDate),
		OpeningBalance:       formatBalance(txnData[0].RunningTotal),
		MoneyIn:              formatBalance(fmt.Sprintf("%f", moneyIn)),
		MoneyOut:             formatBalance(fmt.Sprintf("%f", moneyOut)),
		ClosingBalance:       formatBalance(txnData[len(txnData)-1].RunningTotal),
		Transactions:         transactions,
	}

	var htmlContent bytes.Buffer
	err = tmpl.Execute(&htmlContent, data)
	if err != nil {
		return nil, fmt.Errorf("unable to execute template: %v", err)
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("unable to create PDF generator: %v", err)
	}

	page := wkhtmltopdf.NewPageReader(&htmlContent)

	footerHTMLPath := "templates/footer.html"
	footerHTML, err := os.ReadFile(footerHTMLPath)
	if err != nil {
		log.Fatalf("Failed to read footer HTML file: %v", err)
	}

	footerFile, err := os.CreateTemp("", "footer*.html")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(footerFile.Name())

	_, err = footerFile.Write(footerHTML)
	if err != nil {
		log.Fatal(err)
	}
	footerFile.Close()

	page.FooterHTML.Set(footerFile.Name())

	pdfg.AddPage(page)

	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(10)

	err = pdfg.Create()
	if err != nil {
		return nil, fmt.Errorf("unable to create PDF: %v", err)
	}

	return pdfg.Bytes(), nil
}

func formatBalance(bal string) string {

	balance, err := strconv.ParseFloat(bal, 64)
	if err != nil {
		fmt.Printf("Converted to float64: %f\n", balance)
		return ""
	}

	// Convert the float to a string with 2 decimal places
	balanceStr := strconv.FormatFloat(balance, 'f', 2, 64)

	// Split the string into integer and decimal parts
	parts := strings.Split(balanceStr, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	// Reverse the integer part to easily add commas
	runes := []rune(integerPart)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	reversedInt := string(runes)

	var result strings.Builder
	for i, char := range reversedInt {
		if i > 0 && i%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(char)
	}

	// Reverse back and combine with decimal part
	runes = []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	formattedInt := string(runes)

	return formattedInt + "." + decimalPart
}

func formatDate(dateStr string) string {
	parsedDate, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}

	formattedDate := parsedDate.Format("2 January 2006")
	return formattedDate
}

func (s *Store) GenerateCSV(ctx context.Context, request *requests.StatementRequest, userID string) ([]byte, error) {

	req := requests.TransactionRequest{
		UserId:   userID,
		FromDate: request.FromDate,
		ToDate:   request.ToDate,
	}
	txnData, err := s.TransactionHistory.GetTransactionData(ctx, req)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	writer := csv.NewWriter(&buffer)

	headers := []string{"Date", "Transaction details", "Credit", "Debit", "Balance"}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("error writing headers to CSV: %v", err)
	}

	records := [][]string{}

	for _, v := range txnData {
		record := []string{v.TransactionDate, v.TransactionDescription, v.TransactionAmount, v.CodeDRCR, v.RunningTotal}
		records = append(records, record)
	}

	if err := writer.WriteAll(records); err != nil {
		return nil, fmt.Errorf("error writing records to CSV: %v", err)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing CSV writer: %v", err)
	}

	return buffer.Bytes(), nil
}

func parseDateData(dateString string) string {
	layout := "02-01-2006 15:04:05"

	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}

	return parsedTime.Format("02/01/2006")
}
