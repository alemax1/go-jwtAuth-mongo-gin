package db

import (
	"amaksimov/engineService/models"
	"amaksimov/pkg/grpc/pb"
	"log"

	"github.com/lib/pq"
)

func GetAllEngines() ([]models.Engine, error) {
	rows, err := connection.db.Query(`
	SELECT
		id, volume
	FROM engines`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		e       models.Engine
		engines []models.Engine
	)

	for rows.Next() {
		if err := rows.Scan(&e.ID, &e.Volume); err != nil {
			log.Printf("error trying scan car IDs: %v", err)

			continue
		}

		engines = append(engines, e)
	}

	return engines, nil
}

func GetEngineByID(ID int32) (*models.Engine, error) {
	row := connection.db.QueryRow(
		`SELECT 
			id, volume
		FROM engines
		WHERE id = $1`, ID)

	if err := row.Err(); err != nil {
		return nil, err
	}

	var e models.Engine

	if err := row.Scan(&e.ID, &e.Volume); err != nil {
		return nil, err
	}

	return &e, nil
}

func GetEnginesByIDs(IDs []int32) ([]models.Engine, error) {
	rows, err := connection.db.Query(`
	SELECT DISTINCT
		id, volume
	FROM engines
	WHERE id = any($1)`, pq.Array(IDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		e       models.Engine
		engines []models.Engine
	)

	for rows.Next() {
		if err := rows.Scan(&e.ID, &e.Volume); err != nil {
			log.Printf("error trying scan car IDs: %v", err)

			continue
		}

		engines = append(engines, e)
	}

	return engines, nil
}

func CreateEngine(e *pb.Engine) (*models.Engine, error) {
	query := `
	INSERT INTO 
		ENGINES(volume)
	VALUES($1)
	RETURNING id, volume`

	stmt, err := connection.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var (
		ID     int32
		volume float32
	)

	if err = stmt.QueryRow(e.Volume).Scan(&ID, &volume); err != nil {
		return nil, err
	}

	return &models.Engine{ID: ID, Volume: volume}, nil
}

func UpdateEngine(engineID int32, v float32) (*models.Engine, error) {
	query := `
	UPDATE engines
	SET volume = $2
	WHERE id = $1
	RETURNING id, volume`

	stmt, err := connection.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var (
		ID     int32
		volume float32
	)

	if err = stmt.QueryRow(engineID, volume).Scan(&ID, &volume); err != nil {
		return nil, err
	}

	return &models.Engine{ID: ID, Volume: volume}, nil
}

func DeleteEngine(engineID int32) error {
	if _, err := connection.db.Exec(`
	DELETE FROM engines
	where id = $1`, engineID); err != nil {
		return err
	}

	return nil
}
