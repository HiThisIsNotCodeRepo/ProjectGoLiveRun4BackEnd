# Go School Project Go Live Run 4
## Project Title: Pao Tui(跑腿)
**Back End API**
>Author: Qin Chenfeng
>
> Email:freddy.qin@gmail.com

## Project Description
Due to Covid 19, an ad hoc job posting platform has emerged to provide more job opportunities, anyone can post their urgent task with personalised requirements ,for instance deliver food, buy necessity, send documents etc, it also needs to include max acceptable rate and expected delivery time. Anyone who is interested in making pocket money can bid for the job with their minimum acceptable rate and finish time. Among those service providers, job posters should pick one. Once a job is assigned, the service provider should do their best to complete the task before the deadline to avoid any penalty.

## Database Design Consideration
1. The task and user table is based on 1-to-1 relationship. It means one task one record and one user one record in respective table, it has primary key for `task_id` and `uid`. I prefer to use UUID in these primary key as it's random and makes it safer.
2. The task_bid table is differnt it's n-to-n relationship, because one task can have multiple task bidder and one task bidder can bid mutliple task. So there's no primary key set in this table.
3. To make searching speed quicker in task_bid table, `task_id` and `task_bidder_rate` has been set to `Index`.
4. To minimise the space occupation, some field has been set with unsigned number, for instance the mobile number I set it to `INT UNSIGNED` it can provide 0~4_294_967_295 10 digits while it's 4 bytes.

The sql file can be found [here](https://github.com/qinchenfeng/ProjectGoLiveRun4BackEnd/blob/master/doc/sql/mysql.sql).


## Core feature
1. User registeration
2. User login
3. Front end and back end communication
4. Database CRUD
5. SSL certificate
6. Container orchestration(k8s)

