package repositories

import (
	"errors"

	"github.com/ad3n/golang-testable/models"

	"gorm.io/gorm"
)

type Account struct {
	Storage *gorm.DB
}

func (r *Account) Find(Id int) (*models.Account, error) {
	account := models.Account{}
	err := r.Storage.First(&account, "account_number = ?", Id).Error

	if account.ID == 0 {
		return &account, errors.New("account not found")
	}

	return &account, err
}

func (r *Account) Saves(accounts ...*models.Account) error {
	tx := r.Storage.Begin()
	defer func() {
		r := recover()
		if r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Error
	if err != nil {
		return err
	}

	for _, m := range accounts {
		err = tx.Save(m).Error
		if err != nil {
			tx.Rollback()

			return err
		}
	}

	return tx.Commit().Error
}
