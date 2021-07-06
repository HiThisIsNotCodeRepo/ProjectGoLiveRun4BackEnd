# restart image download
docker stop my-sql
docker system prune -a --volumes # 2nd time onward->docker rm my-sql
docker run --name my-sql -dp 3306:3306 -e MYSQL_ROOT_PASSWORD=password  mysql:latest
docker cp mysql.sql my-sql:/tmp/
docker exec -it my-sql bash
ln -sf /usr/share/zoneinfo/Asia/Singapore /etc/localtime
exit
docker restart my-sql
docker exec -it my-sql bash
mysql -uroot -ppassword
source /tmp/mysql.sql


# 将含表的数据库重新制作为镜像
# 2d2d40728497为容器id
docker commit 2d2d40728497 magicpowerworld/paotui_mysql:20210706
# 制作完镜像之后推送
docker push magicpowerworld/paotui_mysql:20210706