CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(20),
    surname VARCHAR(20),
    age INTEGER
);

INSERT INTO users(name, surname, age) 
VALUES('Ivan', 'Sergeev', 20), 
('Fedor','Zaytsev', 35), 
('Vladimir','Ivanov', 45), 
('Petr','Ivanov', 24), 
('John','Smith', 60);

CREATE TABLE IF NOT EXISTS users_cars(
    user_id INTEGER,
    car_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO users_cars(user_id, car_id) 
VALUES(1, 1), (1, 2), (2, 2), (2, 3), (3, 3), (3, 4), (4, 4), (4, 5), (5, 5), (5, 1);

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS users_cars;

CREATE TABLE IF NOT EXISTS car_configurations(
    id SERIAL PRIMARY KEY,
    concern VARCHAR(20),
    model VARCHAR(20),
    year INTEGER,
    engine_id INTEGER
);

INSERT INTO car_configurations(concern, model, year, engine_id) 
VALUES('BMW', 'M5', 1999, 2),
('BMW', 'M3', 2002, 1),
('Mersedes', 'C-class', 2000, 3),
('Mersedes', 'E-class', 2022, 6),
('Toyota', 'Camry', 2022, 4),
('Opel', 'Astra', 2010, 5),
('Kia', 'Optima', 2022, 1);

CREATE TABLE IF NOT EXISTS cars(
    id SERIAL PRIMARY KEY,
    configuration_id integer,
    used BOOLEAN,
    price integer,
    FOREIGN KEY (configuration_id) REFERENCES car_configurations(id)
);

INSERT INTO cars(configuration_id, used, price)
VALUES(1, true, 2000000),
(2, false, 2400000),
(3, false, 1800000),
(4, true, 8000000),
(5, false, 4000000),
(6, true, 800000),
(7, false, 220000);

DROP TABLE IF EXISTS cars;

DROP TABLE IF EXISTS car_configurations;

CREATE TABLE IF NOT EXISTS engines(
    id SERIAL PRIMARY KEY,
    volume DECIMAL(2, 1)
);

INSERT INTO engines(volume) 
VALUES(1.6), (2.5), (4.0), (3.5), (6.0), (3.0);

DROP TABLE IF EXISTS engines;
