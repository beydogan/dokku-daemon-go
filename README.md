我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

# dokku-daemon-go
Dokku Daemon written with Go to interact with Dokku

# Requirements 

A server running Ubuntu 14.04 or later with Dokku installed.

# Installing 

As a user with access to `sudo`

```
wget https://github.com/beydogan/dokku-daemon-go/releases/download/v0.1.1/dokku-daemon-go-linux64
chmod +x dokku-daemon-go-linux64
sudo ./dokku-daemon-go-linux64 install
sudo service dokku-daemon start
```

# Usage

Daemon will create a UNIX socket at `/var/run/dokku-daemon/dokku-daemon.sock` owned by `dokku` user.

To test you can use;

```
sudo socat - UNIX-CONNECT:/var/run/dokku-daemon/dokku-daemon.sock
```

and type `apps`. It should return the output in JSON format.

```
{"status":"success","output":{"message":"=====\u003e My Apps\nhello\n"}}
```
