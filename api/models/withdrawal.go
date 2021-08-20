package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Withdrawal struct {
	gorm.Model
	Amount            float64         `json:"amount"`
	DonationProgram   DonationProgram `gorm:"foreignKey:donation_program_id" json:"donation"`
	DonationProgramID uint            `json:"donation_program_id"`
	User              User            `gorm:"foreignKey:UserID" json:"user"`
	UserID            uint            `json:"user_id"`
	Status            Status          `gorm:"type:Status; default:'pending'" json:"status,omitempty"`
}

type WithdrawalModel struct {
	ID        uint      `json:"ID"`
	Amount    float64   `json:"amount"`
	Title     string    `json:"title"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (w *Withdrawal) CreateWithdrawal(db *gorm.DB) (*Withdrawal, error) {
	if err := db.Debug().Create(&w).Error; err != nil {
		return &Withdrawal{}, err
	}
	return w, nil
}

func (w *Withdrawal) Prepare() {
	w.User = User{}
	w.DonationProgram = DonationProgram{}
}

func (w *Withdrawal) Validate() error {
	if w.Amount <= 0 {
		return errors.New("amount is invalid")
	}
	return nil
}

func GetWithdrawalById(id int, db *gorm.DB) (*Withdrawal, error) {
	withdrawal := &Withdrawal{}
	if err := db.Debug().Preload("DonationProgram").Preload("User").Table("withdrawals").Where("id = ?", id).First(withdrawal).Error; err != nil {
		return nil, err
	}
	return withdrawal, nil
}

func (w *Withdrawal) VerifyWithdrawal(id int, db *gorm.DB) (*Withdrawal, error) {
	if err := db.Debug().Table("withdrawals").Where("id = ?", id).Updates(Withdrawal{Status: "verified"}).Error; err != nil {
		return &Withdrawal{}, err
	}
	return w, nil
}

func GetUnverifiedWithdrawal(db *gorm.DB) (*[]WithdrawalModel, error) {
	withdrawals := []WithdrawalModel{}
	if err := db.Debug().Table("withdrawals").Select("withdrawals.*, u.name AS name, dp.title AS title").
		Where("withdrawals.status = ?", "pending").
		Joins("LEFT JOIN users u ON withdrawals.user_id = u.id").
		Joins("LEFT JOIN donation_programs dp ON dp.id = withdrawals.donation_program_id").Find(&withdrawals).Error; err != nil {
		return nil, err
	}
	return &withdrawals, nil
}
