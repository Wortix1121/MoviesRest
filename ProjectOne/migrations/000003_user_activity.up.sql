CREATE TABLE
    public."user_watch_history" (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        content_type VARCHAR(20) NOT NULL CHECK (content_type IN ('movie', 'episode')),
        content_id INT NOT NULL,
        progress_seconds INT NOT NULL,
        watched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    public."user_reviews" (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        movie_id INT NOT NULL REFERENCES movies (id) ON DELETE CASCADE,
        rating INT NOT NULL CHECK (
            rating >= 1
            AND rating <= 5
        ),
        comment TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (user_id, movie_id)
    );

CREATE TABLE
    public."user_favorites" (
        user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        movie_id INT NOT NULL REFERENCES movies (id) ON DELETE CASCADE,
        PRIMARY KEY (user_id, movie_id)
    );