package main

import (
	"database/sql"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotZero(t, id)

	storedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, parcel.Client, storedParcel.Client)
	assert.Equal(t, parcel.Status, storedParcel.Status)
	assert.Equal(t, parcel.Address, storedParcel.Address)

	err = store.Delete(id)
	assert.NoError(t, err)

	_, err = store.Get(id)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	assert.NoError(t, err)

	updatedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, newAddress, updatedParcel.Address)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	assert.NoError(t, err)

	err = store.SetStatus(id, ParcelStatusSent)
	assert.NoError(t, err)

	updatedParcel, err := store.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, updatedParcel.Status)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	for _, parcel := range parcels {
		id, err := store.Add(parcel)
		assert.NoError(t, err)
		parcel.Number = id
	}

	storedParcels, err := store.GetByClient(client)
	assert.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))

	for i, storedParcel := range storedParcels {
		assert.Equal(t, parcels[i].Client, storedParcel.Client)
		assert.Equal(t, parcels[i].Address, storedParcel.Address)
		assert.Equal(t, parcels[i].Status, storedParcel.Status)
	}
}
