package main

import (
	"log"

	"mcp-log-server/internal/config"
	"mcp-log-server/internal/handlers"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 加载配置
	cfg := config.GetDefaultConfig()

	// 创建MCP服务器
	mcpServer := server.NewMCPServer(
		"Log Command Generator with AI",
		"2.2.0",
		server.WithToolCapabilities(true),
	)

	// 创建处理器
	logHandler := handlers.NewLogCommandHandler(cfg)

	// 注册日志相关工具
	registerLogTools(mcpServer, logHandler)

	// 启动服务器 - 使用stdio方式
	log.Println("Starting Log Command Generator with AI MCP Server v2.2.0...")
	log.Printf("Ollama服务地址: %s", cfg.Ollama.BaseURL)
	log.Printf("默认模型: %s", cfg.Ollama.DefaultModel)
	log.Println("核心功能: 按设备ID和各种参数查询日志")
	log.Println("服务模式: sse")

	sse := server.NewSSEServer(mcpServer)
	// 启动stdio服务器
	if err := sse.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// registerLogTools 注册日志相关工具
func registerLogTools(mcpServer *server.MCPServer, handler *handlers.LogCommandHandler) {
	// 按设备ID和时间范围查询日志 (核心功能)
	mcpServer.AddTool(mcp.NewTool("query_device_logs_by_time",
		mcp.WithDescription("按设备ID和时间范围查询日志，执行命令并使用AI分析结果"),
		mcp.WithString("environment",
			mcp.Description("环境 (dev/test/staging/prod)"),
			mcp.Required(),
		),
		mcp.WithString("log_type",
			mcp.Description("日志类型 (blackhole/oms)"),
			mcp.Required(),
		),
		mcp.WithString("device_id",
			mcp.Description("设备ID"),
			mcp.Required(),
		),
		mcp.WithString("keyword",
			mcp.Description("搜索关键词（可选）"),
		),
		mcp.WithString("lines",
			mcp.Description("查询日志行数（默认2000行）"),
			mcp.DefaultString("2000"),
		),
		mcp.WithString("start_time",
			mcp.Description("开始时间 (格式: 2025-07-24 11:59:38.369)"),
		),
		mcp.WithString("end_time",
			mcp.Description("结束时间 (格式: 2025-07-24 11:59:38.369)"),
		),
		mcp.WithString("model",
			mcp.Description("AI分析模型（可选）"),
		),
	), handler.QueryDeviceLogsByTimeRange)
}
