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
		wallet = &Wallet{Amount: 0, UserID: uint(user_id)}
		wallet, _ = wallet.CreateWallet(db)
	}
	return wallet, nil
}

func (w *Wallet) UpdateWallet(db *gorm.DB) (*Wallet, error) {
	if err := db.Debug().Table("wallets").Where("id = ?", w.ID).Updates(Wallet{
		Amount: w.Amount,
	}).Error; err != nil {
		return &Wallet{}, err
	}
	return w, nil
}

func (w *Wallet) CreateWallet(db *gorm.DB) (*Wallet, error) {
	err := db.Debug().Create(&w).Error
	if err != nil {
		return &Wallet{}, err
	}
	return w, nil
}
