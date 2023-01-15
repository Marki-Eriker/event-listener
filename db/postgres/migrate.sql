CREATE TABLE IF NOT EXISTS events
(
    id UUID NOT NULL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    created TIMESTAMP NOT NULL,
    system_name VARCHAR NOT NULL,
    message VARCHAR NOT NULL,
    incident BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    login VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    role INTEGER NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO users (login, password, role, verified)
    VALUES
        ('admin', '$2a$10$IciO74Nwo1P5xsxTaCE9gOdISPu/6gOtAP.dlG19DaY1UOjCjbxzG', 100, true),
        ('analyst', '$2a$10$IciO74Nwo1P5xsxTaCE9gOdISPu/6gOtAP.dlG19DaY1UOjCjbxzG', 200, true),
        ('analyst2', '$2a$10$IciO74Nwo1P5xsxTaCE9gOdISPu/6gOtAP.dlG19DaY1UOjCjbxzG', 200, false)
    ON CONFLICT DO NOTHING;

INSERT INTO events (id, event_id, created, system_name, message, incident)
    VALUES
        ('aa16d9e0-4cd3-4087-b701-a21e89954ffe', '5', '2023-01-15 16:53:10.283928', 'system_1', 'mock encoded data', false),
        ('ada14b45-db09-40ec-a595-8728ab80e0dd', '5', '2023-01-15 16:53:10.283928', 'system_2', 'mock encoded data', false),
        ('c195cd96-3fcd-4b35-830a-6a527f2cc8bf', '6', '2023-01-15 16:53:10.283928', 'system_3', 'mock encoded data', false),
        ('fc0075d0-d506-4ef7-b96c-b859bac6eca3', '6', '2023-01-15 16:53:10.283928', 'system_2', 'mock encoded data', false),
        ('86f147b4-82c7-4f13-88c3-1fb5269609cd', '7', '2023-01-15 16:53:10.283928', 'system_1', 'mock encoded data', false),
        ('6d5bbf90-bd54-46c3-a8df-39d3bca985a9', '7', '2023-01-15 16:53:10.283928', 'system_1', 'mock encoded data', false),
        ('28311c84-2b89-421d-849c-81eb62ad5a2e', '80', '2023-01-15 16:53:10.283928', 'system_2', 'mock encoded data', false),
        ('5de75ce6-0537-4a43-9d8c-b7eca99fbcd5', '80', '2023-01-15 16:53:10.283928', 'system_1', 'mock encoded data', false),
        ('3496e22b-9fe1-424b-b308-e587ad873eaf', '80', '2023-01-15 16:53:10.283928', 'system_1', 'mock encoded data', false),
        ('9b95be1e-e84e-4561-9768-865cfaa8381a', '80', '2023-01-15 16:53:10.283928', 'system_3', 'mock encoded data', false)
    ON CONFLICT DO NOTHING;
