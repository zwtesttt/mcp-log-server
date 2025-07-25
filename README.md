# 日志命令生成器 MCP 服务器 v2.0

这是一个基于 [mcp-go](https://github.com/mark3labs/mcp-go) SDK 构建的智能日志命令生成器，支持通过自然语言生成日志查看命令，并集成Ollama AI模型进行日志分析。

## 🚀 功能特性

- 🔧 **智能命令生成**: 根据环境、日志类型、关键词等参数生成SSH日志查看命令
- 🌍 **多环境支持**: 支持dev、test、staging、prod等环境
- 📝 **多日志类型**: 支持blackhole、oms、api、error等日志文件
- 🔍 **关键词过滤**: 支持根据关键词过滤日志内容
- 🤖 **AI智能分析**: 集成Ollama模型，支持日志内容智能分析
- 💬 **AI问答**: 支持通过AI回答运维相关问题
- 📊 **模块化设计**: 清晰的包结构，易于维护和扩展

## 📁 项目结构

```
mcp-log-server/
├── main.go                    # 主程序入口
├── internal/
│   ├── config/
│   │   └── config.go         # 配置管理
│   ├── handlers/
│   │   ├── log_commands.go   # 日志命令处理器
│   │   └── ollama_handlers.go # AI处理器
│   └── ollama/
│       └── client.go         # Ollama客户端
├── go.mod
├── go.sum
└── README.md
```

## 🛠️ 安装和运行

### 1. 环境要求

- Go 1.21 或更高版本
- Ollama (可选，用于AI功能)

### 2. 安装Ollama (可选)

如果需要使用AI分析功能，请先安装Ollama：

```bash
# macOS
brew install ollama

# 或从官网下载: https://ollama.ai

# 启动Ollama服务
ollama serve

# 下载模型
ollama pull qwen2.5:0.5b
```

### 3. 构建服务器

```bash
go mod tidy
go build -o mcp-server main.go
```

### 4. 配置MCP客户端

#### 在 Cursor 中使用

编辑 `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "log-command-generator": {
      "command": "/path/to/your/mcp-server"
    }
  }
}
```

#### 在 Claude Desktop 中使用

编辑配置文件:

```json
{
  "mcpServers": {
    "log-command-generator": {
      "command": "/path/to/your/mcp-server"
    }
  }
}
```

## 🔧 可用工具

### 1. generate_log_command
生成日志查看命令

**参数**:
- `environment` (必需): 环境名称 (dev/test/staging/prod)
- `log_type` (必需): 日志类型 (blackhole/黑洞/oms/api/error/错误)
- `keyword` (可选): 搜索关键词
- `lines` (可选): 显示行数，默认1000

### 2. generate_and_analyze_log
生成日志查看命令，自动执行并使用AI分析结果

**参数**:
- `environment` (必需): 环境名称 (dev/test/staging/prod)
- `log_type` (必需): 日志类型 (blackhole/黑洞/oms/api/error/错误)
- `keyword` (可选): 搜索关键词
- `lines` (可选): 显示行数，默认100行
- `issue_description` (可选): 问题描述，帮助AI更好地分析
- `model` (可选): AI模型名称

### 3. query_device_logs_by_time ⭐ 新功能
按设备ID和时间范围查询日志并进行AI分析

**参数**:
- `environment` (必需): 环境名称 (dev/test/staging/prod)
- `log_type` (必需): 日志类型 (blackhole/黑洞/oms/api/error/错误)
- `device_id` (必需): 设备ID
- `days` (可选): 查询天数，默认2天
- `issue_description` (可选): 问题描述，帮助AI更好地分析
- `model` (可选): AI模型名称

### 4. ask_ollama
调用Ollama模型进行问答

**参数**:
- `prompt` (必需): 问题或内容
- `model` (可选): 模型名称，默认qwen2.5:0.5b
- `context` (可选): 上下文信息

### 5. analyze_log_with_ai
使用AI分析日志内容

**参数**:
- `log_content` (必需): 要分析的日志内容
- `issue_description` (可选): 问题描述
- `model` (可选): 模型名称

### 6. list_environments
列出所有可用环境

### 7. list_log_types  
列出所有支持的日志类型

## 💡 使用示例

### 基本日志查看
**你**: "帮我生成一个查看dev环境blackhole日志的命令"

**AI**: 会生成类似这样的命令：
```bash
ssh deploy@dev-server.company.com 'tail -1000 /var/log/app/blackhole.log'
```

### 带关键词过滤
**你**: "看test环境的error日志，搜索包含'timeout'的记录"

**AI**: 会生成：
```bash
ssh deploy@test-server.company.com 'tail -1000 /var/log/app/error.log | grep -i "timeout"'
```

### 设备日志查询 ⭐ 新功能
**你**: "帮我输出oms开发环境 1234设备最近两天的日志"

**AI**: 会调用 `query_device_logs_by_time` 工具:
- 自动生成查询命令: `ssh develop@mm01.sca.im -p 59822 'find $(dirname /data/develop/oms/logs/oms.log) -name "*.log*" -mtime -2 -exec grep -l "1234" {} \; | head -10 | xargs grep "1234" | tail -500'`
- 自动执行获取日志内容
- AI智能分析后提供:
  - 🔍 设备状态分析: 设备1234在最近两天的运行情况
  - 🎯 问题识别: 发现的错误或异常情况
  - 💡 解决建议: 针对性的修复建议
  - 📈 趋势分析: 设备运行趋势和性能指标

**使用参数**:
- 环境: `dev` (开发环境)
- 日志类型: `oms` (OMS系统日志)
- 设备ID: `1234`
- 时间范围: 最近2天

### AI日志分析
**你**: "分析这段日志内容：[粘贴日志内容]"

**AI**: 会调用Ollama模型分析日志，提供：
- 错误识别
- 问题原因分析
- 解决建议
- 系统状态总结

### AI问答
**你**: "什么情况下会出现数据库连接超时？"

**AI**: 会基于运维知识回答相关问题

## ⚙️ 配置自定义

可以修改 `internal/config/config.go` 中的配置：

```go
// 修改环境配置
"dev": {
    Name: "开发环境",
    Host: "your-dev-server.com",
    User: "your-user",
    Port: "22",
},

// 修改日志文件配置
"blackhole": {
    Name:        "blackhole",
    Path:        "/your/log/path/blackhole.log",
    Description: "黑洞服务日志",
    Aliases:     []string{"黑洞"},
},

// 修改Ollama配置
Ollama: OllamaConfig{
    BaseURL:      "http://localhost:11434",
    DefaultModel: "your-preferred-model",
    Timeout:      30,
},
```

## 🔄 典型工作流

1. **问题报告**: "客户说数据有问题"
2. **生成命令**: "帮我看下prod环境的blackhole日志"
3. **执行命令**: 复制生成的SSH命令到终端执行
4. **AI分析**: "分析这段日志内容，重点关注数据同步问题"
5. **获得建议**: AI提供问题诊断和解决方案

## 🚀 版本历史

- **v2.2.0**: 新增设备日志查询功能，支持按设备ID和时间范围查询日志
- **v2.1.0**: 新增智能日志分析工作流，支持一键生成命令→自动执行→AI智能分析
- **v2.0.0**: 重构代码结构，支持模块化设计，集成Ollama AI功能
- **v1.0.0**: 基础日志命令生成功能

## 📄 许可证

MIT License

## 🆕 v2.1.0 新功能

### 智能日志分析工作流

现在支持**一键生成命令→自动执行→AI智能分析**的完整工作流！

#### 新增工具: `generate_and_analyze_log`

这个强大的新工具会：
1. 🔧 根据你的参数生成SSH日志查看命令
2. ⚡ 自动执行命令获取日志内容  
3. 🤖 使用Ollama AI模型智能分析日志
4. 📊 提供完整的分析报告和建议

**参数**:
- `environment` (必需): 环境名称
- `log_type` (必需): 日志类型  
- `keyword` (可选): 搜索关键词
- `lines` (可选): 显示行数，默认100行
- `issue_description` (可选): 问题描述，帮助AI更精准分析
- `model` (可选): 指定AI模型

### 使用示例

#### 场景1: 快速问题诊断
**你**: "客户反馈登录失败，帮我分析dev环境的api日志，关键词是login"

**AI**: 会调用 `generate_and_analyze_log` 工具:
- 生成命令: `ssh deploy@dev-server.company.com 'tail -100 /var/log/app/api.log | grep -i "login"'`
- 自动执行获取日志
- AI分析后提供:
  - 🔍 错误识别: 发现认证失败
  - 🎯 问题原因: Token过期或用户权限问题
  - 💡 解决建议: 检查认证服务状态，更新用户权限配置
  - 📈 系统状态: 整体稳定，仅个别用户受影响

#### 场景2: 性能问题分析
**你**: "prod环境响应很慢，帮我看下oms日志，重点关注性能问题"

**AI**: 会:
- 执行: `ssh deploy@prod-server.company.com 'tail -100 /var/log/app/oms.log'`
- 分析性能相关指标
- 提供优化建议

### 工具对比

| 工具名称 | 功能 | 适用场景 |
|---------|------|----------|
| `generate_log_command` | 仅生成命令 | 需要手动执行命令时 |
| `generate_and_analyze_log` | 生成+执行+AI分析 | 快速问题诊断和分析 |
| `analyze_log_with_ai` | 分析已有日志内容 | 已获取日志，需要AI分析 |

### 安全说明

⚠️ **重要**: `generate_and_analyze_log` 工具会实际执行SSH命令，请确保:
- SSH密钥配置正确
- 服务器连接安全可靠  
- 具有相应的日志文件读取权限
- 在受信任的网络环境中使用

如果不希望自动执行命令，请使用 `generate_log_command` 工具仅生成命令。
