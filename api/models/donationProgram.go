package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

type DonationProgram struct {
	gorm.Model
	Title    string     `json:"title"`
	Detail   string     `json:"detail"`
	Amount   float64    `json:"amount"`
	User     User       `gorm:"foreignKey:UserID" json:"user"`
	UserID   uint       `json:"user_id"`
	Status   Status     `gorm:"type:Status; default:'unverified'" json:"status"`
	Donation []Donation `gorm:"foreignKey:DonationProgramID;references:ID" json:"donations"`
}

func (d *DonationProgram) Prepare() {
	d.Title = strings.TrimSpace(d.Title)
	d.Detail = strings.TrimSpace(d.Detail)
	d.User = User{}
}

func (d *DonationProgram) Validate() error {
	if d.Title == "" {
		return errors.New("title is required")
	}
	if d.Detail == "" {
		return errors.New("detail is required")
	}
	if d.Amount <= 0 {
		return errors.New("amount is invalid")
	}
	return nil
}

func (d *DonationProgram) Save(db *gorm.DB) (*DonationProgram, error) {
	err := db.Debug().Create(&d).Error
	if err != nil {
		return &DonationProgram{}, err
	}
	return d, nil
}

func GetDonationPrograms(db *gorm.DB) (*[]DonationProgram, error) {
	donationProgram := []DonationProgram{}
	if err := db.Debug().Table("donation_programs dp").
		Unscoped().
		Preload("User").
		Preload("Donation").
		Select("dp.*, SUM(d.amount) AS donasi_terkumpul, (dp.amount-SUM(d.amount)) AS donasi_kekurangan").
		Joins("LEFT OUTER JOIN donations d ON d.donation_program_id = dp.id").
		Group("dp.id").
		Find(&donationProgram).Error; err != nil {
		return &[]DonationProgram{}, err
	}
	return &donationProgram, nil
}

func GetDonationProgramById(id int, db *gorm.DB) (*DonationProgram, error) {
	donationProgram := &DonationProgram{}
	if err := db.Debug().Preload("Donation").Preload("User").Table("donation_programs").Where("id = ?", id).First(donationProgram).Error; err != nil {
		return nil, err
	}
	return donationProgram, nil
}

func UpdateDonationProgramById(id int, amount float64 db *gorm.DB) (*DonationProgram, error) {
	d, err := GetDonationProgramById(id,db)
	if err != nil {
		return &DonationProgramID, errors.New("donation_program not found")
	}
	
	d.Amount += amount
	if err := db.Debug().Table("donation_programs").Where("id = ?", id).Updates(DonationProgram{
		Amount: d.Amount,	
	}).Error; err != nil {
		return &DonationProgram{}, err
	}
	return &d, nil
}

func (d *DonationProgram) UpdateDonationProgram(id int, db *gorm.DB) (*DonationProgram, error) {
	if err := db.Debug().Table("donation_programs").Where("id = ?", id).Updates(DonationProgram{
		Title:  d.Title,
		Detail: d.Detail,
		Amount: d.Amount,
	}).Error; err != nil {
		return &DonationProgram{}, err
	}
	return d, nil
}

func DeleteDonationProgram(id int, db *gorm.DB) error {
	if err := db.Debug().Table("donation_programs").Where("id = ?", id).Delete(&DonationProgram{}).Error; err != nil {
		return err
	}
	return nil
}

func GetDonationProgramByFundraiser(user_id int, db *gorm.DB) (*[]DonationProgram, error) {
	donationPrograms := []DonationProgram{}
	if err := db.Debug().
		Preload("User").
		Table("donation_programs").
		Where("user_id = ?", user_id).
		Find(&donationPrograms).Error; err != nil {
		return &[]DonationProgram{}, err
	}
	return &donationPrograms, nil
}

func (dp *DonationProgram) VerifyDonationProgram(id int, db *gorm.DB) (*DonationProgram, error) {
	if err := db.Debug().Table("donation_programs").Where("id = ?", id).Updates(DonationProgram{Status: "verified"}).Error; err != nil {
		return &DonationProgram{}, err
	}
	return dp, nil
}
