package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER,
			status TEXT,
			address TEXT,
			created_at TEXT
		)
	`)
	require.NoError(t, err)
	return db
}

func TestAddGetDelete(t *testing.T) {
	db := setupDB(t)
	store := NewParcelStore(db)
	parcel := Parcel{
		Client:    1,
		Status:    ParcelStatusRegistered,
		Address:   "123 Test Street",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// Get
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, storedParcel.Client)

	// Delete
	err = store.Delete(id)
	require.NoError(t, err)

	// Verify Deletion
	_, err = store.Get(id)
	require.Error(t, err)
}
