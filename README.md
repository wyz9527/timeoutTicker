# timeoutTicker
go+redis 定时器

## 适用业务场景
在某一个时间点的延迟任务非常多的情况。比如商家发布了一个优惠活动，活动的有效期是三天后，用户参与活动需要领取一个优惠卡券，那么可能会存在有非常多的领取卡券记录是需要在同一个时间点把状态改成过期。

## 配置
当前值支持json格式的文件
默认情况下，配置文件为 conf/config.json
```javascript
{
  "redis": "redis://localhost:6379/", //redisdns信息
  "Connections": 4, //最大redis链接数
  "tapKey": [//监听的 redis key数组
    {
      "key" :"TEST:queue:try_%d", //监听的key %d为秒级时间戳
      "class" : "Test", //绑定的class名称
      "infoLogFile": "/var/log/info.log", //信息日志文件
      "errorLogFile": "/var/log/error.log" //错误日志文件
    }
  ],
  "pid" : "/var/log/timeoutTicker/unicom.pid" //pid保存文件
}
```

## 安装
```shell
go get -u github.com/wyz9527/timeoutTicker
```

## 示列
模拟场景：将一个未付款订单信息，1天后状态改成关闭
### test.php
```php
<?php
	$key = 'TEST:queue:try_%d';
	$redis = new Redis();
	$redis->connect('127.0.0.1', 6379, 60));
	$orderId = 1;
	$expireTime = time() + 86400;
	//将需要处理的订单id加入到需要监听的redis 队列
	$redis->rpush(sprintf($key,$expireTime), $orderId)
```

### test.go
```go
package main

import (
	"fmt"
	"github.com/wyz9527/timeoutTicker"
)

func main()  {
	//注册执行方法 将方法和class绑定
	timeoutTicker.Register("Test", myTest ) //这里的Test 要和config里面key 里面的class 名称一致

	if err := timeoutTicker.Run(); err != nil{
		fmt.Println( "error:", err )
	}
	timeoutTicker.Close()
}

//具体的业务执行方法
func myTest( queues string) error {
  fmt.Println(queues) //queues 就是监控的rediskey 的具体值
  //拿到队列queues, 获取队列里面的订单id，执行具体业务，将订单状态改成关闭
  //也可以用其他语言处理具体业务，再使用 exec.Command()去调用
	return  nil
}
```

### 编译
```shell
go build
```

### 已守护进程的方式 start
```shell
./test -opt=start -d=true
```

### stop
```shell
./test -opt=stop
```

### restart
```shell
./test -opt=restart
```


