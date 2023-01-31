package db

import "github.com/labstack/gommon/log"

func GetCarsIDsByID(userID int32) ([]int32, error) {
	rows, err := connection.db.Query(` -- TODO: почему у тебя все запросы не исходят от той таблицы по которой идет условие? - перепиши запрос
	SELECT 
		uc.car_id
	FROM users_cars uc
		INNER JOIN users u ON u.id=uc.user_id 
	WHERE uc.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		carIDs []int32
		carID  int32 // TODO: везде определяй переменную для Scan() внутри Next чтобы была на виду
	)

	for rows.Next() {
		if err := rows.Scan(&carID); err != nil {
			log.Printf("error trying scan car IDs: %v", err)

			continue // TODO: почему ты везде скипаешь ошибки, в чем идея?
		}

		carIDs = append(carIDs, carID)
	}

	return carIDs, nil
}
