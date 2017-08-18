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
	"unicode"

	"io/ioutil"
	"path/filepath"

	"gopkg.in/ini.v1"

	"xutils/xerr"
	"xutils/xfile"
	"xutils/xtime"
)

const (
	_CONFIG_FILE = "config.ini"
)

var (
	pConfig = &_AppConfig{LogDir: "log/"}
)

type _AppConfig struct {
	Keyword      string
	LogDir       string
	TimeInterval int
	WhiteList    []string
}

func main() {
	fmt.Println("Hello World!")
	dir := pConfig.LogDir
	dir_list, err := ioutil.ReadDir(dir)
	xerr.ThrowPanic(err)

	// work_dir(dir)
	// work_through(dir)

	for _, f := range dir_list {
		if f.IsDir() {
			foo(f.Name() + string(filepath.Separator))
		}
	}
}

func foo(dir string) {
	defer xerr.CatchPanic()
	dir_list, err := ioutil.ReadDir(dir)
	xerr.ThrowPanic(err)

	for i, v := range dir_list {
		fname := v.Name()
		if !v.IsDir() && strings.HasSuffix(fname, `.java`) {
			log.Println(i, "=", v.Name())
			proc(dir + v.Name())
		}
	}

}

func proc(fname string) {
	defer xerr.CatchPanic()
	f, err := os.Open(fname) //打开文件
	xerr.ThrowPanic(err)
	defer f.Close() //打开文件出错处理

	if nil == err {
		buff := bufio.NewReader(f) //读入缓存
		new_fname := strings.TrimSuffix(fname, filepath.Ext(fname))
		log.Println(new_fname)
		pNewFile, err := os.Create(new_fname + `.txt`)
		xerr.ThrowPanic(err)
		defer pNewFile.Close()

		for {
			line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
			if err != nil || io.EOF == err {
				_, err = pNewFile.WriteString(line + "\n")
				xerr.ThrowPanic(err)
				break
			}

			bPrivateString := strings.Contains(line, "private String")
			bPrivateInt := strings.Contains(line, "private Integer")
			bPrivateLong := strings.Contains(line, "private Long")
			line = strings.TrimRightFunc(line, unicode.IsSpace)
			// skip empty line
			if len(line) <= 0 {
				continue
			}

			// replace string
			if bPrivateString && !strings.Contains(line, "=") && strings.HasSuffix(line, `;`) {
				line = strings.Replace(line, `;`, ` = "";`, 1) + "\n"
				_, err = pNewFile.WriteString(line)
				xerr.ThrowPanic(err)
				continue
			}

			// replac int
			if bPrivateInt && !strings.Contains(line, "=") && strings.HasSuffix(line, `;`) {
				line = strings.Replace(line, `;`, ` = 0;`, 1) + "\n"
				_, err = pNewFile.WriteString(line)
				xerr.ThrowPanic(err)
				continue
			}

			// replace long
			if bPrivateLong && !strings.Contains(line, "=") && strings.HasSuffix(line, `;`) {
				line = strings.Replace(line, `;`, ` = 0L;`, 1) + "\n"
				_, err = pNewFile.WriteString(line)
				xerr.ThrowPanic(err)
				continue
			}

			_, err = pNewFile.WriteString(line + "\n")
		}

		pNewFile.Close()
	}
}

func work_through(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		println(path)
		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func work_dir(dir string) {
	defer xerr.CatchPanic()
	log.Println(dir)
	files, err := ioutil.ReadDir(dir)
	xerr.ThrowPanic(err)

	for _, finfo := range files {
		fpath := filepath.Join(dir, finfo.Name())
		if !finfo.IsDir() && strings.HasSuffix(fpath, ".jpg") {
			log.Println("\t File: ", fpath)
			dest_dir := finfo.ModTime().Format(xtime.FORMAT_yyyyMM)
			if !xfile.FileExist(dest_dir) {
				os.MkdirAll(dest_dir, os.ModeDir)
			}

			pSrcFile, err := os.Open(fpath)
			xerr.ThrowPanic(err)
			dest_file := filepath.Join(dest_dir, finfo.Name())
			if pDestFile, ok := xfile.CreateFile(dest_file); ok {
				_, err = io.Copy(pDestFile, pSrcFile)
				xerr.ThrowPanic(err)
				pDestFile.Close()

				pSrcFile.Close()
				os.Remove(fpath)
			}
		} else if finfo.IsDir() {
			work_dir(fpath)
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
	//	filepath := time.Now().Format(xtime.FORMAT_yyyyMMddHHmmss) + ".log"
	//	if pLogFile, ok := xfile.CreateFile(filepath); ok {
	//		log.SetOutput(pLogFile)
	//	}

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
