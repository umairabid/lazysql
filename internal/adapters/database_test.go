package adapters

import (
	"testing"
)

func TestDbConnectionString(t *testing.T) {
	c := DbConnection{
		Username: "postgres",
		Password: "secret",
		Host:     "localhost",
		Port:     "5432",
	}
	got := c.String()
	want := "user=postgres password=secret host=localhost port=5432"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestValidateConnection(t *testing.T) {
	tests := []struct {
		name    string
		conn    DbConnection
		wantErr bool
	}{
		{
			name: "valid connection",
			conn: DbConnection{
				Host:     "localhost",
				Port:     "5432",
				Username: "postgres",
				Password: "postgres",
				Driver:   "pgx",
			},
			wantErr: false,
		},
		{
			name:    "missing all fields",
			conn:    DbConnection{},
			wantErr: true,
		},
		{
			name: "missing host",
			conn: DbConnection{
				Port:     "5432",
				Username: "postgres",
				Password: "postgres",
				Driver:   "pgx",
			},
			wantErr: true,
		},
		{
			name: "missing port",
			conn: DbConnection{
				Host:     "localhost",
				Username: "postgres",
				Password: "postgres",
				Driver:   "pgx",
			},
			wantErr: true,
		},
		{
			name: "missing username",
			conn: DbConnection{
				Host:     "localhost",
				Port:     "5432",
				Password: "postgres",
				Driver:   "pgx",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			conn: DbConnection{
				Host:     "localhost",
				Port:     "5432",
				Username: "postgres",
				Driver:   "pgx",
			},
			wantErr: true,
		},
		{
			name: "missing driver",
			conn: DbConnection{
				Host:     "localhost",
				Port:     "5432",
				Username: "postgres",
				Password: "postgres",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.conn.validateConnection()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnectWithDatabase(t *testing.T) {
	t.Run("successful connection", func(t *testing.T) {
		c := DbConnection{
			Host:     "localhost",
			Port:     "5432",
			Username: "postgres",
			Password: "postgres",
			Driver:   "pgx",
		}
		db, err := c.InitConnection()
		if err != nil {
			t.Fatalf("InitConnection() unexpected error: %v", err)
		}
		if db == nil {
			t.Error("InitConnection() returned nil db")
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		c := DbConnection{
			Host:     "localhost",
			Port:     "5432",
			Username: "postgres",
			Password: "wrongpassword",
			Driver:   "pgx",
		}
		_, err := c.InitConnection()
		if err == nil {
			t.Error("InitConnection() expected error for invalid credentials, got nil")
		}
	})

	t.Run("validation failure returns error", func(t *testing.T) {
		c := DbConnection{
			Host:   "localhost",
			Port:   "5432",
			Driver: "pgx",
			// missing Username and Password
		}
		_, err := c.InitConnection()
		if err == nil {
			t.Error("InitConnection() expected validation error, got nil")
		}
	})
}
