package ollama

import (
	"testing"
)

func TestOllamaConnection(t *testing.T) {
	// 创建客户端
	client := NewClient("http://192.168.194.90:11434", 360)

	// 测试简单生成
	response, err := client.Generate("gemma3:27b", "你好")
	if err != nil {
		t.Errorf("连接失败: %v", err)
		return
	}

	if response == "" {
		t.Error("响应为空")
		return
	}

	t.Logf("✅ Ollama连接正常")
	t.Logf("完整响应: %s", response)
}
