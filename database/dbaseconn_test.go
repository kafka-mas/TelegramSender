package database_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kafka-mas/TelegramSender/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	os.Create(dbPath)
	os.Chmod(dbPath, 0444)
	db := database.NewSQLite(dbPath)
	err := db.Create()
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	dbase, dbpath := setupDB(t)

	tests := []struct {
		name      string
		dbFactory func(t *testing.T) database.DB
		wantErr   bool
	}{
		{
			name:      "Success delete",
			dbFactory: func(t *testing.T) database.DB { return dbase },
			wantErr:   false,
		},
		{
			name: "Error delete",
			dbFactory: func(t *testing.T) database.DB {
				os.Chmod(dbpath, 0444)
				return dbase
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.dbFactory(t)
			err := db.Delete()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddRecord(t *testing.T) {
	goodDB, goodDBpath := setupDB(t)

	tests := []struct {
		name      string
		dbFactory func(t *testing.T) database.DB
		userID    int64
		fname     string
		wantID    int
		wantErr   bool
	}{
		{
			name:      "Succes test  user 1",
			dbFactory: func(t *testing.T) database.DB { return goodDB },
			userID:    12345678901,
			fname:     "Alice",
			wantID:    1,
			wantErr:   false,
		},
		{
			name:      "Succes test  user 2",
			dbFactory: func(t *testing.T) database.DB { return goodDB },
			userID:    12345678902,
			fname:     "Alice",
			wantID:    2,
			wantErr:   false,
		},
		{
			name:      "Error test  not unique user_id",
			dbFactory: func(t *testing.T) database.DB { return goodDB },
			userID:    12345678901,
			fname:     "Alice",
			wantID:    -1,
			wantErr:   true,
		},
		{
			name: "Error test  bad database path",
			dbFactory: func(t *testing.T) database.DB {
				badPath := filepath.Join(t.TempDir(), "missing", "test.db")
				return database.NewSQLite(badPath)
			},
			userID:  999,
			fname:   "Fail",
			wantID:  -1,
			wantErr: true,
		},
		{
			name: "Error test  bad db.Exec",
			dbFactory: func(t *testing.T) database.DB {
				os.Chmod(goodDBpath, 0444)
				return goodDB
			},
			userID:  998,
			fname:   "Fail",
			wantID:  -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.dbFactory(t)
			gotID, err := db.AddRecord(tt.userID, tt.fname)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, -1, gotID)
			} else {
				assert.NoError(t, err)
				assert.Greater(t, gotID, 0)
			}
		})
	}
}

func TestGetRecord(t *testing.T) {
	db, _ := setupDB(t)
	tests := []struct {
		name      string
		userID    int64
		dbFactory func(t *testing.T) database.DB
		fname     string
		wantErr   bool
	}{
		{
			name:      "Error test",
			userID:    12345678901,
			dbFactory: func(t *testing.T) database.DB { return db },
			fname:     "Alice",
			wantErr:   true,
		},
		{
			name:   "Succes test",
			userID: 12345678902,
			dbFactory: func(t *testing.T) database.DB {
				db.AddRecord(12345678902, "Alice")
				return db
			},
			fname:   "Alice",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.dbFactory(t)
			result, err := db.GetRecord(1)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, result.UserID, "Must be equal")
				assert.Equal(t, tt.fname, result.Fname, "Must be equal")
			}

		})
	}
}

