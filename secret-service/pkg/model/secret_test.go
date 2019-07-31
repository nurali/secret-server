package model_test

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/nurali/secret-server/secret-service/pkg/app"
	"github.com/nurali/secret-server/secret-service/pkg/config"
	"github.com/nurali/secret-server/secret-service/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SecretStoreSuite struct {
	suite.Suite
	store model.SecretStore
}

func TestSecretStore(t *testing.T) {
	cfg := config.New()

	db, err := app.OpenDB(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := app.SetupDB(db); err != nil {
		t.Fatal(err)
	}

	s := &SecretStoreSuite{
		store: model.NewSecretDBStore(db),
	}
	suite.Run(t, s)
}

func (s *SecretStoreSuite) TestReadOnce() {

	s.T().Run("ok", func(t *testing.T) {
		want := newSecret(5)
		want, err := s.store.Create(want)
		require.NoError(t, err)

		// test
		got, err := s.store.ReadOnce(want.Hash)

		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, want.RemainingViews-1, got.RemainingViews)
	})

	s.T().Run("check_remaining_views", func(t *testing.T) {
		want := newSecret(5)
		want, err := s.store.Create(want)
		require.NoError(t, err)

		for i := 1; i <= 5; i++ {
			// test
			got, err := s.store.ReadOnce(want.Hash)
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, want.RemainingViews-i, got.RemainingViews)
		}
	})

}

func (s *SecretStoreSuite) TestReadOnceParallel() {
	// create
	remainingViews := 400
	want := newSecret(remainingViews)
	want, err := s.store.Create(want)
	require.NoError(s.T(), err)

	// this sub-test is needed to make sure "validate" executes after all worker test finishes
	s.T().Run("ok", func(t *testing.T) {
		workers := 4
		calls := remainingViews / workers
		for w := 1; w <= workers; w++ {
			name := fmt.Sprintf("worker-%d", w)
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				for i := 1; i <= calls; i++ {
					got, err := s.store.ReadOnce(want.Hash)
					assert.NoError(t, err)
					assert.NotNil(t, got)
				}
			})
		}
	})

	// validate
	got, err := s.store.ReadOnce(want.Hash)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), got)
}

func newSecret(remainingViews int) *model.Secret {
	now := time.Now()
	return &model.Secret{
		SecretText:     "mytext1",
		CreatedAt:      now,
		ExpiresAt:      now.Add(5 * time.Minute),
		RemainingViews: remainingViews,
	}
}
