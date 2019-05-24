# timeoutTicker
go+redis 定时器

## 适用业务场景
需要在某一个时间点去处理某项业务。比如一个未付款订单的有效期是3天，3天后要将订单状态改成关闭。

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
test.go
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
  //执行具体业务，
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


