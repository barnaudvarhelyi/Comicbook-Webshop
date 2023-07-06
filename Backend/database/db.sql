CREATE DATABASE comicbooks
DEFAULT CHARACTER SET utf8
COLLATE utf8_hungarian_ci;

DROP TABLE IF EXISTS `users`;
CREATE TABLE users (
    `id` INT AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL DEFAULT "",
    `password` VARCHAR(255) NOT NULL DEFAULT "",
    `email` VARCHAR(255) NOT NULL DEFAULT "",
    PRIMARY KEY(`id`)
);

INSERT INTO `users` 
VALUES (1, 'Admin', 'Admin', 'admin@gmail.com');

DROP TABLE IF EXISTS `user_tokens`;
CREATE TABLE user_tokens (
    `userId` INT NOT NULL DEFAULT 0,
    `token` VARCHAR(255) NOT NULL DEFAULT ""
);

DROP TABLE IF EXISTS `comicbooks`;
DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
    `id` INT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL DEFAULT ""
);

CREATE TABLE `comicbooks` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL DEFAULT "",
    `price` DECIMAL(10, 2),
    `category_id` INT,
    FOREIGN KEY (category_id) REFERENCES category(id)
);