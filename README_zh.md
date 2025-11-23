# Ollama 蜜罐

[English](README.md)

一个轻量级的 Go 应用程序，用于模拟 Ollama API 以检测和记录未经授权的访问 Ollama 服务的尝试。

## 问题陈述

Ollama 是一个流行的本地运行 LLM 的工具。随着其采用率的增长，它成为攻击者利用暴露的 API 端点目标。用户需要一种方式来检测他们的端口是否被扫描或被恶意行为者针对，这些行为者试图使用他们的计算资源或窃取数据。

## 解决方案

`ollama_honeypot` 是一个轻量级的 Go 应用程序，用于模拟 Ollama API。它监听标准的 Ollama 端口 (11434) 并使用假数据响应常见的 API 请求（如生成文本或列出模型）。至关重要的是，它记录每个请求，捕获攻击者的意图细节，如请求的模型和使用的提示。

## 功能

- **模拟 API 端点：**
  - `GET /`: 健康检查/状态。
  - `POST /api/generate`: 模拟文本生成，使用流式 NDJSON 响应。
  - `POST /api/chat`: 模拟聊天完成。
  - `GET /api/tags`: 列出假的可用模型。
  - `POST /api/pull`: 模拟模型拉取，具有可配置的下载速度。
  - `DELETE /api/delete`: 模拟模型删除。
  - `GET /api/ps`: 模拟运行模型状态列表。
  - `GET /api/show`: 显示模型信息。
  - `GET /api/version`: 获取版本信息。
  - `GET /v1/models`: 模拟 OpenAI 的模型列表 API。
  - `POST /v1/chat/completions`: 模拟 OpenAI 的聊天完成 API，支持流式传输。
  - `POST /chat/completions`: 模拟 OpenAI 的聊天完成 API，支持流式传输。

- **请求日志记录：** 使用结构化 JSON 日志记录捕获 IP、时间戳、端点、方法和请求体（提示）。
- **配置：** 通过环境变量或 CLI 标志支持自定义端口、日志偏好和模拟数据路径。
- **低资源使用：** 设计为轻量级，可与其他服务并行运行。

## 安装

1. 确保安装了 Go 1.25.4 或更高版本。

2. 克隆仓库：
   ```bash
   git clone https://github.com/pedxyuyuko/ollama_honeypot.git
   cd ollama_honeypot
   ```

3. （可选）预填充模拟数据文件：
   ```bash
   mkdir -p mock
   curl -o mock/tags.json https://raw.githubusercontent.com/pedxyuyuko/ollama_honeypot/refs/heads/master/mock/tags.json
   curl -o mock/response.json https://raw.githubusercontent.com/pedxyuyuko/ollama_honeypot/refs/heads/master/mock/response.json
   curl -o mock/version.json https://raw.githubusercontent.com/pedxyuyuko/ollama_honeypot/refs/heads/master/mock/version.json
   ```

4. 构建应用程序：
   ```bash
   go build -o ollama_honeypot .
   ```

## Docker

您可以使用 Docker 运行蜜罐，以实现轻松部署和隔离。

### 使用 Docker Compose

1. 确保安装了 Docker 和 Docker Compose。

2. 复制 `.example.env` 到 `.env` 并根据需要配置环境变量。

3. 运行蜜罐：

   ```bash
   docker-compose up
   ```

   这将在端口 11434 上启动蜜罐，挂载 `mock` 和 `logs` 目录以实现持久数据存储。

### 手动构建和运行

1. 构建 Docker 镜像：

   ```bash
   docker build -t ollama_honeypot .
   ```

2. 运行容器：

   ```bash
   docker run -p 11434:11434 -v $(pwd)/mock:/app/mock -v $(pwd)/logs:/app/logs --env-file .env ollama_honeypot
   ```

   - `-p 11434:11434`：将容器的端口 11434 映射到主机的端口 11434。
   - `-v $(pwd)/mock:/app/mock`：挂载本地 `mock` 目录以持久化模拟数据。
   - `-v $(pwd)/logs:/app/logs`：挂载本地 `logs` 目录以持久化日志。
   - `--env-file .env`：从 `.env` 文件加载环境变量。

## 使用

运行蜜罐服务器：

```bash
./ollama_honeypot serve
```

默认情况下，它在端口 11434（Ollama 的标准端口）上启动。

### CLI 选项

- `--port`: 指定要绑定的端口（默认：11434）
- `--log-path`: 审计日志文件的路径（可选，如果未设置，则仅记录到控制台）
- `--mock-path`: 包含模拟数据文件的目录路径（默认：./mock）
- `--help`: 显示帮助信息

## 配置

可以通过环境变量或 CLI 标志进行配置。环境变量可以在 `.env` 文件中设置。

复制 `.example.env` 到 `.env` 并根据需要修改值：

```bash
cp .example.env .env
```

### 环境变量

- `PORT`: 要绑定服务器的端口（默认：11434）
- `MOCK_PATH`: 包含模拟数据文件的目录路径（默认：./mock）
- `LOG_PATH`: 审计日志文件的路径（可选，如果未设置，则仅记录到控制台）
- `DEBUG`: 启用调试模式（设置为 1 以启用调试日志，默认：0）
- `DOWNLOAD_SPEED`: 模拟模型拉取的下载速度（字节/秒，默认：52428800，即 50MB/s）
- `DOWNLOAD_SPEED_VARIANCE`: 随机速度波动的方差因子（0.0 到 1.0，默认：0.2）
- `DOWNLOAD_SPEED_WAVE_PERIOD`: 正弦速度变化的周期（秒，默认：1.0）
- `DOWNLOAD_SPEED_WAVE_AMPLITUDE`: 正弦速度变化的振幅因子（0.0 到 1.0，默认：0.5）

## 日志记录

应用程序使用 Logrus 进行结构化日志记录。请求细节（IP、时间戳、方法、路径、主体）被记录以供分析。

- 控制台输出：文本格式，便于阅读。
- 文件输出（如果设置了 LOG_PATH）：JSON 格式，便于解析。

## 依赖

- `github.com/gin-gonic/gin`: HTTP Web 框架。
- `github.com/spf13/cobra`: CLI 应用程序结构。
- `github.com/sirupsen/logrus`: 结构化日志记录库。
- `github.com/joho/godotenv`: 环境变量加载。

## 贡献

欢迎贡献！请在 GitHub 上打开问题或提交拉取请求。

## 许可证

此项目根据 LICENSE 文件中指定的条款获得许可。