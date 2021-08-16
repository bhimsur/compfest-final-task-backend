package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type TopUp struct {
	Amount float64 `json:"amount"`
	User   User    `gorm:"foreignKey:UserID" json:"user"`
	UserID uint    `json:"user_id"`
}

func (t *TopUp) Validate() error {
	if t.Amount < 0 {
		return errors.New("amount is invalid")
	}
	if t.UserID < 0 {
		return errors.New("user_id is invalid")
	}
	return nil
}

func (t *TopUp) CreateTopUp(db *gorm.DB) (*TopUp, error) {
	w, err := UpdateWalletByUserId(t.UserID,t.Amount,db)
	if err != nil {
		return &TopUp, err
	}
	
	if err := db.Debug().Create(&t).Error; err != nil {
		return &TopUp{}, err
	}
	return t, nil
}
