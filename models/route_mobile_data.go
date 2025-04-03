package models

import (
	"time"
)

type RouteMobileData struct {
	Id           int64     `json:"id"`
	UserId       string    `json:"user_id"`
	MobileNumber string    `json:"mobile_number"`
	Operator     string    `json:"operator"`
	Circle       string    `json:"circle"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewRouteMobileData() *RouteMobileData {
	return &RouteMobileData{}
}

func (d *RouteMobileData) Create(userId, mobileNumber, operator, circle string) {
	d.UserId = userId
	d.MobileNumber = mobileNumber
	d.Operator = operator
	d.Circle = circle
}
