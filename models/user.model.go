package models

import (
	"database/sql"
	"time"

	"github.com/seanburman/seanburman.com/db"
	"github.com/seanburman/seanburman.com/utils"
	"gorm.io/gorm"
)

type UserError string

func (e UserError) Error() string {
	return string(e)
}

const (
	ErrPasswordInvalid            UserError = "password must contain at least 8 characters, one capital letter, one number, and one special character"
	ErrPasswordMissingLetter      UserError = "password must contain at least one letter"
	ErrPasswordMissingCapital     UserError = "password must contain at least one capital letter"
	ErrPasswordMissingNumber      UserError = "password must contain at least one number"
	ErrPasswordMissingSpecialChar UserError = "password must contain at least one special character"
	ErrPasswordsDoNotMatch        UserError = "passwords do not match"
	ErrUsernameExists             UserError = "username already exists"
	ErrEmailExists                UserError = "email already exists"
	ErrUserNotFound               UserError = "user not found"
	ErrWrongPassword              UserError = "invalid username or password"
)

type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Verified  bool      `json:"verified"`
}

func NewUser(db *db.Database, username string, password string, email string) (*User, error) {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
	}, nil
}

func (u *User) Create(db *db.Database) (affected int64, err error) {
	var user User = User{}
	result := db.Postgres.Where(
		"username = @user OR email = @email",
		sql.Named("user", u.Username),
		sql.Named("email", u.Email),
	).Find(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	if user.Username != "" && user.Username == u.Username {
		return 0, ErrUsernameExists
	}

	if user.Email != "" {
		return 0, ErrEmailExists
	}

	result = db.Postgres.Create(u)
	return result.RowsAffected, result.Error
}

func (u *User) Get(db *db.Database) (*User, error) {
	var user User
	result := db.Postgres.Where(User{Username: u.Username}).Find(&user)
	if user == (User{}) {
		return nil, ErrUserNotFound
	}
	return &user, result.Error
}

func (u *User) Update(db *db.Database) error {
	return db.Postgres.Save(u).Error
}

func (u *User) Delete(db *db.Database, v any) error {
	return db.Postgres.Delete(v).Error
}

func (u *User) Verify(db *db.Database) error {
	u.Verified = true
	return u.Update(db)
}

func (u *User) Authenticate(password string) bool {
	return utils.ValidatePassword(u.Password, password)
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	u.CreatedAt = time.Now()
	err := u.Validate()
	if err != nil {
		return err
	}
	hash, err := utils.HashString(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash
	return nil
}

func (u *User) BeforeUpdate(db *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return u.Validate()
}

func (u *User) Validate() error {
	if len(u.Password) < 8 {
		return ErrPasswordInvalid
	}
	var hasLetter, hasCapital, hasNumber, hasSpecialChar bool
	for _, r := range u.Password {
		if r >= 'A' && r <= 'Z' {
			hasCapital = true
		}
		if r >= 'a' && r <= 'z' {
			hasLetter = true
		}
		if r >= '0' && r <= '9' {
			hasNumber = true
		}
		if r >= 33 && r <= 47 {
			hasSpecialChar = true
		}
	}
	if !hasLetter && !hasCapital && !hasNumber && !hasSpecialChar {
		return ErrPasswordInvalid
	}
	return nil
}
