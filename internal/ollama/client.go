package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Request Ollama API请求结构
type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Response Ollama API响应结构
type Response struct {
	Model      string `json:"model"`
	CreatedAt  string `json:"created_at"`
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	DoneReason string `json:"done_reason,omitempty"`
}

// Client Ollama客户端
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient 创建新的Ollama客户端
func NewClient(baseURL string, timeoutSeconds int) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

// Generate 调用Ollama生成API - 使用非流式响应
func (c *Client) Generate(model, prompt string) (string, error) {
	// 构建请求体 - 使用非流式响应
	reqBody := Request{
		Model:  model,
		Prompt: prompt,
		Stream: false, // 改为非流式响应
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %v", err)
	}

	// 发送HTTP请求
	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/api/generate", c.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 处理非流式响应 - 直接读取完整响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var ollamaResp Response
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("JSON解码失败: %v", err)
	}

	return ollamaResp.Response, nil
}

// AnalyzeLog 专门用于日志分析的方法
func (c *Client) AnalyzeLog(model, logContent, issueDescription string) (string, error) {
	prompt := fmt.Sprintf(`你是一个专业的运维工程师，请分析以下日志内容：

日志内容:
%s

📊 **分析要求:**
请按以下结构进行详细分析：

## 🔍 异常识别与分类
- **错误日志**: 列出所有ERROR级别的日志，包括错误代码、时间戳和描述
- **警告信息**: 统计WARN级别的警告，分析潜在风险
- **异常模式**: 识别重复出现的异常或错误趋势

## 📈 系统状态评估  
- **设备运行状态**: 分析各设备ID的运行情况
- **性能指标**: 提取CPU、内存、电压、电流等关键指标
- **通信状态**: 评估设备间通信和外部服务连接状态

## ⚠️ 风险评估
- **紧急程度**: 按严重程度对问题进行分级（高/中/低）
- **影响范围**: 分析问题可能影响的系统模块或业务功能
- **建议措施**: 针对发现的问题提供具体的处理建议

## 📋 日志统计摘要
- **时间范围**: 日志的时间跨度
- **日志总数**: 各级别日志的数量统计
- **设备清单**: 涉及的所有设备ID列表

请用中文回答，格式清晰易读。`, logContent)

	if issueDescription != "" {
		prompt += fmt.Sprintf("\n\n特别关注的问题：%s", issueDescription)
	}

	return c.Generate(model, prompt)
}

// AskQuestion 通用问答方法
func (c *Client) AskQuestion(model, question, context string) (string, error) {
	prompt := question
	if context != "" {
		prompt = fmt.Sprintf("上下文信息:\n%s\n\n问题: %s", context, question)
	}

	return c.Generate(model, prompt)
}
