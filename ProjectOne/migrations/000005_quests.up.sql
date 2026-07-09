CREATE TABLE
    public."quests" (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        type VARCHAR(50) NOT NULL CHECK (
            type IN ('watch_movie', 'write_review', 'invite_friend')
        ),
        required_count INT NOT NULL,
        reward_points INT NOT NULL
    );

CREATE TABLE
    public."user_quests" (
        user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        quest_id INT NOT NULL REFERENCES quests (id) ON DELETE CASCADE,
        progress INT NOT NULL DEFAULT 0,
        complete_at TIMESTAMP NULL DEFAULT NULL,
        PRIMARY KEY (user_id, quest_id)
    );