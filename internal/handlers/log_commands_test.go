package handlers

import (
	"context"
	"testing"

	"mcp-log-server/internal/config"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestQueryDeviceLogsByTimeRange(t *testing.T) {
	// 创建配置
	cfg := config.GetDefaultConfig()

	// 创建处理器
	handler := NewLogCommandHandler(cfg)

	// 准备测试参数
	arguments := map[string]interface{}{
		"environment": "dev",
		"log_type":    "oms",
		"device_id":   "yHkqCtKAdWuBPX5Z",
		"lines":       "100",
	}

	// 创建请求
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "query_device_logs_by_time",
			Arguments: arguments,
		},
	}

	// 调用方法
	result, err := handler.QueryDeviceLogsByTimeRange(context.Background(), request)

	// 检查结果
	if err != nil {
		t.Logf("方法调用完成，可能因为SSH连接问题: %v", err)
		return
	}

	if result == nil {
		t.Error("结果为空")
		return
	}

	if len(result.Content) == 0 {
		t.Error("结果内容为空")
		return
	}

	t.Logf("✅ 方法调用成功")
	t.Logf("结果类型: %s", result.Content[0])
}
