## novel-spider v2.0
曾经一度痴迷于看小说，小说给了我第二个世界，脱离于现实生活的世界。虽然因为看小说耽误了很多光阴，但是它也给我带来了很多的快乐，伴随主人公的开心而开心，伴随主人公的伤心而伤心。  
只是现在看小说相比之前更为麻烦了，各种书荒，各种收费亦或者各种广告，所以萌生了开发novel的想法，旨在不侵犯他的前提下满足自己开开心心看小说的愿望。  

### auth
admin@nightc.com

### demo
演示地址： [novel](https://novel.nightc.com) 
### depend on
 - mysql
 - redis
### app && web
you can build web(single web app) or app(cordova) from [novel-mobile](https://gitee.com/cnjack/novel-mobile)
### how to run

```
go get -u gitee.com/cnjack/novel-spider
go build gitee.com/cnjack/novel-spider
cp $GOPATH/src/gitee.com/cnjack/novel-spider/config.yaml .
edit config.yaml
(linux)
cp $GOPATH/src/gitee.com/cnjack/novel-spider/control .
mkdir var
./control start
(win)
运行novel-spider.exe
```