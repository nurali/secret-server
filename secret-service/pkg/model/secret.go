package model

import (
	"sync"
	"time"

	errs "github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type SecretStore interface {
	Create(secret *Secret) (*Secret, error)
	Save(secret *Secret) (*Secret, error)
	ReadOnce(hash uuid.UUID) (*Secret, error)
	UnreadOnce(hash uuid.UUID) (*Secret, error)
}

type secretDBStore struct {
	db        *gorm.DB
	viewMutex *sync.Mutex
}

type Secret struct {
	Hash           uuid.UUID `sql:"type:uuid default uuid_generate_v4()" gorm:"primary_key"`
	SecretText     string
	CreatedAt      time.Time
	ExpiresAt      time.Time
	RemainingViews int
}

func NewSecretDBStore(db *gorm.DB) SecretStore {
	return &secretDBStore{
		db:        db,
		viewMutex: &sync.Mutex{},
	}
}

func (r *secretDBStore) Create(secret *Secret) (*Secret, error) {
	err := r.db.Create(secret).Error
	if err != nil {
		log.Errorf("create secret failed with error, %v", err)
		return nil, errs.WithStack(err)
	}
	return secret, nil
}

func (r *secretDBStore) load(hash uuid.UUID) (*Secret, error) {
	secret := Secret{}
	err := r.db.Model(&Secret{}).Where("hash = ?", hash).First(&secret).Error
	if err != nil {
		log.Errorf("load secret with hash '%v' failed with error, %v", hash, err)
		return nil, err
	}
	return &secret, nil
}

func (r *secretDBStore) Save(secret *Secret) (*Secret, error) {
	err := r.db.Model(&Secret{}).Where("hash = ?", secret.Hash).Save(secret).Error
	if err != nil {
		log.Errorf("save secret failed with error, %v", err)
		return nil, err
	}
	return secret, nil
}

func (r *secretDBStore) ReadOnce(hash uuid.UUID) (*Secret, error) {
	r.viewMutex.Lock()
	defer r.viewMutex.Unlock()

	secret, err := r.load(hash)
	if err != nil {
		return nil, err
	}

	if secret.RemainingViews > 0 {
		secret.RemainingViews = secret.RemainingViews - 1
	} else {
		return nil, errs.New("no more view allowed for the secret as it had already viewed for max allowed views")
	}
	return r.Save(secret)
}

func (r *secretDBStore) UnreadOnce(hash uuid.UUID) (*Secret, error) {
	r.viewMutex.Lock()
	defer r.viewMutex.Unlock()

	secret, err := r.load(hash)
	if err != nil {
		return nil, err
	}

	secret.RemainingViews = secret.RemainingViews + 1
	return r.Save(secret)
}
