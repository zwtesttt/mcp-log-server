package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"mcp-log-server/internal/config"
	"mcp-log-server/internal/ollama"

	"github.com/mark3labs/mcp-go/mcp"
)

// 安全常量定义
const (
	MAX_DEVICE_ID_LENGTH = 50
	MAX_KEYWORD_LENGTH   = 100
	MAX_LINES_LIMIT      = 99999
	COMMAND_TIMEOUT      = 30 * time.Second
)

// LogCommandHandler 日志命令处理器
type LogCommandHandler struct {
	config       *config.Config
	ollamaClient *ollama.Client
}

// NewLogCommandHandler 创建日志命令处理器
func NewLogCommandHandler(cfg *config.Config) *LogCommandHandler {
	ollamaClient := ollama.NewClient(cfg.Ollama.BaseURL, cfg.Ollama.Timeout)
	return &LogCommandHandler{
		config:       cfg,
		ollamaClient: ollamaClient,
	}
}

// QueryDeviceLogsByTimeRange 按设备ID和时间范围查询日志并进行AI分析
func (h *LogCommandHandler) QueryDeviceLogsByTimeRange(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()

	// 获取环境
	envName, ok := arguments["environment"].(string)
	if !ok {
		return nil, fmt.Errorf("环境参数无效")
	}

	env, exists := h.config.GetEnvironment(envName)
	if !exists {
		return nil, fmt.Errorf("未知环境: %s", envName)
	}

	// 获取日志类型
	logTypeName, ok := arguments["log_type"].(string)
	if !ok {
		return nil, fmt.Errorf("日志类型参数无效")
	}

	logFile, exists := h.config.GetLogFileForEnvironment(logTypeName, envName)
	if !exists {
		return nil, fmt.Errorf("未知日志类型: %s", logTypeName)
	}

	// 获取设备ID
	deviceID, ok := arguments["device_id"].(string)
	if !ok || deviceID == "" {
		return nil, fmt.Errorf("设备ID参数无效")
	}

	// 验证设备ID安全性
	if err := validateDeviceID(deviceID); err != nil {
		return nil, fmt.Errorf("设备ID验证失败: %v", err)
	}

	// 获取关键词（可选）
	keyword, _ := arguments["keyword"].(string)
	if keyword != "" {
		// 验证关键词长度
		if len(keyword) > MAX_KEYWORD_LENGTH {
			return nil, fmt.Errorf("关键词长度不能超过%d个字符", MAX_KEYWORD_LENGTH)
		}
		// 转义关键词中的危险字符
		keyword = escapeShellArg(keyword)
	}

	// 获取日志行数（默认2000行）
	lines := "2000"
	if l, ok := arguments["lines"].(string); ok && l != "" {
		// 验证行数是否为有效数字且在合理范围内
		if matched, _ := regexp.MatchString(`^\d+$`, l); !matched {
			return nil, fmt.Errorf("日志行数必须为数字")
		}
		// 转换为整数进行范围检查
		var linesInt int
		if _, err := fmt.Sscanf(l, "%d", &linesInt); err != nil || linesInt <= 0 || linesInt > MAX_LINES_LIMIT {
			return nil, fmt.Errorf("日志行数必须在1到%d之间", MAX_LINES_LIMIT)
		}
		lines = l
	}

	// 获取开始时间和结束时间（可选）
	startTime, hasStartTime := arguments["start_time"].(string)
	endTime, hasEndTime := arguments["end_time"].(string)

	// 验证时间格式（如果提供）
	timePattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{3})?$`
	if hasStartTime {
		if matched, _ := regexp.MatchString(timePattern, startTime); !matched {
			return nil, fmt.Errorf("开始时间格式无效，应为: 2025-07-24 11:59:38.369")
		}
	}
	if hasEndTime {
		if matched, _ := regexp.MatchString(timePattern, endTime); !matched {
			return nil, fmt.Errorf("结束时间格式无效，应为: 2025-07-24 11:59:38.369")
		}
	}

	// 获取AI模型
	model := h.config.Ollama.DefaultModel
	if m, ok := arguments["model"].(string); ok && m != "" {
		model = m
	}

	// 转换行数为整数
	linesInt, _ := strconv.Atoi(lines)

	// 直接读取本地日志文件 - 安全且高效
	var timeRangeStart, timeRangeEnd string
	if hasStartTime && hasEndTime {
		timeRangeStart = startTime
		timeRangeEnd = endTime
	}

	logContent, err := h.readLocalLogFile(logFile.Path, deviceID, keyword, linesInt, timeRangeStart, timeRangeEnd)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("❌ **读取日志文件失败**\n\n**日志文件**: %s\n**设备ID**: %s\n**错误**: %s\n\n**建议**:\n- 检查日志文件是否存在\n- 确认文件读取权限\n- 验证设备ID格式是否正确",
						logFile.Path, deviceID, err.Error()),
				},
			},
		}, nil
	}

	// 构建执行的"命令"描述（用于显示）
	var commandDesc string
	if hasStartTime && hasEndTime {
		if keyword != "" {
			commandDesc = fmt.Sprintf("读取文件 %s，过滤设备ID '%s'、关键词 '%s'、时间范围 %s 到 %s，最多 %s 行",
				logFile.Path, deviceID, keyword, startTime, endTime, lines)
		} else {
			commandDesc = fmt.Sprintf("读取文件 %s，过滤设备ID '%s'、时间范围 %s 到 %s，最多 %s 行",
				logFile.Path, deviceID, startTime, endTime, lines)
		}
	} else {
		if keyword != "" {
			commandDesc = fmt.Sprintf("读取文件 %s，过滤设备ID '%s'、关键词 '%s'，最近 %s 行",
				logFile.Path, deviceID, keyword, lines)
		} else {
			commandDesc = fmt.Sprintf("读取文件 %s，过滤设备ID '%s'，最近 %s 行",
				logFile.Path, deviceID, lines)
		}
	}

	// 检查日志内容是否为空
	if strings.TrimSpace(logContent) == "" {
		timeRangeDesc := fmt.Sprintf("最近%s条", lines)
		if hasStartTime && hasEndTime {
			timeRangeDesc = fmt.Sprintf("%s 到 %s", startTime, endTime)
		}

		keywordDesc := ""
		if keyword != "" {
			keywordDesc = fmt.Sprintf("，关键词: %s", keyword)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("⚠️ **未找到相关日志**\n\n**查询条件**:\n- 环境: %s (%s)\n- 日志类型: %s\n- 设备ID: %s\n- 时间范围: %s%s\n\n**执行的操作**: %s\n\n**可能原因**:\n- 该设备ID在指定条件下无日志记录\n- 设备ID或关键词格式不正确\n- 日志文件路径错误\n- 权限不足无法读取日志文件",
						env.Name, envName, logFile.Description, deviceID, timeRangeDesc, keywordDesc, commandDesc),
				},
			},
		}, nil
	}

	// 使用AI分析日志内容
	issueDescription := fmt.Sprintf("设备ID: %s", deviceID)
	if keyword != "" {
		issueDescription += fmt.Sprintf(", 关键词: %s", keyword)
	}
	if hasStartTime && hasEndTime {
		issueDescription += fmt.Sprintf(", 时间范围: %s 到 %s", startTime, endTime)
	}

	analysis, err := h.ollamaClient.AnalyzeLog(model, logContent, issueDescription)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("🔧 **设备日志查询结果** (AI分析失败)\n\n**环境**: %s (%s)\n**日志文件**: %s\n**设备ID**: %s\n**执行操作**: %s\n\n❌ **AI分析错误**: %s\n\n**建议**: 请手动分析日志内容，或检查Ollama服务是否正常运行",
						env.Name, envName, logFile.Path, deviceID, commandDesc, err.Error()),
				},
			},
		}, nil
	}

	// 构建时间范围描述
	timeRangeDesc := fmt.Sprintf("最近%s条", lines)
	if hasStartTime && hasEndTime {
		timeRangeDesc = fmt.Sprintf("%s 到 %s", startTime, endTime)
	}

	keywordDesc := ""
	if keyword != "" {
		keywordDesc = fmt.Sprintf("，关键词: %s", keyword)
	}

	// 返回完整的分析结果
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("🔍 **设备日志智能分析报告**\n\n**查询条件**:\n- 环境: %s (%s)\n- 日志文件: %s (%s)\n- 设备ID: %s\n- 时间范围: %s%s\n- AI模型: %s\n\n**执行操作**: %s\n\n---\n\n%s\n\n---\n\n",
					env.Name, envName, logFile.Path, logFile.Description, deviceID, timeRangeDesc, keywordDesc, model, commandDesc, analysis),
			},
		},
	}, nil
}

// readLocalLogFile 读取本地日志文件并过滤
func (h *LogCommandHandler) readLocalLogFile(logPath, deviceID, keyword string, lines int, startTime, endTime string) (string, error) {
	// 添加调试信息
	currentDir, _ := os.Getwd()
	fmt.Printf("DEBUG: 当前工作目录: %s\n", currentDir)
	fmt.Printf("DEBUG: 尝试访问文件: %s\n", logPath)

	// 检查文件是否存在
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return "", fmt.Errorf("日志文件不存在: %s (当前目录: %s)", logPath, currentDir)
	}

	file, err := os.Open(logPath)
	if err != nil {
		return "", fmt.Errorf("无法打开日志文件: %v", err)
	}
	defer file.Close()

	var matchedLines []string
	scanner := bufio.NewScanner(file)

	// 读取所有行并过滤
	for scanner.Scan() {
		line := scanner.Text()

		// 检查是否包含设备ID
		if !strings.Contains(line, deviceID) {
			continue
		}

		// 如果有关键词，检查是否包含关键词
		if keyword != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
			continue
		}

		// 如果有时间范围，检查时间
		if startTime != "" && endTime != "" {
			if !h.isTimeInRange(line, startTime, endTime) {
				continue
			}
		}

		matchedLines = append(matchedLines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取文件时出错: %v", err)
	}

	// 获取最后N行
	if len(matchedLines) > lines {
		matchedLines = matchedLines[len(matchedLines)-lines:]
	}

	return strings.Join(matchedLines, "\n"), nil
}

// isTimeInRange 检查日志行的时间是否在指定范围内
func (h *LogCommandHandler) isTimeInRange(logLine, startTime, endTime string) bool {
	// 简单的时间比较，假设日志格式为 "2025-07-24 15:43:47.123"
	if len(logLine) < 23 {
		return false
	}

	logTime := logLine[:23]
	return logTime >= startTime && logTime <= endTime
}

// 安全的设备ID验证
func validateDeviceID(deviceID string) error {
	// 只允许字母数字和特定字符
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, deviceID)
	if !matched || len(deviceID) > MAX_DEVICE_ID_LENGTH {
		return fmt.Errorf("无效的设备ID格式，只允许字母数字下划线和连字符，长度不超过%d", MAX_DEVICE_ID_LENGTH)
	}
	return nil
}

// 安全的关键词转义
func escapeShellArg(arg string) string {
	// 移除或转义危险字符
	arg = strings.ReplaceAll(arg, "'", "'\"'\"'")
	arg = strings.ReplaceAll(arg, "`", "\\`")
	arg = strings.ReplaceAll(arg, "$", "\\$")
	arg = strings.ReplaceAll(arg, "\\", "\\\\")
	arg = strings.ReplaceAll(arg, "\"", "\\\"")
	arg = strings.ReplaceAll(arg, ";", "\\;")
	arg = strings.ReplaceAll(arg, "&", "\\&")
	arg = strings.ReplaceAll(arg, "|", "\\|")
	arg = strings.ReplaceAll(arg, "<", "\\<")
	arg = strings.ReplaceAll(arg, ">", "\\>")
	arg = strings.ReplaceAll(arg, "(", "\\(")
	arg = strings.ReplaceAll(arg, ")", "\\)")
	arg = strings.ReplaceAll(arg, "!", "\\!")
	return arg
}
