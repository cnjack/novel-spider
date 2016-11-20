## novel-spider v1.0
曾经一度痴迷于看小说，小说给了我第二个世界，脱离于现实生活的世界。虽然因为看小说耽误了很多光阴，但是它也给我带来了很多的快乐，伴随主人公的开心而开心，伴随主人公的伤心而伤心。  
只是现在看小说相比之前更为麻烦了，各种书荒，各种收费亦或者各种广告，所以萌生了开发novel的想法，旨在不侵犯他的前提下满足自己开开心心看小说的愿望。  

### auth
admin@nightc.com

### demo
演示地址： [novel](http://novel.nightc.com) 

### TODO v1.1
 - 配合[novel-view](http://git.oschina.net/cnjack/novel-view)和[novel-mobile](http://git.oschina.net/cnjack/novel-mobile)进行优化
 - task优化 多次更新
 - sql表优化
 - 错误处理完善

### how to run
```
go get -u git.oschina.net/cnjack/novel-spider
go build git.oschina.net/cnjack/novel-spider
cp $GOPATH/src/git.oschina.net/cnjack/novel-spider/config.ini .
edit config.ini
(linux)
cp $GOPATH/src/git.oschina.net/cnjack/novel-spider/control .
mkdir var
./control start
(win)
运行novel-spider.exe
```