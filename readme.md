# SsrMicroClient  
[![](https://img.shields.io/github/license/asutorufa/ssrmicroclient.svg)](https://raw.githubusercontent.com/Asutorufa/SsrMicroClient/master/LICENSE)
[![](https://img.shields.io/github/release-pre/asutorufa/ssrmicroclient.svg)](https://github.com/Asutorufa/SsrMicroClient/releases)
[![codebeat badge](https://codebeat.co/badges/ce94a347-64b1-4ee3-9b18-b95858e1c6b4)](https://codebeat.co/projects/github-com-asutorufa-ssrmicroclient-master)
![](https://img.shields.io/github/languages/top/asutorufa/ssrmicroclient.svg)  

How to use:
- download the [releases](https://github.com/Asutorufa/SsrMicroClient/releases) binary file.if not have your platform ,please build it by yourself.
- if you use windows,you need to read [how to install libsodium to windows](https://github.com/Asutorufa/SsrMicroClient/blob/master/windows_use_ssr_python.md).
- build

```
git clone https://github.com/Asutorufa/SsrMicroClient.git
go get -u github.com/mattn/go-sqlite3
cd SsrMicroClient
go build SSRSub.go
./SSRSub
```
- because the ssr_python deamon not support windows,so i use vgs to make ssr run in deamon and wirte pid to file,but the windows cmd a little slow to get process pid,so i set a 500ms wait to get pid.
- config file  
  it will auto create at first run,path at `~/.config/SSRSub`,windows at Documents/SSRSub.
<!--
```
#config path at ~/.config/SSRSub
#config file,first run auto create,# to note
#python_path /usr/bin/python3
#ssr_path /shadowsocksr-python/shadowsocks/local.py
#local_port 1080
#local_address 127.0.0.1
#connect-verbose-info
workers 8
fast-open
deamon
#pid-file /home/xxx/.config/SSRSub/shadowsocksr.pid
#log-file /dev/null
```
-->
![](https://raw.githubusercontent.com/Asutorufa/SsrMicroClient/master/img/SSRSubV0.2.2beta.png)
<!--
issue:
- [ ] now only can run in bash,cmd is not test.
- [ ] not test path exist or not(now everything is normal).
-->
[日本語](https://github.com/Asutorufa/SSRSubscriptionDecode/blob/master/readme_jp.md) [中文](https://github.com/Asutorufa/SSRSubscriptionDecode/blob/master/readme_cn.md) [other progrmammer language vision](https://github.com/Asutorufa/SSRSubscriptionDecode/blob/master/readme_others.md)   

# Thanks
[Golang](https://golang.org)  
[mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)  
[breakwa11/shadowsokcsr](https://github.com/shadowsocksr-backup/shadowsocksr)  
[akkariiin/shadowsocksrr](https://github.com/shadowsocksrr/shadowsocksr/tree/akkariiin/dev)  

# already know issue
ssr python version at mac may be not support,please test by yourself.

# Others
Make a simple gui([Now Dev](https://github.com/Asutorufa/SsrMicroClient/tree/dev)):
![](https://raw.githubusercontent.com/Asutorufa/SsrMicroClient/dev/img/gui_dev.png)  
Todo:
- [x] (give up)use shadowsocksr write by golang(sun8911879/shadowsocksR),or use ssr_libev share libraries.  
      write a half of [http proxy](https://github.com/Asutorufa/SsrMicroClient/blob/OtherLanguage/Old/SSR_http_client/client.go) find sun8911879/shadowsocksR is not support auth_chain*...oof.  
      when i use ssr_libev i cant run it in the golang that has so many error,i fix a little but more and more error appear. 
```
      # command-line-arguments
    /tmp/go-build379176400/b001/_x002.o：在函数‘main’中：
    ./local.c:1478: `main'被多次定义
    # command-line-arguments
    .........
    .........
    .........
    ./local.c:438:36: warning: comparison between pointer and       integer
                         if (perror == EINPROGRESS) {
                                    ^~
``` 
- [x] (give up)add bypass(↑↑↑ i cant run ssr in golang ↑↑↑)
- [x] ss link compatible. 
- [ ] need more ss link template.
- [ ] support http proxy.
- [ ] create shortcut at first run,auto move or copy file to config path.
- [ ] add `-h` argument to show help.


fixed issue:
- process android is not linux.
- sh should use which to get.  
- support windows.
- can setting timeout.