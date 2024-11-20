package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    randRange.Intn(1000),
		Status:    ParcelStatusRegistered,
		Address:   "Test Address",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func setupDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "test_tracker.db")
	require.NoError(t, err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS parcel (
        number INTEGER PRIMARY KEY AUTOINCREMENT,
        client INTEGER,
        status TEXT,
        address TEXT,
        created_at TEXT
    )`)
	require.NoError(t, err)

	return db
}

func TestAddGetDelete(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	fetchedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, fetchedParcel.Client)
	require.Equal(t, parcel.Status, fetchedParcel.Status)
	require.Equal(t, parcel.Address, fetchedParcel.Address)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)

	newAddress := "New Test Address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	fetchedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, fetchedParcel.Address)
}

func TestSetStatus(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	fetchedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, fetchedParcel.Status)
}

func TestGetByClient(t *testing.T) {
	db := setupDatabase(t)
	defer db.Close()

	store := NewParcelStore(db)

	clientID := randRange.Intn(1000)
	parcels := []Parcel{
		{Client: clientID, Status: ParcelStatusRegistered, Address: "Address 1", CreatedAt: time.Now().UTC().Format(time.RFC3339)},
		{Client: clientID, Status: ParcelStatusRegistered, Address: "Address 2", CreatedAt: time.Now().UTC().Format(time.RFC3339)},
	}

	for _, p := range parcels {
		_, err := store.Add(p)
		require.NoError(t, err)
	}

	fetchedParcels, err := store.GetByClient(clientID)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(fetchedParcels))
}
