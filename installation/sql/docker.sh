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
