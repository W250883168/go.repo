package xfile

import (
	"fmt"
	"os"
	"os/exec"
	FilePath "path/filepath"
	"runtime"

	"dev.project/BackEndCode/devserver/commons/xdebug"
)

// 获取当前EXE执行文件路径
func GetExeFilePath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := FilePath.Abs(file)

	return path
}

// 获取当前工作目录
func GetPWD() (dir string) {
	dir, err := FilePath.Abs(FilePath.Dir(os.Args[0]))
	if err != nil {
		xdebug.LogError(err)
	}

	return dir
}

func FileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func CreateFile(fileName string) (file *os.File, ok bool) {
	defer xdebug.DoRecover()

	pFile, err := os.Create(fileName)
	xdebug.HandleError(err)

	file = pFile
	ok = true
	return file, ok
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	fmt.Println(GetExeFilePath())
	fmt.Println(GetPWD())
}
