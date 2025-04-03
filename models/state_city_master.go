package models

import (
	"bankapi/config"
	"fmt"
)

type StateCityMaster struct {
	StateName string `json:"state_name"`
	CityName  string `json:"city_name"`
	CityCode  string `json:"citycode"`
}

// CheckZipCodeExists checks if a given zip code exists in the database
func CheckZipCodeExists(zipCode string) (bool, error) {
	db := config.GetDB()
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM state_city_master
			WHERE zip_code = $1
		) AS zip_code_exists
	`

	err := db.QueryRow(query, zipCode).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking zip code: %w", err)
	}

	return exists, nil
}

func GetCitiesByZipCode(zipCode string) ([]StateCityMaster, error) {
	db := config.GetDB()
	var cities []StateCityMaster
	query := `
		SELECT state_name, city_name, city_code from state_city_master WHERE zip_code = $1
	`

	rows, err := db.Query(query, zipCode)
	if err != nil {
		return nil, fmt.Errorf("error checking zip code: %w", err)
	}

	for rows.Next() {
		var city StateCityMaster
		if err := rows.Scan(&city.StateName, &city.CityName, &city.CityCode); err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	defer rows.Close()

	return cities, nil
}
