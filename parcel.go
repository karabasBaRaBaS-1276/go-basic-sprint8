package main

import (
	"database/sql"
	"fmt"
)

const (
	addParcelQuery            = `INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :createdAt)`
	getParcelByNumberQuery    = `SELECT client, status, address, created_at FROM parcel WHERE number = :number`
	getParcelByClientQuery    = `SELECT client, status, address, created_at FROM parcel WHERE client = :client`
	setAddressQuery           = `UPDATE parcel SET address = :newAddress WHERE number = :number AND status = :statusRegistered`
	setStatusQuery            = `UPDATE parcel SET status = :newStatus WHERE number = :number`
	deleteParcelByNumberQuery = `DELETE from parcel WHERE number = :number AND status = :statusRegistered`
)

// Информация о хранилище посылок
type ParcelStore struct {
	db *sql.DB
}

// Инициализация нового хранилища посылок
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Добавление новой записи о посылке в хранилище
// Принимает на вход структуру Parcel
// Возвращает:
//   - идентификатор новой записи в хранилище (int)
//   - ошибка, если что-то пошло не так. (error)
func (s ParcelStore) Add(p Parcel) (int, error) {
	// добавление строки в таблицу parcel, используйте данные из переменной p
	result, err := s.db.Exec(
		addParcelQuery,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdAt", p.CreatedAt),
	)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении записи: %w", err)
	}

	number, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении записи: %w", err)
	}

	return int(number), nil
}

// Получение информации о посылке по её номеру
// Возвращает:
//   - структуру (Parcel)
//   - ошибка, если что-то пошло не так. (error)
func (s ParcelStore) Get(number int) (Parcel, error) {

	if number == 0 {
		return Parcel{}, fmt.Errorf("ошибка при чтении записи: недопустимое значение id посылки (%d)", number)
	}

	row := s.db.QueryRow(
		getParcelByNumberQuery,
		sql.Named("number", number),
	)

	var (
		client    int
		status    string
		address   string
		createdAt string
	)

	err := row.Scan(&client, &status, &address, &createdAt)
	if err != nil {
		return Parcel{}, fmt.Errorf("ошибка при чтении записи: %w", err)
	}

	p := Parcel{
		Client:    client,
		Status:    status,
		Address:   address,
		CreatedAt: createdAt,
	}

	return p, nil
}

// Получение списка посылок по клиенту
// Возвращает:
//   - массив структур (Parcel)
//   - ошибка, если что-то пошло не так. (error)
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	if client == 0 {
		return nil, fmt.Errorf("ошибка при поиске списка посылок по клиенту: недопустимое значение клиента (%d)", client)
	}

	rows, err := s.db.Query(
		getParcelByClientQuery,
		sql.Named("client", client),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при поиске списка посылок по клиенту: %w", err)
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		var (
			client    int
			status    string
			address   string
			createdAt string
		)

		err := rows.Scan(&client, &status, &address, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении записи: %w", err)
		}

		res = append(
			res,
			Parcel{
				Client:    client,
				Status:    status,
				Address:   address,
				CreatedAt: createdAt,
			},
		)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при поиске списка посылок по клиенту: %w", err)
	}

	return res, nil
}

// Изменение статуса посылки
// Принимает на вход
//   - number идентификатор посылки
//   - status новый статус посылки
//
// Возвращает:
//   - ошибка, если что-то пошло не так. (error)
func (s ParcelStore) SetStatus(number int, status string) error {

	if status == "" {
		return fmt.Errorf("ошибка при изменении статуса: статус не может быть пустым")
	}
	if number == 0 {
		return fmt.Errorf("ошибка при изменении статуса: недопустимое значение id посылки (%d)", number)
	}

	result, err := s.db.Exec(
		setStatusQuery,
		sql.Named("newStatus", status),
		sql.Named("number", number),
	)
	if err != nil {
		return fmt.Errorf("ошибка при изменении статуса: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при изменении статуса: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("ошибка при изменении статуса: ожидалось обновить 1 запись, но обновлено %d", rowsAffected)
	}

	return nil
}

// Изменение адреса доставки посылки. Только для посылок в статусе 'registered'
// Принимает на вход
//   - number идентификатор посылки
//   - address новый адрес доставки
//
// Возвращает:
//   - ошибка, если что-то пошло не так. (error)
func (s ParcelStore) SetAddress(number int, address string) error {

	if address == "" {
		return fmt.Errorf("ошибка при изменении адреса: адрес не может быть пустым")
	}
	if number == 0 {
		return fmt.Errorf("ошибка при изменении адреса: недопустимое значение id посылки (%d)", number)
	}

	result, err := s.db.Exec(
		setAddressQuery,
		sql.Named("newAddress", address),
		sql.Named("number", number),
		sql.Named("statusRegistered", ParcelStatusRegistered),
	)

	if err != nil {
		return fmt.Errorf("ошибка при изменении адреса: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при изменении адреса: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("ошибка при изменении адреса: ожидалось обновить 1 запись, но обновлено %d", rowsAffected)
	}

	return nil
}

// Удаление записи о посылке в статусе registered по её номеру
// Возвращает:
//   - ошибка, если что-то пошло не так. (error)
func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	if number == 0 {
		return fmt.Errorf("ошибка при удалении записи: недопустимое значение number посылки (%d)", number)
	}

	result, err := s.db.Exec(
		deleteParcelByNumberQuery,
		sql.Named("number", number),
		sql.Named("statusRegistered", ParcelStatusRegistered),
	)

	if err != nil {
		return fmt.Errorf("ошибка при удалении записи: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при удалении записи: %w", err)
	}

	if rowsAffected != 1 {
		fmt.Printf("ошибка при удалении записи: ожидалось удаление 1 записи, но удалено %d\n", rowsAffected)
	}

	return nil
}
