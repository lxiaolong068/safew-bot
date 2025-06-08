package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"safew-bot/bot"
)

// 版本信息 (由编译时注入)
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// 命令行参数
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "显示版本信息")
	flag.BoolVar(&showVersion, "v", false, "显示版本信息 (简写)")
	flag.Parse()

	// 显示版本信息
	if showVersion {
		fmt.Printf("SafeW Bot\n")
		fmt.Printf("版本: %s\n", Version)
		fmt.Printf("构建时间: %s\n", BuildTime)
		return
	}

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 打印启动信息
	log.Printf("SafeW Bot v%s (构建时间: %s)", Version, BuildTime)

	// 加载配置
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	log.Printf("配置加载成功 - 日志级别: %s, 轮询超时: %d秒", config.LogLevel, config.PollTimeout)

	// 创建Bot实例
	safewBot := bot.NewBot(config.BotToken, config.PollTimeout)

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置优雅关闭信号处理
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		
		sig := <-sigChan
		log.Printf("接收到信号 %v，正在关闭Bot...", sig)
		
		// 取消上下文，触发Bot停止
		cancel()
	}()

	// 启动Bot
	log.Println("正在启动SafeW Bot...")
	if err := safewBot.Start(ctx); err != nil {
		if err == context.Canceled {
			log.Println("Bot已优雅关闭")
		} else {
			log.Fatalf("Bot运行时出错: %v", err)
		}
	}

	log.Println("SafeW Bot已停止")
} 