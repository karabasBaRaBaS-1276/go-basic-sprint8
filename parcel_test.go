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

	store := NewParcelStore(db)        // Структура - хранилище посылок
	service := NewParcelService(store) // Структура - сервис по работе с посылками

	// add
	testParcel := getTestParcel()
	number, err := service.store.Add(testParcel)
	require.NoError(t, err, "Нет ошибок при регистрации посылки")
	require.NotEmpty(t, number, "Новой посылке присвоен номер")

	// get
	p, err := service.store.Get(number)
	assert.NoError(t, err, "Нет ошибок при получении информации о посылке")
	assert.Equal(t, testParcel.Client, p.Client, "Значение клиента совпадает")
	assert.Equal(t, testParcel.Address, p.Address, "Значение адреса совпадает")
	assert.Equal(t, testParcel.Status, p.Status, "Значение статуса совпадает")
	assert.Equal(t, testParcel.CreatedAt, p.CreatedAt, "Значение даты создания совпадает")

	// delete
	err = service.store.Delete(number)
	assert.NoError(t, err, "Нет ошибок при удалении информаци о посылке")

	// проверьте, что посылку больше нельзя получить из БД
	_, err = service.store.Get(number)
	assert.ErrorIsf(t, err, sql.ErrNoRows, "Посылки с номером больше нет в БД")

}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	//db, err := // настройте подключение к БД

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	//newAddress := "new test address"

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	//db, err := // настройте подключение к БД

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	// set status
	// обновите статус, убедитесь в отсутствии ошибки

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	/*
		db, err := // настройте подключение к БД

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
			id, err := // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

			// обновляем идентификатор добавленной у посылки
			parcels[i].Number = id

			// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
			parcelMap[id] = parcels[i]
		}

		// get by client
		storedParcels, err := // получите список посылок по идентификатору клиента, сохранённого в переменной client
		// убедитесь в отсутствии ошибки
		// убедитесь, что количество полученных посылок совпадает с количеством добавленных

		// check
		for _, parcel := range storedParcels {
			// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
			// убедитесь, что все посылки из storedParcels есть в parcelMap
			// убедитесь, что значения полей полученных посылок заполнены верно
		}
	*/
}
