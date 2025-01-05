CREATE TABLE IF NOT EXISTS "user" (
                                      id SERIAL PRIMARY KEY,
                                      firstname VARCHAR(50) NOT NULL,
    lastname VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    age INTEGER CHECK (age >= 0),  -- Ограничение на возраст (неотрицательное значение)
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
