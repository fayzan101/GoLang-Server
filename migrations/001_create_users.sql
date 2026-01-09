-- Example migration: create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100),
    dob VARCHAR(20),
    university VARCHAR(100),
    semester VARCHAR(20),
    program VARCHAR(100),
    roll_no VARCHAR(20),
    email VARCHAR(100) UNIQUE,
    password VARCHAR(255),
    type VARCHAR(20)
);
