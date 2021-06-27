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
## Authrorisation with token
The main strategy is to encrypt a piece of text hold user info and expire time. Everytime when front end access restricted website it needs to send token for verification. Here I verified if username in url and token is the same and token has expired or not.

## Update on api

| No  | API URL| 
| -------- | -------- | 
|1| /spending/yesterday/buy-necessity/{userID}|
|2| /spending/yesterday/food-delivery/{userID}|
|3| /spending/yesterday/send-document/{userID}|
|4| /spending/yesterday/other/{userID}|
|5| /spending/two-days-ago/buy-necessity/{userID}|
|6| /spending/two-days-ago/food-delivery/{userID}|
|7| /spending/two-days-ago/send-document/{userID}|
|8| /spending/two-days-ago/other/{userID}|
|9| /spending/three-days-ago/buy-necessity/{userID}|
|10|/spending/three-days-ago/food-delivery/{userID}|
|11|/spending/three-days-ago/send-document/{userID}|
|12|/spending/three-days-ago/other/{userID}|
|13|/spending/this-week/summary/{userID}|
|14|/spending/last-week/summary/{userID}|
|15|/spending/tasks/{userID}|
|16|/earning/tasks/{userID}|
|17|/earning/past-days/{userID}|
|18|/earning/last-week/radar/{userID}|
|19|/earning/this-week/radar/{userID}|
|20|/tasks/task|
|21|/auth/login|
|22|/auth/token-verify/{userID}
|

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

