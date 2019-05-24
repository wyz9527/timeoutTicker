package timeoutTicker

import (
	"fmt"
	"sync"
	"time"
)

type workerFunc func(string) error

var (
	workers map[string]workerFunc
)

func init() {
	workers = make(map[string]workerFunc)
}

func Register(class string, worker workerFunc) {
	workers[class] = worker
}

//执行具体任务
func runWorker( key string, class string, monitor *sync.WaitGroup, quit <-chan bool, infoLogFile *string, errorLogFile *string ) error  {

	monitor.Add(1 )
	go func() {
		defer func() {
			defer monitor.Done()
			if err := recover(); err !=nil{
				errorLog( errorLogFile, err )
			}
		}()


		if workerFunc, ok := workers[class]; ok {
			//创建定时器
			tick := time.NewTicker(time.Second*1)

			for {
				select {
					case <-quit:
						tick.Stop()
						infoLog( infoLogFile, "退出执行 key:", key )
						return
					case <-tick.C:
						getKey := fmt.Sprintf( key,  time.Now().Unix() )
						conn, err := GetConn()
						if err != nil{
							tick.Stop()
							errorLog(errorLogFile, err )
							return
						}
						reply, err := conn.Do("EXISTS", getKey )
						if err != nil {
							errorLog( errorLogFile, err )
							tick.Stop()
							return
						}
						conn.Flush()
						PutConn(conn)
						if reply.(int64) == 1 {
							monitor.Add(1)
							go func() {
								defer func() {
									monitor.Done()
								}()
								infoLog( infoLogFile, getKey, "被执行")
								workerFunc(getKey)
							}()
						}else {
							infoLog( infoLogFile, getKey, "不存在")
						}
				}
			}

		}else{
			panic( fmt.Sprintf("【%s】 方法不存在或未注册", class) )
		}

	}()

	return nil

}