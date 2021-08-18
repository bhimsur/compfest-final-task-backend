package models

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// roles
type Roles string

const (
	admin      Roles = "admin"
	donor      Roles = "donor"
	fundraiser Roles = "fundraiser"
)

type Status string

const (
	verified   Status = "verified"
	pending    Status = "pending"
	unverified Status = "unverified"
)

// user model
type User struct {
	gorm.Model
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     Roles  `gorm:"type:Roles" json:"role"`
	Status   Status `gorm:"type:Status; default:'pending'" json:"status,omitempty"`
}

type UserDetail struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status Status `json:"status"`
	Role   Roles  `json:"role"`
}

// hash password from user input
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// password and hashed check match
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("password incorrect")
	}
	return nil
}

// hash user password before save
func (u *User) BeforeSave() error {
	password := strings.TrimSpace(u.Password)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// strips any whitespace
func (u *User) Prepare() {
	u.Email = strings.TrimSpace(u.Email)
	u.Username = strings.TrimSpace(u.Username)
	u.Password = strings.TrimSpace(u.Password)
}

// validate user input
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Username == "" {
			return errors.New("username is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		return nil
	default:
		if u.Username == "" {
			return errors.New("username is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
		if u.Email == "" {
			return errors.New("email is required")
		}
		if u.Name == "" {
			return errors.New("name is required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) GetUser(db *gorm.DB) (*User, error) {
	account := &User{}
	if err := db.Debug().Table("users").Where("username = ?", u.Username).First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func GetAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	if err := db.Debug().Table("users").Find(&users).Error; err != nil {
		return &[]User{}, err
	}
	return &users, nil

}

func GetUserById(id int, db *gorm.DB) (*UserDetail, error) {
	user := &UserDetail{}
	if err := db.Debug().Table("users").Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) VerifyFundraiser(id int, db *gorm.DB) (*User, error) {
	if err := db.Debug().Table("users").Where("id = ?", id).Updates(User{
		Status: "verified",
	}).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateUser(id int, db *gorm.DB) (*User, error) {
	password := strings.TrimSpace(u.Password)
	hashedPassword, _ := HashPassword(password)
	u.Password = string(hashedPassword)
	if err := db.Debug().Table("users").Where("id = ?", id).Updates(User{
		Email:    u.Email,
		Name:     u.Name,
		Password: u.Password,
	}).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}

func GetUnverifiedUser(db *gorm.DB) (*[]User, error) {
	users := []User{}
	if err := db.Debug().Table("users").Where("status = ?", "pending").Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}
