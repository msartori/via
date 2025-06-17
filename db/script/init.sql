CREATE TABLE IF NOT EXISTS operator (
    id SERIAL PRIMARY KEY,
    account UNIQUE NOT NULL VARCHAR(200),
    enabled boolean, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS guide (
    id SERIAL PRIMARY KEY,
    via_guide_id CHAR(12) UNIQUE NOT NULL,
    recipient VARCHAR(100),
    status VARCHAR(30),
    operator_id INTEGER REFERENCES operator(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS guide_history (
    id SERIAL PRIMARY KEY,
    guide_id INTEGER NOT NULL REFERENCES guide(id),
    status VARCHAR(30), 
    operator_id INTEGER REFERENCES operator(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO operator (account, enabled) values ('miguel.sartori@gmail.com', 'true')
