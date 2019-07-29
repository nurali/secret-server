package model

import (
	"time"

	errs "github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	Create(secret *Secret) (*Secret, error)
	Load(hash uuid.UUID) (*Secret, error)
}

type GormRepository struct {
	DB *gorm.DB
}

type Secret struct {
	Hash           uuid.UUID `sql:"type:uuid default uuid_generate_v4()" gorm:"primary_key"`
	SecretText     string
	CreatedAt      time.Time
	ExpiresAt      time.Time
	RemainingViews int
}

func (r *GormRepository) Create(secret *Secret) (*Secret, error) {
	err := r.DB.Create(secret).Error
	if err != nil {
		log.Error("failed to insert secret")
		return nil, errs.WithStack(err)
	}
	return secret, nil
}

func (r *GormRepository) Load(hash uuid.UUID) (*Secret, error) {
	secret := Secret{}
	tx := r.DB.Model(&Secret{}).Where("hash = ?", hash).First(&secret)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &secret, nil
}
