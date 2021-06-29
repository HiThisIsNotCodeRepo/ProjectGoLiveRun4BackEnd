# Installation Guide
## Using Source Code
I give out source code running manual first, I have tested it on my Mac and windows and working properly.
## Prerequisite
## Docker and MySQL
1. Download the latest Docker from [here](https://www.docker.com/)
2. Prepare MySQL container with following script.
```shell
docker run --name my-sql -dp 3306:3306 -e MYSQL_ROOT_PASSWORD=password  mysql:latest
```
3. Reset MySQL container localtime and restart
```shell
docker exec -it my-sql bash
ln -sf /usr/share/zoneinfo/Asia/Singapore /etc/localtime
exit
docker restart my-sql
```
4. Using below script to prepare MySQL database
```sql
drop database if exists paotui;
CREATE DATABASE paotui;
use paotui;
drop table if exists user;
CREATE TABLE user
(
    uid VARCHAR(36) COMMENT 'user id UUID',
    name VARCHAR(20) COMMENT 'user login name',
    password  VARCHAR(60) COMMENT 'user login password',
    email VARCHAR(40) COMMENT 'user email address',
    mobile_number INT UNSIGNED COMMENT 'user mobile number',
    last_login DATETIME COMMENT 'user last login datetime',
    PRIMARY KEY (uid)
);
drop table if exists task;
CREATE TABLE task
(
    task_id VARCHAR(36) COMMENT 'task id UUID',
    task_title VARCHAR(50) COMMENT 'task title required',
    task_description VARCHAR(100) COMMENT 'task description required',
    task_category_id TINYINT UNSIGNED COMMENT 'task category id 0~255 required',
    task_from VARCHAR(100) COMMENT 'task from address required',
    task_to VARCHAR(100) COMMENT 'task to address required',
    task_create DATETIME COMMENT 'task create datetime',
    task_start DATETIME COMMENT 'task start datetime required',
    task_complete DATETIME COMMENT 'task complete datetime',
    task_duration SMALLINT UNSIGNED COMMENT 'task duration in mins 0~65535 required',
    task_step TINYINT UNSIGNED COMMENT 'task current step 0~255',
    task_owner_id  VARCHAR(36) COMMENT 'task owner user id UUID required',
    task_owner_rate SMALLINT UNSIGNED COMMENT 'task owner expected rate 0~65535 required',
    task_deliver_id VARCHAR(36) COMMENT 'task deliver user id UUID',
    task_deliver_rate SMALLINT UNSIGNED COMMENT 'task deliver rate, final rate 0~65535',
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


```
## Clone repository
1. Ensure the MySQL container is on and `git clone https://github.com/qinchenfeng/ProjectGoLiveRun4BackEnd.git` to local `GOPATH`
2. `go mod tidy`
3. Populate MySQL
```sql
cd installation/prepare_database
go run .
```
4. Run the main app
```sql
cd app/
go run .
```