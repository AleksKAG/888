func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// Get
	storedParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, storedParcel.Client)
	require.Equal(t, parcel.Status, storedParcel.Status)
	require.Equal(t, parcel.Address, storedParcel.Address)

	// Delete
	err = store.Delete(id)
	require.NoError(t, err)

	// Verify deletion
	_, err = store.Get(id)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	// аналогично TestAddGetDelete, только с изменением адреса
}

func TestSetStatus(t *testing.T) {
	// аналогично TestAddGetDelete, только с изменением статуса
}

func TestGetByClient(t *testing.T) {
	// аналогично TestAddGetDelete, только с добавлением нескольких посылок одного клиента
}
