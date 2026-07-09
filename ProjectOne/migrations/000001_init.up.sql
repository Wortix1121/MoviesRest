CREATE TABLE
    public."users" (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        role VARCHAR(50) NOT NULL DEFAULT 'user',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    public."movies" (
        id SERIAL PRIMARY KEY,
        title VARCHAR(200) NOT NULL UNIQUE,
        description TEXT,
        duration_minutes INT NOT NULL,
        release_date DATE NOT NULL,
        poster_url VARCHAR(500) NOT NULL,
        age_rating VARCHAR(10) NOT NULL,
        type VARCHAR(20) NOT NULL CHECK (type IN ('movie', 'series')),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    public."genres" (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL UNIQUE
    );

CREATE TABLE
    public."movie_genres" (
        movie_id INT NOT NULL REFERENCES movies (id) ON DELETE CASCADE,
        genre_id INT NOT NULL REFERENCES genres (id) ON DELETE CASCADE,
        PRIMARY KEY (movie_id, genre_id)
    );

CREATE INDEX idx_movie_genres_movie_id ON movie_genres (movie_id);

CREATE INDEX idx_movie_genres_genre_id ON movie_genres (genre_id);