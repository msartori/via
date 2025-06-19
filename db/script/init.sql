CREATE TABLE IF NOT EXISTS operators (
    id SERIAL PRIMARY KEY,
    account VARCHAR(200) UNIQUE NOT NULL,
    enabled boolean NOT NULL DEFAULT false, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS guides (
    id SERIAL PRIMARY KEY,
    via_guide_id CHAR(12) NOT NULL,
    recipient VARCHAR(100) NOT NULL,
    status VARCHAR(30) NOT NULL,
    operator_id INTEGER NOT NULL REFERENCES operators(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS guide_histories (
    id SERIAL PRIMARY KEY,
    guide_id INTEGER NOT NULL REFERENCES guides(id),
    status VARCHAR(30) NOT NULL, 
    operator_id INTEGER NOT NULL REFERENCES operators(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION insert_guide_histories()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO guide_histories (guide_id, status, operator_id)
  VALUES (NEW.id, NEW.status, NEW.operator_id);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_insert_guide_histories
AFTER INSERT OR UPDATE ON guides
FOR EACH ROW
EXECUTE FUNCTION insert_guide_histories();

-- Insert SYSTEM operator only if not exists
INSERT INTO operators (id, account, enabled)
SELECT 1, 'SYSTEM', true
WHERE NOT EXISTS (SELECT 1 FROM operators WHERE id = 1);

-- Advance the sequence safely
SELECT setval(pg_get_serial_sequence('operators', 'id'), GREATEST(MAX(id), 1)) FROM operators;

INSERT INTO operators (account, enabled) values ('miguel.sartori@gmail.com', 'true');
