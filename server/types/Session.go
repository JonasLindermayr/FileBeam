package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID uint `json:"id"`
	Token string `json:"key" gorm:"unique"`
	EmployeeID uuid.UUID `json:"employeeId"`
}