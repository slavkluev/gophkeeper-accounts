CREATE TABLE IF NOT EXISTS accounts
(
    id       INTEGER PRIMARY KEY,
    login    TEXT    NOT NULL,
    pass     TEXT    NOT NULL,
    info     TEXT    NOT NULL DEFAULT '',
    user_uid INTEGER NOT NULL
);
