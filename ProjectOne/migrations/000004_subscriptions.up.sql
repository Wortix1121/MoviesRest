CREATE TABLE
    public."subscription_plans" (
        id SERIAL PRIMARY KEY,
        name VARCHAR(250) NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        features JSON
    );

CREATE TABLE
    public."subscriptions" (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        plan_id INT NOT NULL REFERENCES subscription_plans (id) ON DELETE CASCADE,
        start_date DATE NOT NULL,
        end_date DATE NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT true,
        CHECK (end_date >= start_date)
    );

CREATE UNIQUE INDEX idx_active_subscriptions_per_user ON subscriptions (user_id)
WHERE
    is_active = true;