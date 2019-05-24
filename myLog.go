package timeoutTicker

import (
	"fmt"
	"log"
	"os"
	"time"
)


func init()  {
	//设置日志的格式
	log.SetPrefix("【timeoutTicker】") //设置日志的前缀
	log.SetFlags(log.Ldate|log.Lshortfile) //显示日期，文件名和行号
}

func InfoLog( myLogFile *string, msg ...interface{} ) {
	infoLog( myLogFile, msg... )
}

//记录信息日志
func infoLog( myLogFile *string, msg ...interface{} ){

	t := time.Now()
	y,m,d := t.Date()
	ymd := fmt.Sprintf("%d%d%d", y,m,d)

	errFile,err:=os.OpenFile( *myLogFile+"."+ymd, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	defer errFile.Close()
	if err!=nil{
		log.Fatalln("打开日志文件失败：",err)
	}

	myInfo := log.New(errFile,"Info:\n",log.Ldate | log.Ltime | log.Lshortfile)
	myInfo.Println( msg... )
}

//记录错误日志
func errorLog( myLogFile *string, msg ...interface{} ){
	t := time.Now()
	y,m,d := t.Date()
	ymd := fmt.Sprintf("%d%d%d", y,m,d)


	errFile,err:=os.OpenFile( *myLogFile+"."+ymd, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	defer errFile.Close()
	if err!=nil{
		log.Fatalln("打开日志文件失败：",err)
	}

	myError := log.New(errFile,"Error:",log.Ldate | log.Ltime | log.Lshortfile)
	myError.Println( msg... )

}

//记录警告日志
func warnLog( myLogFile *string, msg ...interface{})  {
	t := time.Now()
	y,m,d := t.Date()
	ymd := fmt.Sprintf("%d%d%d", y,m,d)

	errFile,err:=os.OpenFile( *myLogFile+"."+ymd, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	defer errFile.Close()
	if err!=nil{
		log.Fatalln("打开日志文件失败：",err)
	}

	myWarn := log.New(errFile,"Warin:",log.Ldate | log.Ltime | log.Lshortfile)
	myWarn.Println( msg... )
}