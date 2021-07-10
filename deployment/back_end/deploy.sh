# 构建最小镜像,需要在项目根目录下执行
docker build -t magicpowerworld/paotui_back_end:20210710 -f .\DockerFile.mini .

# 制作完镜像之后推送
docker push magicpowerworld/paotui_back_end:20210710

# 后端运行容器
docker run --name paotui_back_end --net=host -d magicpowerworld/paotui_back_end:20210710

# 本地测试
docker run --name paotui_back_end --net=host magicpowerworld/paotui_back_end:20210710