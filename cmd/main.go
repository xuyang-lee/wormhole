package main

import (
	"github.com/xuyang-lee/wormhole/hole"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	hole.Init()

	// 关闭: 捕获 Ctrl+C/SIGINT/SIGTERM 信号（允许 Docker/K8S 优雅下线）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞，直到信号到来

	log.Println("正在关闭服务...")

}
