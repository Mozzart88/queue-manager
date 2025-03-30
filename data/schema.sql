BEGIN TRANSACTION;

DROP TABLE IF EXISTS message_status;
CREATE TABLE message_status (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

INSERT INTO message_status
VALUES
(0,'done'),
(1,'new'),
(2,'active');

DROP TABLE IF EXISTS publisher;
CREATE TABLE publisher (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

DROP TABLE IF EXISTS message;
CREATE TABLE message (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    publisher_id INTEGER NOT NULL,
    status_id INTEGER DEFAULT 1,
    FOREIGN KEY (publisher_id) REFERENCES publisher (id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES message_status (id) ON DELETE CASCADE
);

DROP VIEW IF EXISTS queue;
CREATE VIEW queue AS
SELECT
    message.id as id,
    message.status_id as status_id,
    message_status.name as status,
    message.publisher_id as publisher_id,
    publisher.name as publisher,
    message.content as msg
FROM message
JOIN message_status ON message.status_id = message_status.id
JOIN publisher ON message.publisher_id = publisher.id;

COMMIT;
