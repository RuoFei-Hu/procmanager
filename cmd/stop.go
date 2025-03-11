package cmd

import (
	// fmt包用于格式化输出
	"fmt"
	// os包提供了操作系统功能的接口
	"os"
	// syscall包提供了操作系统底层接口的访问
	"syscall"
	// time包提供了时间相关的功能
	"time"

	// cobra是一个强大的现代CLI应用程序框架
	"github.com/spf13/cobra"
)

var (
	// force 标志是否强制终止程序（发送SIGKILL信号）
	force bool
	// timeout 等待程序终止的超时时间（秒）
	timeout int
)

// stopCmd 表示停止命令
// 这个命令用于停止正在运行的程序，需要提供PID文件路径
var stopCmd = &cobra.Command{
	// 定义命令的使用方式
	Use:   "stop",
	// 简短描述
	Short: "停止正在运行的程序",
	// 详细描述
	Long: `stop命令用于停止指定的程序，
需要提供PID文件路径。可以选择是否强制终止程序。`,
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

		// 查找指定PID的进程
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("无法找到进程 (PID: %d): %v\n", pid, err)
			// 清理过时的PID文件
			os.Remove(pidFile) 
			os.Exit(1)
		}

		// 如果不是强制终止，先尝试优雅地终止进程
		if !force {
			fmt.Printf("正在停止进程 (PID: %d)...\n", pid)
			// 发送SIGTERM信号，请求进程正常终止
			err = process.Signal(syscall.SIGTERM)
			if err != nil {
				fmt.Printf("无法发送终止信号: %v\n", err)
				os.Exit(1)
			}

			// 等待进程终止，最多等待timeout秒
			for i := 0; i < timeout; i++ {
				time.Sleep(time.Second)
				// 检查进程是否已经终止
				if !isProcessRunning(pid) {
					fmt.Println("程序已停止")
					// 清理PID文件
					os.Remove(pidFile) 
					return
				}
			}

			// 如果超时仍未终止，提示将尝试强制终止
			fmt.Println("程序未在超时时间内停止，尝试强制终止...")
		}

		// 强制终止进程（发送SIGKILL信号）
		fmt.Printf("正在强制终止进程 (PID: %d)...\n", pid)
		err = process.Signal(syscall.SIGKILL)
		if err != nil {
			fmt.Printf("无法强制终止进程: %v\n", err)
			os.Exit(1)
		}

		// 等待进程被强制终止，最多等待5秒
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			// 检查进程是否已经终止
			if !isProcessRunning(pid) {
				fmt.Println("程序已强制终止")
				// 清理PID文件
				os.Remove(pidFile) 
				return
			}
		}

		// 如果仍然无法终止，提示用户手动检查
		fmt.Println("无法终止程序，请手动检查")
	},
}

// init函数在包被导入时自动执行，用于初始化命令
func init() {
	// 将stopCmd添加为rootCmd的子命令
	rootCmd.AddCommand(stopCmd)

	// 定义命令行标志
	// --pid/-p 标志：指定PID文件路径
	stopCmd.Flags().StringVarP(&pidFile, "pid", "p", "", "PID文件路径")
	// --force/-f 标志：是否强制终止程序
	stopCmd.Flags().BoolVarP(&force, "force", "f", false, "强制终止程序（发送SIGKILL信号）")
	// --timeout/-t 标志：等待程序终止的超时时间
	stopCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "等待程序终止的超时时间（秒）")
	// 标记pid标志为必需
	stopCmd.MarkFlagRequired("pid")
}