CREATE TABLE
    public."seasons" (
        id SERIAL PRIMARY KEY,
        movie_id INT NOT NULL REFERENCES movies (id) ON DELETE CASCADE,
        season_number INT NOT NULL,
        title VARCHAR(200) NOT NULL,
        UNIQUE (movie_id, season_number)
    );

CREATE TABLE
    public."episodes" (
        id SERIAL PRIMARY KEY,
        season_id INT NOT NULL REFERENCES seasons (id) ON DELETE CASCADE,
        episode_number INT NOT NULL,
        title VARCHAR(200) NOT NULL,
        duration_minutes INT NOT NULL,
        video_url VARCHAR(255) NOT NULL,
        UNIQUE (season_id, episode_number)
    );