func TestGetAllRecords(t *testing.T) {
	db, _ := setupDB(t)
	tests := []struct {
		name      string
		dbFactory func(t *testing.T) database.DB
		user_s    []database.DBUser
		wantErr   bool
	}{
		{
			name: "Success test",
			dbFactory: func(t *testing.T) database.DB {
				return db
			},
			user_s:  []database.DBUser{},
			wantErr: false,
		},
		{
			name: "Succes test",
			dbFactory: func(t *testing.T) database.DB {
				db.AddRecord(12345678901, "UserA")
				db.AddRecord(12345678902, "UserB")
				db.AddRecord(12345678903, "UserC")
				return db
			},
			user_s: []database.DBUser{
				{
					ID:        1,
					UserID:    12345678901,
					Fname:     "UserA",
					ISdefault: false,
				},
				{
					ID:        2,
					UserID:    12345678902,
					Fname:     "UserB",
					ISdefault: false,
				},
				{
					ID:        3,
					UserID:    12345678903,
					Fname:     "UserC",
					ISdefault: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		db := tt.dbFactory(t)
		result, err := db.GetAllRecords()
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.user_s, *result)
		}
	}
}

func TestRemoveRecord(t *testing.T) {
	db, _ := setupDB(t)
	tests := []struct {
		name      string
		dbID      int
		dbFactory func(t *testing.T) database.DB
		wantErr   bool
	}{
		{
			name: "Test empty DB",
			dbID: 1,
			dbFactory: func(t *testing.T) database.DB {
				return db
			},
			wantErr: false,
		},
		{
			name: "Succes test",
			dbID: 1,
			dbFactory: func(t *testing.T) database.DB {
				db.AddRecord(123, "A")
				return db
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.dbFactory(t)
			err := db.RemoveRecord(tt.dbID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestSetDefault(t *testing.T) {
	db, dbpath := setupDB(t)
	tests := []struct {
		name        string
		defaultDbID int
		dbFactory   func(t *testing.T) database.DB
		wantErr     bool
	}{
		{
			name:        "Success test",
			defaultDbID: 1,
			dbFactory: func(t *testing.T) database.DB {
				db.AddRecord(123, "Alice")
				db.AddRecord(124, "Alice")
				return db
			},
			wantErr: false,
		},
		{
			name:        "Success test - empty DB",
			defaultDbID: 1,
			dbFactory:   func(t *testing.T) database.DB { return db },
			wantErr:     false,
		},
		{
			name:        "Error test",
			defaultDbID: 1,
			dbFactory: func(t *testing.T) database.DB {
				db.AddRecord(123, "Alice")
				db.AddRecord(124, "Alice")
				os.Chmod(dbpath, 0444)
				return db
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		db := tt.dbFactory(t)
		err := db.SetDefault(tt.defaultDbID)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			users, _ := db.GetAllRecords()
			flag, count := false, 0
			for _, u := range *users {
				if u.ISdefault == true {
					flag = true
					count++
				}
			}
			assert.Equal(t, 1, count)
			assert.Equal(t, true, flag)
		}
	}
}

func TestRemoveDefault(t *testing.T) {
	db, dbpath := setupDB(t)
	tests := []struct {
		name      string
		dbFactory func(t *testing.T) database.DB
		wantErr   bool
	}{
		{
			name: "Success test",
			dbFactory: func(t *testing.T) database.DB {
				db.AddRecord(123, "Alice")
				db.AddRecord(124, "Alice")
				return db
			},
			wantErr: false,
		},
		{
			name:      "Success test - empty DB",
			dbFactory: func(t *testing.T) database.DB { return db },
			wantErr:   false,
		},
		{
			name: "Error test",
			dbFactory: func(t *testing.T) database.DB {
				db.AddRecord(123, "Alice")
				db.AddRecord(124, "Alice")
				db.SetDefault(1)
				os.Chmod(dbpath, 0444)
				return db
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		db := tt.dbFactory(t)
		err := db.RemoveDefault()
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			users, _ := db.GetAllRecords()
			flag, count := false, 0
			for _, u := range *users {
				if u.ISdefault == true {
					flag = true
					count++
				}
			}
			assert.Equal(t, 0, count)
			assert.Equal(t, false, flag)
		}
	}
}

func setupDB(t *testing.T) (database.DB, string) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	db := database.NewSQLite(dbPath)
	err := db.Create()
	require.NoError(t, err)
	return db, dbPath
}
