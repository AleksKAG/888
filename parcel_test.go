package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	store := NewParcelStore("sqlite", "tracker.db")
	defer store.db.Close()

	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, number)

	addedParcel, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, addedParcel.Client, parcel.Client)
	assert.Equal(t, addedParcel.Status, parcel.Status)
	assert.Equal(t, addedParcel.Address, parcel.Address)
	assert.Equal(t, addedParcel.CreatedAt, parcel.CreatedAt)

	err = store.Delete(number)
	assert.NoError(t, err)

	_, err = store.Get(number)
	assert.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	store := NewParcelStore("sqlite", "tracker.db")
	defer store.db.Close()
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	addedParcel, _ := store.Get(id)
	assert.Equal(t, addedParcel.Address, newAddress)

}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	store := NewParcelStore("sqlite", "tracker.db")
	defer store.db.Close()

	parcel := getTestParcel()

	number, err := store.Add(parcel)
	assert.NoError(t, err)
	assert.NotEmpty(t, number)

	err = store.SetStatus(number, ParcelStatusSent)
	assert.NoError(t, err)

	err = store.SetStatus(number, ParcelStatusDelivered)
	assert.NoError(t, err)
	addedParcel, _ := store.Get(number)
	assert.Equal(t, addedParcel.Status, ParcelStatusDelivered)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	store := NewParcelStore("sqlite", "tracker.db")
	defer store.db.Close()

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		number, err := store.Add(parcels[i])
		assert.NoError(t, err)
		assert.NotEmpty(t, number)

		parcels[i].Number = number
		parcelMap[number] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)

	assert.NoError(t, err)
	assert.Equal(t, len(storedParcels), len(parcels))

	for _, parcel := range storedParcels {

		assert.Equal(t, parcel.Number, parcelMap[parcel.Number].Number)
		assert.Equal(t, parcel.Client, parcelMap[parcel.Number].Client)
		assert.Equal(t, parcel.Status, parcelMap[parcel.Number].Status)
		assert.Equal(t, parcel.Address, parcelMap[parcel.Number].Address)
		assert.Equal(t, parcel.CreatedAt, parcelMap[parcel.Number].CreatedAt)
	}
}
