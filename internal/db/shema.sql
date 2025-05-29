CREATE TABLE users (
    id int  PRIMARY KEY auto increment, 
    public_id VARCHAR(36) NOT NULL UNIQUE,--public facing id (UUID)
    username varchar(20) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    password varchar(255) NOT NULL
);

CREATE TABLE posts (
    id int PRIMARY KEY auto increment,
    title varchar(50) NOT NULL,
    text text ,
    poster int NOT NULL,
    posted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (poster) REFERENCES users(public_id) ON DELETE CASCADE
);