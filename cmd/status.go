package cmd

import (
	// fmt包用于格式化输出
	"fmt"
	// os包提供了操作系统功能的接口
	"os"

	// cobra是一个强大的现代CLI应用程序框架
	"github.com/spf13/cobra"
)

// statusCmd 表示状态命令
// 这个命令用于查看指定程序的运行状态，需要提供PID文件路径
var statusCmd = &cobra.Command{
	// 定义命令的使用方式
	Use:   "status",
	// 简短描述
	Short: "查看程序的运行状态",
	// 详细描述
	Long: `status命令用于查看指定程序的运行状态，
需要提供PID文件路径。`,
	// 命令执行函数
	Run: func(cmd *cobra.Command, args []string) {
		// 检查是否提供了PID文件路径
		if pidFile == "" {
			fmt.Println("错误: 必须指定PID文件路径")
			os.Exit(1)
		}

		// 检查PID文件是否存在，并读取PID
		pid, err := readPidFile(pidFile)
		if err != nil {
			fmt.Printf("程序未运行: %v\n", err)
			os.Exit(1)
		}

		// 检查进程是否在运行
		if isProcessRunning(pid) {
			// 如果进程正在运行，输出PID信息
			fmt.Printf("程序正在运行，PID: %d\n", pid)
		} else {
			// 如果进程不在运行但PID文件存在，提示并清理PID文件
			fmt.Printf("程序未运行，但PID文件存在 (PID: %d)\n", pid)
			// 清理过时的PID文件
			os.Remove(pidFile)
		}
	},
}

// init函数在包被导入时自动执行，用于初始化命令
func init() {
	// 将statusCmd添加为rootCmd的子命令
	rootCmd.AddCommand(statusCmd)

	// 定义命令行标志
	// --pid/-p 标志：指定PID文件路径
	statusCmd.Flags().StringVarP(&pidFile, "pid", "p", "", "PID文件路径")
	// 标记pid标志为必需
	statusCmd.MarkFlagRequired("pid")
}