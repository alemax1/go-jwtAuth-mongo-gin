package db

import (
	"amaksimov/engineService/models"
	"log"

	"github.com/lib/pq"
)

// TODO: пиши короткие sql запросы в одной строке
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
			log.Printf("error trying scan engine: %v", err)

			continue
		}

		engines = append(engines, e)
	}

	return engines, nil
}

// TODO: пиши короткие sql запросы в одной строке "SELECT DISTINCT id, volume FROM engines WHERE id = ANY($1)"
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
		log.Printf("error trying scan engine: %v", err)
		return nil, err
	}

	return &e, nil
}

// TODO: пиши короткие sql запросы в одной строке "SELECT DISTINCT id, volume FROM engines WHERE id = ANY($1)"
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
			log.Printf("error trying scan engine: %v", err) // TODO: все логи что ты написал в репозитории - удалить, тут их не должно быть

			continue
		}

		engines = append(engines, e)
	}

	return engines, nil
}

// TODO: запрос не длинный, можно записать его в одну строчку сразу как аргумент запроса
// TODO: зачем тут подготавливать запрос, для чего ты это написал?
func CreateEngine(en models.Engine) (*models.Engine, error) {
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

	var e models.Engine

	if err = stmt.QueryRow(en.Volume).Scan(&e.ID, &e.Volume); err != nil {
		log.Printf("error trying scan engine: %v", err)
		return nil, err
	}

	return &e, nil
}

// TODO: запрос не длинный, можно записать его в одну строчку сразу как аргумент запроса
// TODO: зачем тут подготавливать запрос, для чего ты это написал?
func UpdateEngine(en models.Engine) (*models.Engine, error) {
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

	var e models.Engine

	if err = stmt.QueryRow(en.ID, en.Volume).Scan(&e.ID, &e.Volume); err != nil {
		log.Printf("error trying scan engine: %v", err)
		return nil, err
	}

	return &e, nil
}

// TODO: запрос не длинный, можно записать его в одну строчку сразу как аргумент запроса
func DeleteEngine(engineID int32) error {
	if _, err := connection.db.Exec(`
	DELETE FROM engines
	where id = $1`, engineID); err != nil {
		return err
	}

	return nil
}
