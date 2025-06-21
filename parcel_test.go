package main

import (
	"database/sql"
	"log"
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

func add(t *testing.T, store ParcelStore, parcel Parcel) int {
	number, err := store.Add(parcel)
	require.NoError(t, err, "Нет ошибок при регистрации посылки")
	require.NotEmpty(t, number, "Новой посылке присвоен номер")
	return number
}

func delete(t *testing.T, store ParcelStore, number int) {
	err := store.Delete(number)
	assert.NoError(t, err, "Нет ошибок при удалении информаци о посылке")

	// посылку больше нельзя получить из БД
	_, err = store.Get(number)
	assert.ErrorIsf(t, err, sql.ErrNoRows, "Посылки с номером больше нет в БД")
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open(DbDriver, dataSourceName())
	require.NoError(t, err, "Нет ошибок при открытии БД")

	defer func() {
		if err = db.Close(); err != nil {
			log.Println(err)
		}
	}()
	store := NewParcelStore(db) // Структура - хранилище посылок

	// add
	testParcel := getTestParcel()
	number := add(t, store, testParcel)
	testParcel.Number = number

	// get
	p, err := store.Get(number)
	assert.NoError(t, err, "Нет ошибок при получении информации о посылке")
	assert.EqualValues(t, testParcel, p, "Значения полей в посылке совпадают с ожидаемыми")

	// delete
	delete(t, store, number)

}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open(DbDriver, dataSourceName())
	require.NoError(t, err, "Нет ошибок при открытии БД")

	defer func() {
		if err = db.Close(); err != nil {
			log.Println(err)
		}
	}()

	store := NewParcelStore(db) // Структура - хранилище посылок

	// add
	testParcel := getTestParcel()
	number := add(t, store, testParcel)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err, "Нет ошибок при изменении адреса")

	// check
	p, err := store.Get(number)
	assert.NoError(t, err, "Нет ошибок при получении информации о посылке")
	assert.Equal(t, newAddress, p.Address, "Значение адреса совпадает")

	// delete
	delete(t, store, number)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open(DbDriver, dataSourceName())
	require.NoError(t, err, "Нет ошибок при открытии БД")

	defer func() {
		if err = db.Close(); err != nil {
			log.Println(err)
		}
	}()

	store := NewParcelStore(db) // Структура - хранилище посылок

	// add
	testParcel := getTestParcel()
	number := add(t, store, testParcel)

	// set status
	newStatus := ParcelStatusDelivered
	err = store.SetStatus(number, newStatus)
	require.NoError(t, err, "Нет ошибок при изменении статуса")

	// check
	p, err := store.Get(number)
	assert.NoError(t, err, "Нет ошибок при получении информации о посылке")
	assert.Equal(t, newStatus, p.Status, "Значение адреса совпадает")

	// попытка удалить посылку со статусом, отличном от ParcelStatusRegistered
	err = store.Delete(number)
	require.NoError(t, err, "Нет ошибок при изменении статуса")
	_, err = store.Get(number)
	assert.NotErrorIsf(t, err, sql.ErrNoRows, "Нет ошибки о том, что посылки не существует")

	// Чистим за собой:
	// set status
	newStatus = ParcelStatusRegistered
	err = store.SetStatus(number, newStatus)
	require.NoError(t, err, "Нет ошибок при изменении статуса")

	// delete
	delete(t, store, number)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare

	db, err := sql.Open(DbDriver, dataSourceName())
	require.NoError(t, err, "Нет ошибок при открытии БД")

	defer func() {
		if err = db.Close(); err != nil {
			log.Println(err)
		}
	}()

	store := NewParcelStore(db) // Структура - хранилище посылок

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id := add(t, store, parcels[i])

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err, "Нет ошибок при получении посылок по клиенту")
	assert.Equal(t, len(parcels), len(storedParcels), "Найдено ожидаемое количество посылок (%d) по клиенту", len(parcels))

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		require.Equal(t, parcel.Number, parcelMap[parcel.Number].Number, "Найдена ожидаемая посылка")
		expect := parcelMap[parcel.Number]
		assert.EqualValues(t, expect, parcel, "Значения полей в посылке совпадают с ожидаемыми")

		// Чистим за собой
		delete(t, store, parcel.Number)
	}
}
