package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {

	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES ( :client, :status, :address, :createdAt)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdAt", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	p := Parcel{}

	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	var res []Parcel

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))

	if err != nil {

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

		if err != nil {
			return nil, err
		}

		res = append(res, p)

	}

	if err := rows.Err(); err != nil {

		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {

	p := Parcel{}

	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))

	err := row.Scan(&p.Status)

	if err != nil {
		return err
	}

	_, err = s.db.Exec("UPDATE parcel SET address = CASE WHEN status = :status THEN :address ELSE address END WHERE number = :number",
		sql.Named("status", ParcelStatusRegistered),
		sql.Named("address", address),
		sql.Named("number", number))

	return err

}

func (s ParcelStore) Delete(number int) error {

	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status == :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	if err != nil {
		return fmt.Errorf("db exec error: %w", err)

	}

	return nil
}
