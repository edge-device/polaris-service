-- load using command:$ mysql -u root -p --verbose < create_tables.sql

CREATE DATABASE IF NOT EXISTS polaris;

USE polaris;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    user_id VARCHAR(255) NOT NULL PRIMARY KEY,
    created_at BIGINT NOT NULL,
    last_login BIGINT DEFAULT 0,
    firstname VARCHAR(255)  DEFAULT '',
    lastname VARCHAR(255) DEFAULT ''
);

-- devices table will be used for device auth and waiting room
DROP TABLE IF EXISTS devices;
CREATE TABLE devices (
    device_id VARCHAR(255) NOT NULL PRIMARY KEY,
    org_id VARCHAR(255) DEFAULT '',
    created_at BIGINT NOT NULL,
    last_seen BIGINT DEFAULT 0,
    joined_at BIGINT DEFAULT 0,
    device_key VARCHAR(255) NOT NULL, -- Make sure in base64
    profile_url VARCHAR(255) DEFAULT '',
    mac VARCHAR(255)  DEFAULT '',
    ip VARCHAR(255) DEFAULT '',
    properties MEDIUMTEXT
);
INSERT INTO devices (device_id, org_id, created_at, device_key, profile_url) 
VALUES
    ('device1', 'mcdonalds1002', 1677886013, 'VGhpcyBpcyBteSBzZWNyZXQga2V5', 'https://github.com'),
    ('device2', 'mcdonalds1002', 1677886013, 'my key', 'https://github.com'),
    ('device3', 'mcdonalds1001', 1665703610, 'my key', 'https://github.com'),
    ('device4', 'mcdonalds1001', 1677886013, 'my key', 'https://github.com'),
    ('device5', 'mcdonalds1001', 1665703610, 'my key', 'https://github.com');

-- refresh table
DROP TABLE IF EXISTS access_token;
CREATE TABLE access_token (
    user_id VARCHAR(255) NOT NULL,
    token_id INT AUTO_INCREMENT PRIMARY KEY,
    created_at BIGINT NOT NULL,
    expires_at BIGINT NOT NULL
);

-- Add 'api' user and privileges
CREATE USER IF NOT EXISTS 'api'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON polaris . * TO 'api'@'localhost';