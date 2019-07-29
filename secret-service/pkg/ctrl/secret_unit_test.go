package ctrl

import (
	"strings"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	reqBody := `{"secret":"mysec1", "expireAfterViews":10, "expireAfter":5}`

	secret, err := decode(strings.NewReader(reqBody))

	assert.NoError(t, err)
	assert.NotNil(t, secret)
	assert.Equal(t, "mysec1", secret.Secret)
	assert.Equal(t, 10, secret.ExpireAfterViews)
	assert.Equal(t, 5, secret.ExpireAfter)
}

func TestEncode(t *testing.T) {
	now := time.Now()
	secret := &secretResp{
		Hash:           uuid.NewV4().String(),
		SecretText:     "mysec1",
		CreatedAt:      now.Format(time.RFC3339),
		ExpiresAt:      now.Add(time.Minute * 5).Format(time.RFC3339),
		RemainingViews: 10,
	}

	content, err := encode(secret)

	assert.NoError(t, err)
	assert.NotNil(t, content)
	assert.Contains(t, string(content), secret.Hash)
	assert.Contains(t, string(content), secret.SecretText)
	assert.Contains(t, string(content), secret.CreatedAt)
	assert.Contains(t, string(content), secret.ExpiresAt)
}
