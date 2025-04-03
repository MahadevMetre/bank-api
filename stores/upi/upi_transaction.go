package upi

// import (
// 	"bankapi/models"
// 	"database/sql"
// )

// type UPITransactionService struct {
// 	db *sql.DB
// }

// func NewUPITransactionService(db *sql.DB) *UPITransactionService {
// 	return &UPITransactionService{db: db}
// }

// func (s *UPITransactionService) CreateTransaction(userID, transactionID, credType string) (*models.UPITransaction, error) {
// 	transaction := &models.UPITransaction{
// 		UserID:        userID,
// 		TransactionID: transactionID,
// 		CredType:      credType,
// 	}

// 	err := models.CreateUPITransaction(s.db, transaction)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return transaction, nil
// }
