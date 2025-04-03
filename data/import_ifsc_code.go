package main

import (
	"bankapi/config"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"bitbucket.org/paydoh/paydoh-commons/settings"
	_ "github.com/lib/pq"
)

type IFSCRecord struct {
	IFSCCode      string
	BankName      string
	BranchName    string
	BranchCity    string
	BranchState   string
	BranchCountry string
	PaymentMode   string
}

func main() {
	settings.LoadEnvFile()

	db := config.InitDB()

	defer db.Close()

	if err := ImportIFSCData(db, "data/ifsc_list.csv"); err != nil {
		fmt.Printf("Failed to import IFSC codes: %v\n", err)
		os.Exit(1)
	}
}

// ImportIFSCData reads IFSC records from a CSV file and inserts them into the database.
func ImportIFSCData(db *sql.DB, filePath string) error {
	records, err := ReadIFSCRecords(filePath)
	if err != nil {
		return fmt.Errorf("failed to read IFSC records: %v", err)
	}

	if err := InsertIFSCRecords(db, records); err != nil {
		return fmt.Errorf("failed to insert IFSC records: %v", err)
	}

	fmt.Printf("✅ Successfully imported %d records\n", len(records))
	return nil
}

// ReadIFSCRecords reads a CSV file and returns a slice of IFSCRecord.
func ReadIFSCRecords(filePath string) ([]IFSCRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	var records []IFSCRecord
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %v", err)
		}

		if len(record) < 7 {
			return nil, fmt.Errorf("invalid CSV record: expected 7 fields, got %d", len(record))
		}

		records = append(records, IFSCRecord{
			IFSCCode:      record[0],
			BankName:      record[1],
			BranchName:    record[2],
			BranchCity:    record[3],
			BranchState:   record[4],
			BranchCountry: record[5],
			PaymentMode:   record[6],
		})
	}

	return records, nil
}

// InsertIFSCRecords inserts multiple IFSC records into the database in a transaction.
func InsertIFSCRecords(db *sql.DB, records []IFSCRecord) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO ifsc_data (ifsc_code, bank_name, branch_name, branch_city, branch_state, branch_country, payment_mode) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	for i, record := range records {
		_, err = stmt.Exec(
			record.IFSCCode,
			record.BankName,
			record.BranchName,
			record.BranchCity,
			record.BranchState,
			record.BranchCountry,
			record.PaymentMode,
		)
		if err != nil {
			return fmt.Errorf("error inserting record at index %d: %v", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
