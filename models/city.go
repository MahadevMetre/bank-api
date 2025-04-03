package models

import (
	"bankapi/constants"
	"database/sql"
)

type City struct {
	Id       string `json:"id"`
	City     string `json:"city_name"`
	State    string `json:"state_name"`
	CityCode string `json:"city_code"`
	ZipCode  string `json:"zip_code"`
}

func NewCity() *City {
	return &City{}
}

func FindCityByState(db *sql.DB, state string) ([]City, error) {
	cities := make([]City, 0)

	row, err := db.Query("SELECT id, city_name, state_name, city_code, zip_code FROM state_city_master WHERE state_name = $1", state)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		city := NewCity()

		if err := row.Scan(
			&city.Id,
			&city.City,
			&city.State,
			&city.CityCode,
			&city.ZipCode,
		); err != nil {
			return nil, err
		}

		cities = append(cities, *city)
	}

	return cities, nil
}

func FindCityDataByCityName(db *sql.DB, cityname string) ([]City, error) {
	cities := make([]City, 0)

	row, err := db.Query("SELECT id, city_name, state_name, city_code, zip_code FROM state_city_master WHERE LOWER(city_name) = $1", cityname)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		city := NewCity()

		if err := row.Scan(
			&city.Id,
			&city.City,
			&city.State,
			&city.CityCode,
			&city.ZipCode,
		); err != nil {
			return nil, err
		}

		cities = append(cities, *city)
	}

	return cities, nil
}
