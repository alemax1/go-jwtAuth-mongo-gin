package db

import (
	"amaksimov/carService/models"
	"log"

	"github.com/lib/pq"
)

func GetCarsByIDs(IDs []int32) ([]models.Car, error) {
	rows, err := connection.db.Query(`
	SELECT 
		cc.id, cc.concern, cc.model, cc.year, cc.used, cc.engine_id, c.price
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
			&c.Used,
			&c.EngineID,
			&c.Price,
		); err != nil {
			log.Printf("error trying scan car IDs: %v", err)

			continue
		}

		cars = append(cars, c)
	}

	return cars, nil
}

func GetEngineID(carID int32) (int32, error) {
	row := connection.db.QueryRow(`
	SELECT
		cc.engine_id
	FROM car_configurations cc
		INNER JOIN cars c ON cc.id=c.configuration_id
	WHERE c.configuration_id = $1`, carID)

	if err := row.Err(); err != nil {
		return 0, err
	}

	var eID int32

	if err := row.Scan(&eID); err != nil {
		return 0, err
	}

	return eID, nil
}

func GetEnginesByIDs(carIDs []int32) ([]int32, error) {
	rows, err := connection.db.Query(`
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

			continue
		}

		enginesIDs = append(enginesIDs, eID)
	}

	return enginesIDs, nil
}
