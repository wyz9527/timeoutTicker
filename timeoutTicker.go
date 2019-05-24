package timeoutTicker

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/spf13/viper"
	"github.com/youtube/vitess/go/pools"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var (
	logger      seelog.LoggerInterface
	pool        *pools.ResourcePool
	ctx         context.Context
	initMutex   sync.Mutex
	initialized bool
	stop string
	restart string
)

func Init() error {
	initMutex.Lock()
	defer initMutex.Unlock()
	if !initialized {
		pid := os.Getpid()
		err := WritePid( viper.GetString("pid"), pid )
		if err != nil {
			fmt.Println("Error:", err )
			return  err
		}

		ctx = context.Background()

		pool = newRedisPool(viper.GetString("redis"), viper.GetInt("Connections"), viper.GetInt("Connections"), time.Minute )

		initialized = true
	}

	return nil
}

func GetConn() (*RedisConn, error) {
	resource, err := pool.Get(ctx)

	if err != nil {
		return nil, err
	}
	return resource.(*RedisConn), nil
}

func PutConn(conn *RedisConn) {
	pool.Put(conn)
}


func Close() {
	initMutex.Lock()
	defer initMutex.Unlock()
	if initialized {
		pool.Close()
		initialized = false
	}
}

//写文件
func WritePid(name string, pid int) error {
	return  ioutil.WriteFile(name, []byte(fmt.Sprintln(pid)),0666)
}

func Run() error {
	Init()

	quit := signals()

	var quits []chan bool
	i := 0
	var monitor sync.WaitGroup
	keys := viper.Get("tapKey")
	if ks, ok := keys.([]interface{}); ok{
		for _, v := range ks{
			if tk, ok := v.(map[string]interface{}); ok{
				quits = append(quits, make(chan bool))
				infolog := tk["infoLogFile"].(string)
				errorlog := tk["errorLogFile"].(string)
				err := runWorker( tk["key"].(string), tk["class"].(string), &monitor, quits[i], &infolog, &errorlog )
				if err != nil{
					return err
				}
				i++
			}
		}
	}

	monitor.Add(1)
	go func() {
		defer func() {
			monitor.Done()
		}()
		for{
			select {
				case <-quit:
					for _, v := range quits {
						v <- true
					}
					return
			}
		}
	}()
	monitor.Wait()

	return nil
}

