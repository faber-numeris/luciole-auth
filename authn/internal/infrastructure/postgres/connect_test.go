package postgres

import (
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
	"github.com/faber-numeris/luciole-auth/authn/internal/mocks"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	mockCfg := mocks.NewMockIAppConfig(t)
	mockCfg.EXPECT().DBHost().Return("localhost").Maybe()
	mockCfg.EXPECT().DBPort().Return(5432).Maybe()
	mockCfg.EXPECT().DBUser().Return("user").Maybe()
	mockCfg.EXPECT().DBPassword().Return("pass").Maybe()
	mockCfg.EXPECT().DBName().Return("db").Maybe()
	mockCfg.EXPECT().DBSSLMode().Return("disable").Maybe()

	t.Run("success", func(t *testing.T) {
		db_mock, _, _ := sqlmock.New()
		dummyDB := sqlx.NewDb(db_mock, "postgres")
		
		patches := gomonkey.ApplyFunc(config.LoadConfig, func() (config.IAppConfig, error) {
			return mockCfg, nil
		})
		defer patches.Reset()

		patches.ApplyFunc(sqlx.Connect, func(driverName, dataSourceName string) (*sqlx.DB, error) {
			return dummyDB, nil
		})

		db := Connect()
		assert.Equal(t, dummyDB, db)
		// Reset singleton for other tests
		once = sync.Once{}
		DBInstance = nil
	})
}
