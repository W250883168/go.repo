package xfile

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	FilePath "path/filepath"

	xutil "go.repo/xutils/xapp"
	"go.repo/xutils/xdebug"
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

// 文件是否存在
func FileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

// 创建文件
func CreateFile(fileName string) (file *os.File, ok bool) {
	if pFile, err := os.Create(fileName); err == nil {
		file, ok = pFile, true
	}

	return file, ok
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok && xutil.IsDebugMode() {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
