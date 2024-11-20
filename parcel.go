package main

import (
	"database/sql"
)

type ParcelStore interface {
	Add(parcel Parcel) (int, error)
	GetByClient(client int) ([]Parcel, error)
	Get(number int) (Parcel, error)
	SetStatus(number int, status string) error
	SetAddress(number int, address string) error
	Delete(number int) error
}

type SQLiteParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) *SQLiteParcelStore {
	return &SQLiteParcelStore{db: db}
}

func (s *SQLiteParcelStore) Add(parcel Parcel) (int, error) {
	stmt, err := s.db.Prepare("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(parcel.Client, parcel.Status, parcel.Address, parcel.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (s *SQLiteParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var parcel Parcel
		if err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt); err != nil {
			return nil, err
		}
		parcels = append(parcels, parcel)
	}
	return parcels, nil
}

func (s *SQLiteParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = ?", number)
	var parcel Parcel
	if err := row.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt); err != nil {
		return Parcel{}, err
	}
	return parcel, nil
}

func (s *SQLiteParcelStore) SetStatus(number int, status string) error {
	stmt, err := s.db.Prepare("UPDATE parcel SET status = ? WHERE number = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(status, number)
	return err
}

func (s *SQLiteParcelStore) SetAddress(number int, address string) error {
	stmt, err := s.db.Prepare("UPDATE parcel SET address = ? WHERE number = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(address, number)
	return err
}

func (s *SQLiteParcelStore) Delete(number int) error {
	stmt, err := s.db.Prepare("DELETE FROM parcel WHERE number = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(number)
	return err
}
