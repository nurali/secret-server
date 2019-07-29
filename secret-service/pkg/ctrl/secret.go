package ctrl

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/nurali/secret-server/secret-service/pkg/model"
	uuid "github.com/satori/go.uuid"
)

type SecretCtrl interface {
	Create(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
}

type secretCtrl struct {
	secretStore map[string]*model.Secret
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

func NewSecretCtrl() SecretCtrl {
	return &secretCtrl{
		secretStore: make(map[string]*model.Secret),
	}
}

func ToSecret(secretIn *secretReq) *model.Secret {
	hash := uuid.NewV4().String()
	now := time.Now()
	secret := &model.Secret{
		Hash:           hash,
		SecretText:     secretIn.Secret,
		CreatedAt:      now,
		ExpiresAt:      now.Add(time.Minute * time.Duration(secretIn.ExpireAfter)),
		RemainingViews: secretIn.ExpireAfterViews,
	}
	return secret
}

func ToSecretResp(secret *model.Secret) *secretResp {
	res := &secretResp{
		Hash:           secret.Hash,
		SecretText:     secret.SecretText,
		CreatedAt:      secret.CreatedAt.Format(time.RFC3339),
		ExpiresAt:      secret.ExpiresAt.Format(time.RFC3339),
		RemainingViews: secret.RemainingViews,
	}
	return res
}

func (c *secretCtrl) Create(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Invalid secret", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	secretIn, err := decode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateSecretReq(secretIn); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	newSecret := ToSecret(secretIn)
	c.secretStore[newSecret.Hash] = newSecret

	secretOut := ToSecretResp(newSecret)
	content, err := encode(secretOut)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(content)
}

func (c *secretCtrl) Show(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	secret := c.secretStore[hash]
	if secret == nil {
		http.Error(w, "secret not found", http.StatusNotFound)
		return
	}

	secretOut := ToSecretResp(secret)
	content, err := encode(secretOut)
	if err != nil {
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
		return errors.New("Invalid expireAfterViews, it must be greater than 0")
	}
	return nil
}
