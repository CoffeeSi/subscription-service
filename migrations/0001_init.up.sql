CREATE TABLE IF NOT EXISTS subscriptions(
    id SERIAL NOT NULL,
    service_name varchar NOT NULL,
    price integer NOT NULL, 
    user_id uuid NOT NULL,
    start_date date NOT NULL,
    end_date date,
    PRIMARY KEY (id),
    CONSTRAINT subscriptions_price_check CHECK (price >= 0)
);