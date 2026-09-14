package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
	query := c.Query("query")
	fmt.Printf("收到查询请求: %s\n", query)

	// 1. 设置 SSE 必需的响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	// 允许跨域（如果是前后端分离项目需要加上）
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// 2. 模拟大模型生成的文本内容
	content := "你好，我是大模型，正在为你生成回答，请稍候..."
		
	// 3. 逐字推送数据
	for i, char := range content {
		// 检查客户端是否已经断开连接（防止服务端继续无效写入）
		if c.Request.Context().Err() != nil {
			fmt.Println("客户端已断开连接")
			return
		}

		// 构造 SSE 标准格式：data: 内容\n\n
		// 注意：每条消息必须以两个换行符结尾
		msg := fmt.Sprintf("data: %c\n\n", char)
			
		// 写入响应体
		c.Writer.WriteString(msg)
		// 刷新缓冲区，确保数据立即发送到客户端
		c.Writer.Flush()

		// 模拟生成延迟（500毫秒）
		time.Sleep(500 * time.Millisecond)
			
		// 打印日志方便观察
		fmt.Printf("推送第 %d 个字符: %c\n", i+1, char)
		}

		// 4. 发送结束信号（可选，通常用特定的 data 标识结束）
		c.Writer.WriteString("data: [DONE]\n\n")
		c.Writer.Flush()
}