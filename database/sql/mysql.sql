drop database if exists paotui;
CREATE DATABASE paotui DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
use paotui;
drop table if exists user;
CREATE TABLE user
(
    uid VARCHAR(36) COMMENT 'user id UUID',
    name VARCHAR(20) COMMENT 'user login name',
    password  VARCHAR(44) COMMENT 'user login password',
    email VARCHAR(40) COMMENT 'user email address',
    mobile_number INT UNSIGNED COMMENT 'user mobile number',
    last_login DATETIME COMMENT 'user last login datetime',
    user_delete TINYINT UNSIGNED COMMENT 'logical bit represent delete status',
    PRIMARY KEY (uid)
);
drop table if exists task;
CREATE TABLE task
(
    task_id VARCHAR(36) COMMENT 'task id UUID',
    task_title VARCHAR(50) COMMENT 'task title',
    task_description VARCHAR(100) COMMENT 'task description',
    task_category_id TINYINT UNSIGNED COMMENT 'task category id 0~255',
    task_from VARCHAR(100) COMMENT 'task from address',
    task_to VARCHAR(100) COMMENT 'task to address',
    task_create DATETIME COMMENT 'task create datetime',
    task_start DATETIME COMMENT 'task start datetime',
    task_complete DATETIME COMMENT 'task complete datetime',
    task_duration SMALLINT UNSIGNED COMMENT 'task duration in mins 0~65535',
    task_step TINYINT UNSIGNED COMMENT 'task current step 0~255',
    task_owner_id  VARCHAR(36) COMMENT 'task owner user id UUID',
    task_owner_rate SMALLINT UNSIGNED COMMENT 'task owner expected rate 0~65535',
    task_deliver_id VARCHAR(36) COMMENT 'task deliver user id UUID',
    task_deliver_rate SMALLINT UNSIGNED COMMENT 'task deliver rate, final rate 0~65535',
    task_delete TINYINT UNSIGNED COMMENT 'logical bit represent delete status',
    PRIMARY KEY (task_id)
);
CREATE TABLE task_bid
(
    task_id VARCHAR(36) COMMENT 'task id UUID',
    task_bidder_id VARCHAR(36) COMMENT 'task bidder id UUID',
    task_bidder_rate SMALLINT UNSIGNED COMMENT 'task bidder rate 0~65535',
    INDEX(task_id),
    Index(task_bidder_id)
);
CREATE TABLE category
(
    cid TINYINT UNSIGNED COMMENT 'category id',
    title VARCHAR(20) COMMENT 'category title',
    PRIMARY KEY(cid)
);
CREATE USER 'user'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON  paotui.* TO 'user'@'%';
FLUSH PRIVILEGES;

