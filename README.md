# rcmd
linux下通过ssh远程执行命令的工具
 命令用法
rcmd -h ${host} -u ${user} -p ${passwd} -c "cd ${path},./stop.sh,tar -zxvf wowo-bweb-0.0.1-SNAPSHOT.tar.gz, ./startup.sh"

# al
linux macos保存服务器信息的ssh登录工具
特别说明：本工具大部分代码复制[https://github.com/islenbo/autossh](https://github.com/islenbo/autossh)
请需要的朋友到原作者去获取，如果原作者有任何问题请联系我(lcl101@163.com)

使用方式：
- cd al
- ./build.sh
- cd ../releases
- unzip al-linux.zip
- cd al-linux
- ./al baidu
就可以登录al.conf文件中定义的baidu服务器，当然也可以通过./al 进行标签选择