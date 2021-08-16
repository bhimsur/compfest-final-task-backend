package models

import "github.com/jinzhu/gorm"

type Donation struct {
	gorm.Model
	Amount            float64         `json:"amount"`
	UserID            uint            `json:"user_id"`
	DonationProgram   DonationProgram `gorm:"foreignKey:DonationProgramID" json:"donation_program"`
	DonationProgramID uint            `json:"donation_program_id"`
}

func GetDonationHistoryFromUser(user_id int, db *gorm.DB) (*[]Donation, error) {
	donations := []Donation{}
	if err := db.Debug().Preload("DonationProgram").Table("donations").Where("user_id = ?", user_id).Find(&donations).Error; err != nil {
		return &[]Donation{}, err
	}
	return &donations, nil
}

func (d *Donation) SaveDonation(db *gorm.DB) (*Donation, error) {
	w, err := UpdateWalletByUserId(d.UserID,-d.Amount,db)
	if err != nil {
		return &Donation{}, err
	}

	d, err := UpdateDonationProgramById(d.DonationProgramID,d.Amount,db)
	if err != nil {
		return &Donation{}, err
	}

	if err := db.Debug().Create(&d).Error; err != nil {
		return &Donation{}, err
	}
	return d, nil
}
