package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
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
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД,
	number, err := store.Add(parcel)
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	parcel.Number = number
	// и наличии идентификатора
	require.NotEmpty(t, number)

	// get
	// получите только что добавленную посылку,
	storedParcel, err := store.Get(number)
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	require.Equal(t, storedParcel, parcel)

	// delete
	// удалите добавленную посылку,
	err = store.Delete(number)
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(number)
	require.Equal(t, sql.ErrNoRows, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД,
	number, err := store.Add(parcel)
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	parcel.Number = number
	// и наличии идентификатора
	require.NotEmpty(t, number)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	storedParcel, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, newAddress, storedParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД,
	number, err := store.Add(parcel)
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	parcel.Number = number
	// и наличии идентификатора
	require.NotEmpty(t, number)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(number, ParcelStatusSent)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	storedParsel, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, storedParsel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	// срез посылок
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	// пустая мапа под посылки
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	// все посылки от одного (случайного) клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		id, err := store.Add(parcels[i])
		require.NoError(t, err)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	storedParcels, err := store.GetByClient(client)
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
