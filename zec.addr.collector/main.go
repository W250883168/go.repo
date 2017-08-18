package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"io/ioutil"

	"gopkg.in/ini.v1"

	"wshua/xutils/xerr"
	"wshua/xutils/xfile"
	"wshua/xutils/xtime"
)

const (
	_CONFIG_FILE = "config.ini"
)

var (
	pConfig = &_AppConfig{
		Keyword:      "mining.authorize",
		LogDir:       "log/",
		TimeInterval: 30,
		WhiteList:    []string{`t1daJrQ4xG3fBBxfNmZL6PG3qz77pieX1fk`, `t1Nr4Gpuh3YKbaLWu9Cmrvq8CFA1iTJD4Va`}}

	gCount    = 0
	gBlackmap = map[string]string{}
)

type _AppConfig struct {
	Keyword      string
	LogDir       string
	TimeInterval int
	WhiteList    []string
}

func main() {
	fmt.Println("Hello World!")
	for {
		foo()
		time.Sleep(time.Second * time.Duration(pConfig.TimeInterval))
	}

}

func foo() {
	defer xerr.CatchPanic()
	dir := pConfig.LogDir
	dir_list, err := ioutil.ReadDir(dir)
	xerr.ThrowPanic(err)

	gCount = 0
	for i, v := range dir_list {
		fname := v.Name()
		if strings.HasSuffix(fname, `_log.txt`) {
			fmt.Println(i, "=", v.Name())
			proc(dir + v.Name())
		}
	}

	fmt.Printf("Black.List.Size=%d\n", len(gBlackmap))
	log.Printf("Black.List.Size=%d\n", len(gBlackmap))
	for _, v := range gBlackmap {
		fmt.Println(v)
		log.Println(v)
	}

	fmt.Println("")
	log.Println("")
}

func proc(fname string) {
	defer xerr.CatchPanic()
	f, err := os.Open(fname) //打开文件
	xerr.ThrowPanic(err)
	defer f.Close() //打开文件出错处理

	if nil == err {
		buff := bufio.NewReader(f) //读入缓存

	Next:
		for {
			line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
			if err != nil || io.EOF == err {
				break
			}

			if yes := strings.Contains(line, pConfig.Keyword); !yes {
				continue
			}

			for _, v := range pConfig.WhiteList {
				if yes := strings.Contains(line, v); yes {
					continue Next
				}
			}

			//["t1W9HL5Aep6WHsSqHiP9YrjTH2ZpfKR1d3t","x"]}
			gCount++
			fmt.Print(gCount, "	", line) //可以对一行进行处理
			idxBegin := strings.LastIndex(line, `["`) + len(`["`)
			idxEnd := strings.LastIndex(line, `",`)
			wchars := []rune(line)
			addr := string(wchars[idxBegin:idxEnd])

			gBlackmap[addr] = addr
		}

	}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		txt := fmt.Sprintf("初始化： %s /%d\n", fun.Name(), line)
		log.Println(txt)
	}

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	filepath := time.Now().Format(xtime.FORMAT_yyyyMMddHHmmss) + ".log"
	if pLogFile, ok := xfile.CreateFile(filepath); ok {
		log.SetOutput(pLogFile)
	}

	// 加载配置
	file, err := ini.Load(_CONFIG_FILE)
	xerr.ThrowPanic(err)
	section := file.Section("app.config")
	pConfig.Keyword = section.Key(`keyword`).String()
	pConfig.LogDir = section.Key(`log.dir`).String()
	pConfig.TimeInterval, _ = section.Key(`time.interval.second`).Int()
	pConfig.WhiteList = section.Key(`address.white.list`).Strings(",")

	data, _ := json.Marshal(pConfig)
	log.Println(string(data))

}
