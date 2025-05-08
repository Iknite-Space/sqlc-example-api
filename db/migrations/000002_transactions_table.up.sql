
CREATE TABLE "transactions" (
  "reference" VARCHAR(36) PRIMARY KEY,
  "external_reference" VARCHAR(36) NOT NULL,
  "status_id" VARCHAR(36) NOT NULL,
  "amount" FLOAT8 NOT NULL,
  "currency" VARCHAR(36) NOT NULL,
  "operator" VARCHAR(36) NOT NULL,
  "code" VARCHAR(36) NOT NULL,
  "operator_reference" VARCHAR(36) NOT NULL,
  "description_id" VARCHAR(36) NOT NULL,
  "external_user" VARCHAR(36) NOT NULL,
  "reason" VARCHAR(36) NOT NULL,
  "phone_number" VARCHAR(36) NOT NULL,
  "create_at" TIMESTAMP DEFAULT now()
);