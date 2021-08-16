package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Wallet struct {
	gorm.Model
	Amount float64 `json:"amount"`
	User   User    `gorm:"foreignKey:UserID" json:"user"`
	UserID uint    `json:"user_id"`
}

func (w *Wallet) Validate() error {
	if w.Amount < 0 {
		return errors.New("amount is invalid")
	}
	return nil
}

func GetWalletByUserId(user_id int, db *gorm.DB) (*Wallet, error) {
	wallet := &Wallet{}
	if err := db.Debug().Preload("User").Table("wallets").Where("user_id = ?", user_id).First(wallet).Error; err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w *Wallet) UpdateWalletFromTopUpByUserId(user_id int, amount float64, db *gorm.DB) (*Wallet, error) {
	if err := db.Debug().Table("wallets").Where("user_id = ?", user_id).Updates(Wallet{
		Amount: amount,
	}).Error; err != nil {
		return &Wallet{}, err
	}
	return w, nil
}

func UpdateWalletByUserId(user_id int, amount float64, db *gorm.DB) (*Wallet, error) {
	w, err := GetWalletByUserId(user_id,db)
	if err != nil {
		return &Wallet{}, errors.New("wallet not found")
	}
	if w.Amount+amount < 0 {
		return &Wallet{}, errors.New("wallet amount not enough")
	}

	w.Amount += amount
	if err := db.Debug().Table("wallets").Where("user_id = ?", user_id).Updates(Wallet{
		Amount: w.Amount,
	}).Error; err != nil {
		return &Wallet{}, err
	}
	return &w, nil
}