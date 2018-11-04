# rcmd
linux下通过ssh远程执行命令的工具
 命令用法
rcmd -h ${host} -u ${user} -p ${passwd} -c "cd ${path},./stop.sh,tar -zxvf wowo-bweb-0.0.1-SNAPSHOT.tar.gz, ./startup.sh"