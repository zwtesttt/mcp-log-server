# MCP Log Server - 智能日志分析服务 v2.3.0

基于 [mcp-go](https://github.com/mark3labs/mcp-go) SDK 构建的智能日志分析服务，支持多环境日志查询、时间范围过滤和AI智能分析。

## 🚀 核心功能

- 🌍 **多环境支持**: 支持 dev、test、staging、prod 四个环境
- 📁 **本地日志管理**: 本地文件系统存储，支持动态路径配置
- 🔍 **设备级查询**: 按设备ID精准查询日志记录
- ⏰ **时间范围过滤**: 支持精确的时间范围查询
- 🤖 **AI智能分析**: 集成Ollama，提供专业的日志分析报告
- 📊 **结构化输出**: 清晰的markdown格式分析报告

## 📁 项目架构

```
mcp-log-server/
├── main.go                    # 主程序入口
├── internal/
│   ├── config/
│   │   └── config.go         # 多环境配置管理
│   ├── handlers/
│   │   └── log_commands.go   # 日志命令处理器
│   └── ollama/
│       └── client.go         # Ollama AI客户端
├── logs/                      # 日志文件存储
│   ├── dev/                   # 开发环境日志
│   │   ├── oms.log
│   │   └── blackhole.log
│   ├── test/                  # 测试环境日志
│   │   ├── oms.log
│   │   └── blackhole.log
│   ├── staging/               # 预发布环境日志
│   │   ├── oms.log
│   │   └── blackhole.log
│   └── prod/                  # 生产环境日志
│       ├── oms.log
│       └── blackhole.log
├── go.mod
├── go.sum
└── README.md
```

## 🛠️ 快速开始

### 1. 环境要求

- Go 1.21+
- Ollama (用于AI分析功能)

### 2. 安装 Ollama

```bash
# macOS
brew install ollama

# 启动服务
ollama serve

# 下载推荐模型
ollama pull gemma3:27b
```

### 3. 构建项目

```bash
git clone <your-repo-url>
cd mcp-log-server
go mod tidy
go build -o mcp-server .
```

### 4. 配置 MCP 客户端

#### Cursor 配置 (`~/.cursor/mcp.json`)

```json
{
  "mcpServers": {
    "mcp-log-server": {
      "command": "/path/to/mcp-log-server/mcp-server"
    }
  }
}
```

#### Claude Desktop 配置

```json
{
  "mcpServers": {
    "mcp-log-server": {
      "command": "/path/to/mcp-log-server/mcp-server"
    }
  }
}
```

## 🔧 核心工具

### `query_device_logs_by_time` - 主要功能

按设备ID和时间范围查询日志，并提供AI智能分析。

**参数**:
- `environment` (必需): 环境名称 (`dev`/`test`/`staging`/`prod`)
- `log_type` (必需): 日志类型 (`oms`/`blackhole`)
- `device_id` (必需): 设备ID
- `start_time` (可选): 开始时间 (格式: `2025-07-24 11:59:38.369`)
- `end_time` (可选): 结束时间 (格式: `2025-07-24 11:59:38.369`)
- `keyword` (可选): 搜索关键词
- `lines` (可选): 查询行数 (默认: 2000)
- `model` (可选): AI分析模型

## 💡 使用示例

### 基础设备查询

```
用户: "查看设备1234在开发环境的日志"
```

AI会调用工具查询 `dev` 环境中设备 `1234` 的最新日志，并提供智能分析。

### 时间范围查询

```
用户: "查看设备1234在7月25日的OMS日志"
```

AI会查询指定日期范围内的日志记录。

### 跨环境对比

```
用户: "对比设备1234在测试环境和生产环境的表现"
```

AI会分别查询两个环境的日志并进行对比分析。

## 📊 AI分析报告格式

每次查询都会生成结构化的分析报告：

```markdown
## 🔍 异常识别与分类
- 错误日志统计和分类
- 警告信息汇总
- 异常模式识别

## 📈 系统状态评估  
- 设备运行状态
- 性能指标分析
- 通信状态评估

## ⚠️ 风险评估
- 紧急程度分级
- 影响范围评估
- 具体建议措施

## 📋 日志统计摘要
- 时间范围
- 日志级别统计
- 涉及设备清单
```

## 🌍 多环境配置

### 支持的环境

| 环境 | 说明 | 主机配置 |
|------|------|----------|
| `dev` | 开发环境 | localhost |
| `test` | 测试环境 | test.example.com |
| `staging` | 预发布环境 | staging.example.com |
| `prod` | 生产环境 | prod.example.com |

### 日志类型

| 类型 | 描述 | 文件路径 |
|------|------|----------|
| `oms` | OMS系统日志 | `/logs/{env}/oms.log` |
| `blackhole` | 黑洞服务日志 | `/logs/{env}/blackhole.log` |

## ⚙️ 配置自定义

修改 `internal/config/config.go` 添加新环境或日志类型：

```go
// 添加新环境
"staging": {
    Name: "预发布环境",
    Host: "staging.example.com",
    User: "staginguser",
    Port: "22",
},

// 添加新日志类型
"api": {
    Name:        "api",
    Path:        "", // 动态路径
    Description: "API服务日志",
    Aliases:     []string{"接口"},
},
```

## 📈 版本历史

### v2.3.0 - 多环境支持 (当前版本)
- ✨ 支持四个独立环境 (dev/test/staging/prod)
- ✨ 动态日志文件路径配置
- ✨ 跨日期的时间范围查询
- ✨ 改进的AI分析提示词
- 🔧 本地文件系统存储

### v2.2.0 - 设备查询功能
- ✨ 按设备ID查询日志
- ✨ 时间范围过滤支持
- ✨ 结构化AI分析报告

### v2.1.0 - AI集成
- ✨ Ollama AI分析集成
- ✨ 智能日志解读
- ✨ 问题诊断建议

### v2.0.0 - 架构重构
- 🔧 模块化代码结构
- 🔧 配置管理优化
- 🔧 MCP协议支持

## 🔍 故障排除

### 常见问题

1. **日志文件不存在**
   - 检查日志文件路径配置
   - 确认文件读取权限

2. **AI分析失败**
   - 检查Ollama服务状态
   - 确认模型是否正确下载

3. **时间范围查询无结果**
   - 验证时间格式 (`2025-07-24 11:59:38.369`)
   - 检查日志文件中是否有对应时间段的数据

### 调试模式

```bash
# 启动时查看详细日志
./mcp-server --debug
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [mcp-go](https://github.com/mark3labs/mcp-go) - MCP协议Go实现
- [Ollama](https://ollama.ai) - 本地AI模型运行环境

---

**🚀 立即开始使用智能日志分析，让AI帮助您快速定位和解决系统问题！**
