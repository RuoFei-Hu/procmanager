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

// restartCmd 表示重启命令
// 这个命令用于重启指定的程序，先停止正在运行的程序，然后重新启动它
var restartCmd = &cobra.Command{
	// 定义命令的使用方式
	Use:   "restart [程序路径] [参数...]",
	// 简短描述
	Short: "重启指定的程序",
	// 详细描述
	Long: `restart命令用于重启指定的程序，
先停止正在运行的程序，然后重新启动它。`,
	// 至少需要一个参数（程序路径）
	Args: cobra.MinimumNArgs(1),
	// 命令执行函数
	Run: func(cmd *cobra.Command, args []string) {
		// 获取程序路径（第一个参数）
		programPath = args[0]
		// 如果有更多参数，则存储为程序参数
		if len(args) > 1 {
			programArgs = args[1:]
		}

		// 检查是否提供了PID文件路径
		if pidFile == "" {
			fmt.Println("错误: 必须指定PID文件路径")
			os.Exit(1)
		}

		// 检查程序是否在运行
		pid, err := readPidFile(pidFile)
		if err == nil && isProcessRunning(pid) {
			// 如果程序正在运行，先停止它
			fmt.Printf("正在停止进程 (PID: %d)...\n", pid)
			process, _ := os.FindProcess(pid)
			
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
					break
				}
				
				// 如果超时，强制终止
				if i == timeout-1 {
					fmt.Println("程序未在超时时间内停止，尝试强制终止...")
					// 发送SIGKILL信号强制终止进程
					process.Signal(syscall.SIGKILL)
					time.Sleep(time.Second)
					// 再次检查进程是否已经终止
					if !isProcessRunning(pid) {
						fmt.Println("程序已强制终止")
						// 清理PID文件
						os.Remove(pidFile) 
					} else {
						// 如果仍然无法终止，提示用户手动检查
						fmt.Println("无法终止程序，请手动检查")
						os.Exit(1)
					}
				}
			}
		}

		// 启动程序
		fmt.Println("正在启动程序...")
		// 根据是否守护进程模式选择启动方式
		if daemonize {
			startAsDaemon()
		} else {
			startProgram()
		}
	},
}

// init函数在包被导入时自动执行，用于初始化命令
func init() {
	// 将restartCmd添加为rootCmd的子命令
	rootCmd.AddCommand(restartCmd)

	// 定义命令行标志
	// --daemon/-d 标志：是否以守护进程方式运行
	restartCmd.Flags().BoolVarP(&daemonize, "daemon", "d", false, "以守护进程方式运行")
	// --log/-l 标志：指定日志文件路径
	restartCmd.Flags().StringVarP(&logFile, "log", "l", "", "日志文件路径")
	// --pid/-p 标志：指定PID文件路径
	restartCmd.Flags().StringVarP(&pidFile, "pid", "p", "", "PID文件路径")
	// --timeout/-t 标志：等待程序终止的超时时间
	restartCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "等待程序终止的超时时间（秒）")
	// 标记pid标志为必需
	restartCmd.MarkFlagRequired("pid")
}