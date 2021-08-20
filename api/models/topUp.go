package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type TopUp struct {
	gorm.Model
	Amount float64 `json:"amount"`
	User   User    `gorm:"foreignKey:UserID" json:"user"`
	UserID uint    `json:"user_id"`
}

type TopUpAPI struct {
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *TopUp) Validate() error {
	if t.Amount < 0 {
		return errors.New("amount is invalid")
	}
	return nil
}

func (t *TopUp) CreateTopUp(db *gorm.DB) (*TopUp, error) {
	if err := db.Debug().Create(&t).Error; err != nil {
		return &TopUp{}, err
	}
	return t, nil
}

func TopupHistory(user_id int, db *gorm.DB) (*[]TopUpAPI, error) {
	topups := []TopUpAPI{}
	if err := db.Debug().Select("Amount,created_at").Table("top_ups").Where("user_id = ?", user_id).Find(&topups).Error; err != nil {
		return nil, err
	}
	return &topups, nil
}
