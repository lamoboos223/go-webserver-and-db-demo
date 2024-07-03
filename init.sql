-- init.sql

-- Create table
CREATE TABLE IF NOT EXISTS authentication_keys (
    id SERIAL PRIMARY KEY,
    apikey VARCHAR(255) NOT NULL
);

-- Insert initial data
INSERT INTO authentication_keys (apikey) VALUES
('apikey1'),
('apikey2'),
('apikey3');
