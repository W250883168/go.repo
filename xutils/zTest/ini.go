package test

import (
	"fmt"
	"log"
	"runtime"

	ini "gopkg.in/ini.v1"
)

const iniFile string = "config.ini"

var (
	server string = "127.0.0.1"
	port   string = "3306"
	dbname string = "zndx"
)

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	if file, err := ini.Load(iniFile); err == nil {
		section := file.Section("database")
		server = section.Key("server").Value()
		port = section.Key("port").Value()
		dbname = section.Key("dbname").Value()
		log.Printf("server=%s; port=%s; dbname=%s\n", server, port, dbname)
	}
}
