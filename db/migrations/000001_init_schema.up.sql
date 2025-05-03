CREATE TABLE "customer" (
"id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
"customer_name" VARCHAR(100) NOT NULL,
"phone_number" VARCHAR(20) NOT NULL UNIQUE,
"email" VARCHAR(100) NOT NULL UNIQUE,
"created_at" TIMESTAMP DEFAULT now()
);
