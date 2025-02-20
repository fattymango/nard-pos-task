package db

import (
	"testing"

	"multitenant/pkg/config"
)

func TestNewMysqlDB(t *testing.T) {
	testCases := []struct {
		Name    string
		Config  *config.Config
		WantErr bool
	}{
		{
			Name: "OK",
			Config: &config.Config{
				DB: config.DB{
					Host:            "localhost",
					Port:            "3306",
					User:            "task",
					Password:        "task123",
					Name:            "task",
					SSLMode:         "disable",
					MaxIdleConns:    2,
					MaxOpenConns:    5,
					MaxConnLifetime: 10,
				},
			},
			WantErr: false,
		},
		{
			Name: "InvalidURL",
			Config: &config.Config{
				DB: config.DB{
					Host:     "invalid_host",
					Port:     "3306",
					User:     "task",
					Password: "task123",
					Name:     "mysql",
					SSLMode:  "disable",
				},
			},
			WantErr: true,
		},
		{
			Name: "InvalidCredentials",
			Config: &config.Config{
				DB: config.DB{
					Host:     "localhost",
					Port:     "3306",
					User:     "invalid",
					Password: "invalid",
					Name:     "mysql",
					SSLMode:  "disable",
				},
			},
			WantErr: true,
		},
		{
			Name: "InvalidDBName",
			Config: &config.Config{
				DB: config.DB{
					Host:     "localhost",
					Port:     "3306",
					User:     "task",
					Password: "task123",
					Name:     "invalid",
					SSLMode:  "disable",
				},
			},
			WantErr: true,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			db, err := NewDB(testCases[i].Config)
			if testCases[i].WantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !testCases[i].WantErr && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if db != nil {
				db.Close()
			}
		})
	}
}
