package cmd

import (
	// fmt包用于格式化输出
	"fmt"
	// os包提供了操作系统功能的接口
	"os"
	// exec包用于执行外部命令
	"os/exec"
	// filepath包用于处理文件路径
	"path/filepath"
	// strconv包用于字符串和基本数据类型之间的转换
	"strconv"
	// syscall包提供了操作系统底层接口的访问
	"syscall"

	// cobra是一个强大的现代CLI应用程序框架
	"github.com/spf13/cobra"
)

// 定义全局变量，用于存储命令行参数
var (
	// programPath 存储要启动的程序路径
	programPath string
	// programArgs 存储传递给程序的参数列表
	programArgs []string
	// daemonize 标志是否以守护进程方式运行
	daemonize   bool
	// logFile 存储日志文件的路径
	logFile     string
	// pidFile 存储PID文件的路径
	pidFile     string
)

// startCmd 表示启动命令
// 这个命令用于启动指定的程序，可以选择是否以守护进程方式运行
var startCmd = &cobra.Command{
	// 定义命令的使用方式
	Use:   "start [程序路径] [参数...]",
	// 简短描述
	Short: "启动指定的程序",
	// 详细描述
	Long: `start命令用于启动指定的程序，可以选择是否以守护进程方式运行，
并可以指定日志文件和PID文件的路径。`,
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

		// 检查程序是否已经在运行
		if pidFile != "" {
			// 读取PID文件并检查进程是否运行
			if pid, err := readPidFile(pidFile); err == nil && isProcessRunning(pid) {
				fmt.Printf("程序已经在运行，PID: %d\n", pid)
				return
			}
		}

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
	// 将startCmd添加为rootCmd的子命令
	rootCmd.AddCommand(startCmd)

	// 定义命令行标志
	// --daemon/-d 标志：是否以守护进程方式运行
	startCmd.Flags().BoolVarP(&daemonize, "daemon", "d", false, "以守护进程方式运行")
	// --log/-l 标志：指定日志文件路径
	startCmd.Flags().StringVarP(&logFile, "log", "l", "", "日志文件路径")
	// --pid/-p 标志：指定PID文件路径
	startCmd.Flags().StringVarP(&pidFile, "pid", "p", "", "PID文件路径")
}

// startProgram 启动程序并等待其完成
// 这个函数直接在前台启动程序，并将控制权交给该程序
func startProgram() {
	// 创建命令对象，设置程序路径和参数
	cmd := exec.Command(programPath, programArgs...)
	// 默认情况下，继承当前进程的标准输入输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 如果指定了日志文件，重定向输出到日志文件
	if logFile != "" {
		// 打开或创建日志文件，设置为追加写入模式
		logFileHandle, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("无法打开日志文件: %v\n", err)
			os.Exit(1)
		}
		defer logFileHandle.Close()

		// 重定向标准输出和标准错误到日志文件
		cmd.Stdout = logFileHandle
		cmd.Stderr = logFileHandle
	}

	// 启动程序
	err := cmd.Start()
	if err != nil {
		fmt.Printf("启动程序失败: %v\n", err)
		os.Exit(1)
	}

	// 如果指定了PID文件，将进程ID写入PID文件
	if pidFile != "" {
		writePidFile(pidFile, cmd.Process.Pid)
	}

	fmt.Printf("程序已启动，PID: %d\n", cmd.Process.Pid)

	// 等待程序完成执行
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("程序异常退出: %v\n", err)
	}

	// 如果指定了PID文件，程序结束后删除PID文件
	if pidFile != "" {
		os.Remove(pidFile)
	}
}

// startAsDaemon 以守护进程方式启动程序
// 这个函数创建一个分离的后台进程来运行程序
func startAsDaemon() {
	// 创建一个新的进程，该进程将成为守护进程
	// 使用当前程序（procmanager）重新执行start命令，但不带守护进程标志
	cmd := exec.Command(os.Args[0], append([]string{"start", programPath}, programArgs...)...)
	
	// 移除守护进程标志，避免无限递归
	for i, arg := range cmd.Args {
		if arg == "-d" || arg == "--daemon" {
			cmd.Args = append(cmd.Args[:i], cmd.Args[i+1:]...)
			break
		}
	}

	// 设置进程组ID，使子进程不受父进程终止的影响
	// Setsid创建一个新的会话并设置进程组ID
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	// 如果指定了日志文件，重定向输出到日志文件
	if logFile != "" {
		// 打开或创建日志文件
		logFileHandle, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("无法打开日志文件: %v\n", err)
			os.Exit(1)
		}
		defer logFileHandle.Close()

		// 重定向输出到日志文件
		cmd.Stdout = logFileHandle
		cmd.Stderr = logFileHandle
	} else {
		// 如果没有指定日志文件，将输出重定向到/dev/null（丢弃所有输出）
		devNull, err := os.OpenFile("/dev/null", os.O_WRONLY, 0)
		if err != nil {
			fmt.Printf("无法打开/dev/null: %v\n", err)
			os.Exit(1)
		}
		defer devNull.Close()

		// 重定向标准输入输出到/dev/null
		cmd.Stdin = devNull
		cmd.Stdout = devNull
		cmd.Stderr = devNull
	}

	// 启动守护进程
	err := cmd.Start()
	if err != nil {
		fmt.Printf("启动守护进程失败: %v\n", err)
		os.Exit(1)
	}

	// 如果指定了PID文件，写入PID
	if pidFile != "" {
		writePidFile(pidFile, cmd.Process.Pid)
	}

	// 输出成功信息并退出当前进程
	fmt.Printf("程序已在后台启动，PID: %d\n", cmd.Process.Pid)
	os.Exit(0)
}

// readPidFile 从PID文件中读取PID
// 返回读取到的进程ID和可能的错误
func readPidFile(path string) (int, error) {
	// 读取PID文件内容
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	// 将文件内容转换为整数（进程ID）
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}

	return pid, nil
}

// writePidFile 将PID写入PID文件
// 创建PID文件并写入进程ID
func writePidFile(path string, pid int) {
	// 确保目录存在，如果不存在则创建
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	// 将PID转换为字符串并写入文件
	err := os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		fmt.Printf("无法写入PID文件: %v\n", err)
	}
}

// isProcessRunning 检查指定PID的进程是否在运行
// 返回布尔值表示进程是否存在
func isProcessRunning(pid int) bool {
	// 查找指定PID的进程
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// 在Unix系统上，FindProcess总是成功的，所以我们需要发送信号0来检查进程是否存在
	// 信号0不会发送实际信号，只检查进程是否存在
	err = process.Signal(syscall.Signal(0))
	return err == nil
}