package cmd

import (
	// fmt包用于格式化输出
	"fmt"
	// os包提供了操作系统功能的接口
	"os"

	// cobra是一个强大的现代CLI应用程序框架
	"github.com/spf13/cobra"
)

// rootCmd 代表没有调用子命令时的基础命令
// 这是整个命令行应用的根命令，所有其他命令都是它的子命令
var rootCmd = &cobra.Command{
	// Use定义了命令的名称和用法
	Use:   "procmanager",
	// Short是在帮助输出中显示的简短描述
	Short: "一个通用的进程管理工具",
	// Long是在帮助输出中显示的详细描述
	Long: `procmanager 是一个通用的命令行工具，用于管理程序的生命周期。
它提供了启动、重启、查看状态和停止程序的功能。`,
}

// Execute 将所有子命令添加到root命令并适当设置标志。
// 这由main.main()调用。它只需要对rootCmd执行一次。
// 该函数执行根命令，解析命令行参数并运行相应的命令
func Execute() {
	// 执行根命令
	err := rootCmd.Execute()
	// 如果执行过程中出现错误，打印错误信息并以状态码1退出程序
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init函数在包被导入时自动执行，用于初始化
func init() {
	// 在这里，您可以定义标志和配置设置。
	// Cobra支持持久性标志，如果定义了的话，这些标志将会被设置为应用。
	
	// 添加一个全局的配置文件标志，所有子命令都可以访问
	rootCmd.PersistentFlags().StringP("config", "c", "", "配置文件路径")
}