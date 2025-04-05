-- create thread table
CREATE TABLE IF NOT EXISTS thread (
    id SERIAL PRIMARY KEY,
    title VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- modify table message
ALTER TABLE message DROP COLUMN thread;

ALTER TABLE message ADD COLUMN thread_id INT NOT NULL;

-- add foreign key contrain
ALTER TABLE message ADD CONSTRAINT fk_thread_id
FOREIGN KEY (thread_id) REFERENCES thread(id) ON DELETE CASCADE;