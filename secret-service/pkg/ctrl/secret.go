package ctrl

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/nurali/secret-server/secret-service/pkg/model"
)

type SecretCtrl interface {
	Create(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
}

type secretCtrl struct {
	store model.SecretStore
}

type secretReq struct {
	Secret           string `json:"secret"`
	ExpireAfterViews int    `json:"expireAfterViews"`
	ExpireAfter      int    `json:"expireAfter"`
}

type secretResp struct {
	Hash           string `json:"hash"`
	SecretText     string `json:"secretText"`
	CreatedAt      string `json:"createdAt"`
	ExpiresAt      string `json:"expiresAt"`
	RemainingViews int    `json:"remainingViews"`
}

func NewSecretCtrl(db *gorm.DB) SecretCtrl {
	return &secretCtrl{
		store: model.NewSecretDBStore(db),
	}
}

func ToSecret(secretIn *secretReq) *model.Secret {
	now := time.Now()
	secret := &model.Secret{
		SecretText:     secretIn.Secret,
		CreatedAt:      now,
		ExpiresAt:      now.Add(time.Minute * time.Duration(secretIn.ExpireAfter)),
		RemainingViews: secretIn.ExpireAfterViews,
	}
	return secret
}

func ToSecretResp(secret *model.Secret) *secretResp {
	res := &secretResp{
		Hash:           secret.Hash.String(),
		SecretText:     secret.SecretText,
		CreatedAt:      secret.CreatedAt.Format(time.RFC3339),
		ExpiresAt:      secret.ExpiresAt.Format(time.RFC3339),
		RemainingViews: secret.RemainingViews,
	}
	return res
}

func (c *secretCtrl) Create(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		log.Errorf("secret body is missing")
		http.Error(w, "Invalid secret", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	secretIn, err := decode(r.Body)
	if err != nil {
		log.Errorf("json to secret failed with error, %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateSecretReq(secretIn); err != nil {
		log.Errorf("secret validation failed with error, %v", err)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	newSecret := ToSecret(secretIn)
	newSecret, err = c.store.Create(newSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	secretOut := ToSecretResp(newSecret)
	content, err := encode(secretOut)
	if err != nil {
		log.Errorf("secret to json failed with error, %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(content)
}

func (c *secretCtrl) Show(w http.ResponseWriter, r *http.Request) {
	var hash uuid.UUID
	var err error
	if hash, err = uuid.FromString(mux.Vars(r)["hash"]); err != nil {
		log.Errorf("secret hash conversion failed with error, %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secret, err := c.store.ReadOnce(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer func() {
		if err != nil {
			log.Debugf("unread secret due to error, %v", err)
			_, unreadErr := c.store.UnreadOnce(hash)
			if unreadErr != nil {
				log.Error(unreadErr)
			}
		}
	}()

	secretOut := ToSecretResp(secret)
	content, err := encode(secretOut)
	if err != nil {
		log.Errorf("secret to json failed with error, %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(content)
}

func decode(content io.Reader) (*secretReq, error) {
	decoder := json.NewDecoder(content)
	var secret secretReq
	err := decoder.Decode(&secret)
	return &secret, err
}

func encode(secret *secretResp) ([]byte, error) {
	content, err := json.Marshal(secret)
	return content, err
}

func validateSecretReq(secret *secretReq) error {
	if secret == nil {
		return nil
	}
	if secret.ExpireAfterViews <= 0 {
		return errors.New("invalid expireAfterViews, it must be greater than 0")
	}
	return nil
}
