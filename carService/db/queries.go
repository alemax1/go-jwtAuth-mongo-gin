package db

import (
	"amaksimov/carService/models"
	"log"

	"github.com/lib/pq"
)

func GetCarsByIDs(IDs []int32) ([]models.Car, error) {
	rows, err := connection.db.Query(` -- TODO: хочу видеть другой запрос получающий тот же результат
	SELECT 
		cc.id, cc.concern, cc.model, cc.year, cc.engine_id, c.price, c.used
	FROM car_configurations cc
		INNER JOIN cars c ON c.configuration_id=cc.id
	WHERE c.id = any($1)`, pq.Array(IDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		c    models.Car
		cars []models.Car
	)

	for rows.Next() {
		if err := rows.Scan(
			&c.ID,
			&c.Concern,
			&c.Model,
			&c.Year,
			&c.EngineID,
			&c.Price,
			&c.Used,
		); err != nil {
			log.Printf("error trying scan car IDs: %v", err)

			continue // TODO: почему просто скипаем при возникновении ошибки?
		}

		cars = append(cars, c)
	}

	return cars, nil
}

func GetEngineID(carID int32) (int32, error) {
	row := connection.db.QueryRow(` -- тоже самое, хочу видеть тот же результат с другим запросом, речь про то что ты передаешь аргумент carID от этого и отталкиваться
	SELECT
		cc.engine_id
	FROM car_configurations cc
		INNER JOIN cars c ON cc.id=c.configuration_id
	WHERE c.configuration_id = $1`, carID)

	if err := row.Err(); err != nil { // TODO: все вот тут что ты написал можно ограничиться в 3 раза меньше строк кода, при этому не понадобиться враппить ошибки
		return 0, err // TODO: где враппинг?
	}

	var eID int32

	if err := row.Scan(&eID); err != nil {
		return 0, err // TODO: где враппинг?
	}

	return eID, nil
}

func GetEnginesByIDs(carIDs []int32) ([]int32, error) {
	rows, err := connection.db.Query(` -- TODO: снова, жду другой запрос исходя переданного аргумента
	SELECT
		cc.engine_id
	FROM car_configurations cc
		INNER JOIN cars c ON cc.id=c.configuration_id
	WHERE c.id = any($1)`, pq.Array(carIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		eID        int32
		enginesIDs []int32
	)

	for rows.Next() {
		if err := rows.Scan(&eID); err != nil {
			log.Printf("error trying scan car IDs: %v", err)

			continue // TODO: тот же вопрос - почему скипаем при получении ошибки?
		}

		enginesIDs = append(enginesIDs, eID)
	}

	return enginesIDs, nil
}

// это что такое? где ты увидел такую работу с транзакцией? с чего  вдруг ты в defer вызываешь rollback?
func DeleteCarConfiguration(engineID int32) error {
	txn, err := connection.db.Begin()
	if err != nil {
		return err // TODO: где враппинг ошибки?
	}

	defer txn.Rollback() // тут должна быть обработка ошибки и исходя её уже выполнять либо rollback либо commit

	// TODO:
	// TODO: не хочу видеть 2 запроса, хочу видеть один
	// TODO:

	if _, err := txn.Exec(`
	DELETE FROM cars c
		USING car_configurations cc
	WHERE c.configuration_id=cc.id AND 
	cc.engine_id=$1
	`, engineID); err != nil {
		return err // TODO: где враппинг ошибки?
	}

	if _, err := txn.Exec(`
	DELETE 
	FROM car_configurations
	WHERE engine_id = $1`, engineID); err != nil {
		return err // TODO: где враппинг ошибки?
	}

	if err := txn.Commit(); err != nil { // я не вижу причин использовать тут коммит
		return err // TODO: где враппинг ошибки?
	}

	return nil
}
