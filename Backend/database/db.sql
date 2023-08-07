CREATE DATABASE `comicbooks`
DEFAULT CHARACTER SET utf8
COLLATE utf8_hungarian_ci;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id` INT AUTO_INCREMENT,
    `username` VARCHAR(51) NOT NULL DEFAULT "", 
    `email` VARCHAR(255) NOT NULL DEFAULT "",
    `pswHash` VARCHAR(255) NOT NULL DEFAULT "",
    `createdAt` VARCHAR(255) NOT NULL DEFAULT "",
    `active` BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY(`id`)
);

DROP TABLE IF EXISTS `user_email_ver_hash`;
CREATE TABLE `user_email_ver_hash` (
    `username` VARCHAR(255) NOT NULL DEFAULT "",
    `verHash` VARCHAR(255) NOT NULL DEFAULT "",
    `timeout` VARCHAR(255) NOT NULL DEFAULT ""
);

DROP TABLE IF EXISTS `user_tokens`;
CREATE TABLE user_tokens (
    `userId` INT NOT NULL DEFAULT 0,
    `token` VARCHAR(255) NOT NULL DEFAULT ""
);

DROP TABLE IF EXISTS `volume_issues`;
DROP TABLE IF EXISTS `issue_characters`;
DROP TABLE IF EXISTS `volume_characters`;
DROP TABLE IF EXISTS `issues`;
DROP TABLE IF EXISTS `volumes`;
DROP TABLE IF EXISTS `characters`;

CREATE TABLE `volumes`(
    `id` INT NOT NULL DEFAULT 0,
    `name` VARCHAR(255) NOT NULL DEFAULT "",
    `img` VARCHAR(255) NOT NULL DEFAULT "",
    `desc` TEXT,
    `publisher` VARCHAR(255) NOT NULL DEFAULT "",
    PRIMARY KEY(`id`)
);

CREATE TABLE `issues` (
    `id` INT NOT NULL DEFAULT 0,
    `name` VARCHAR(255) NOT NULL DEFAULT "",
    `issue_number` INT NOT NULL DEFAULT 0,
    `img` VARCHAR(255) NOT NULL DEFAULT "",
    `cover_date` VARCHAR(255) NOT NULL DEFAULT "",
    `date_added` VARCHAR(255) NOT NULL DEFAULT "",
    PRIMARY KEY(`id`)
);

CREATE TABLE `characters` (
    `id` INT NOT NULL DEFAULT 0,
    `name` VARCHAR(255) NOT NULL DEFAULT "",
    `img` VARCHAR(255) NOT NULL DEFAULT "",
    PRIMARY KEY(`id`)
);

CREATE TABLE `volume_issues` (
    `volume_id` INT,
    `issue_id` INT
);

CREATE TABLE `volume_characters` (
    `volume_id` INT,
    `character_id` INT
);

CREATE TABLE `issue_characters` (
    `issue_id` INT,
    `character_id` INT
);