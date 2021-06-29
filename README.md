# Go School Project Go Live Run 4

## Project Title: Pao Tui(跑腿)

**Back End API**
> Author: Qin Chenfeng
>
> Email:freddy.qin@gmail.com

## Project Description

Due to Covid 19, an ad hoc job posting platform has emerged to provide more job opportunities, anyone can post their
urgent task with personalised requirements ,for instance deliver food, buy necessity, send documents etc, it also needs
to include max acceptable rate and expected delivery time. Anyone who is interested in making pocket money can bid for
the job with their minimum acceptable rate and finish time. Among those service providers, job posters should pick one.
Once a job is assigned, the service provider should do their best to complete the task before the deadline to avoid any
penalty.

## Database Design Consideration

1. The task and user table is based on 1-to-1 relationship. It means one task one record and one user one record in
   respective table, it has primary key for `task_id` and `uid`. I prefer to use UUID in these primary key as it's
   random and makes it safer.
2. The task_bid table is differnt it's n-to-n relationship, because one task can have multiple task bidder and one task
   bidder can bid mutliple task. So there's no primary key set in this table.
3. To make searching speed quicker in task_bid table, `task_id` and `task_bidder_rate` has been set to `Index`.
4. To minimise the space occupation, some field has been set with unsigned number, for instance the mobile number I set
   it to `INT UNSIGNED` it can provide 0~4_294_967_295 10 digits while it's 4 bytes.

The sql file can be found [here](https://github.com/qinchenfeng/ProjectGoLiveRun4BackEnd/blob/master/doc/sql/mysql.sql).

## Installation Manual


## Authrorisation with token

The main strategy is to encrypt a piece of text hold user info and expire time. Everytime when front end access
restricted website it needs to send token for verification. Here I verified if username in url and token is the same and
token has expired or not.

## REST API

|No| Module  | Method|Request Data|API URL|Query Key|
|---|---|---|---|---|---|
|1|Auth|Post,Options|Yes|/api/v1/auth|option{login,token-verify}|
|2|Category|Get|No|/api/v1/categories|
|3|My Info|Get|No|/api/v1/spending/{userID}|chart-type{card,summary,datasource},date{yesterday,two-days-ago,three-days-ago},category{buy-necessity,food-delivery,send-document,other},date{last-week,this-week}|
|4|My Info|Get|No|/api/v1/earning/{userID}|chart-type{card,radar,datasource},date{last-week,this-week}|
|5|Task|Post,Options|Yes|/api/v1/tasks/task||
|6|Task|Get|No|/api/v1/tasks/{id}|identity{user,task},options{on-going},category{only-me,exclude-me}|
|7|Task|Put,Delete,Option|Yes|/api/v1/tasks/task/{taskID}|option{confirm-task-deliver,update-expected-rate,delete}|
|8|Task Bid|Post,Option|Yes|/api/v1/tasks/bid|

## mySQL connection time issue

When deal with MySQL with Docker if the container timezone is not set properly, the record datetime may not be correct.
By default the timezone in image
is [UTC/Universal Time Coordinated / Universal Coordinated Time](https://www.timeanddate.com/worldclock/timezone/utc).

In order to set it to local time we need to:

1. Get in to container bash and reset `/etc/localtime`

```shell
ln -sf /usr/share/zoneinfo/Asia/Singapore /etc/localtime
exit
```

after that reset container.

```shell
docker restart my-sql
```

2. In go put `charset=utf8&parseTime=True&loc=Local` after database connection url.

## Handle CORS in gorilla

Example in [github](https://github.com/gorilla/mux#handling-cors-requests)

## SQL Note

Get current week's Sunday, if today is Sunday it will fall on next Sunday

```sql
select subdate(curdate(), date_format(curdate(), '%w') - 7)
```

Convert datetime to y-m-d

```sql
select DATE_FORMAT(task_complete, '%Y-%m-%d')
from task
```

## Core feature

1. User registeration
2. User login
3. Front end and back end communication
4. Database CRUD
5. SSL certificate
6. Container orchestration(k8s)

