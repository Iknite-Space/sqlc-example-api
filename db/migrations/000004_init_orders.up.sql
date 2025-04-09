CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    customer_name TEXT NOT NULL,
    amount INT NOT NULL,
    phone_number TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT now()
);