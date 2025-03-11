package main

import (
	// 导入cmd包，该包包含了所有的命令行命令实现
	"github.com/user/procmanager/cmd"
)

// main 函数是程序的入口点
// 它调用cmd包中的Execute函数来启动命令行应用程序
func main() {
	// 执行根命令，这将解析命令行参数并运行相应的命令
	cmd.Execute()
}