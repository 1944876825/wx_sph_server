# 微信视频号 打招呼自动回复 后端

基于 Go + Gin + SQLite 的微信视频号自动回复服务，支持固定文本回复和 **AI 智能回复**。

## ✨ 功能特性

- 🔔 自动监听视频号打招呼消息
- 🤖 **AI 智能回复** — 接入 GLM 大模型，根据用户消息自动生成个性化回复
- 📝 固定文本/图片回复 — 支持多条消息随机回复
- ⚙️ 可配置系统提示词，自定义 AI 角色和回复风格
- 👥 多账号管理
- 🖥️ Web 管理后台（前端嵌入后端二进制，单文件部署）

## 技术栈

- Go 1.22 + Gin + GORM + SQLite
- GLM API (OpenAI 兼容格式)
- 前端 embed 进 Go 二进制

## 快速开始

### 配置

编辑 `config.yaml`：

```yaml
port: 5000
title: 视频号助手
icon: ./assets/vite.svg

# AI 智能回复配置
ai:
  enabled: true
  api_key: "your-api-key"          # 智谱 AI API Key
  base_url: "https://open.bigmodel.cn/api/paas/v4"
  model: "glm-5-turbo"              # 可替换为 glm-5.2 等模型
  default_prompt: "你是一个友好的微信视频号客服助手。请用简洁、亲切的语气回复用户的打招呼消息。回复不超过50个字。"
```

### 运行

```bash
go build -o wx_sph_server .
./wx_sph_server
```

访问 `http://127.0.0.1:5000` 进入管理后台。

## AI 回复模式

在管理后台的消息设置中，可以为每条回复消息开启 **AI 智能回复**：

- **开启后**：系统会根据用户发送的消息内容，调用 GLM 大模型自动生成回复
- **系统提示词**：可自定义 AI 的角色和回复风格，留空则使用全局默认提示词
- **降级机制**：AI 调用失败时自动降级为固定文本回复

## 前端地址

https://github.com/1944876825/wx_sph_web

## LICENSE

[MIT](https://opensource.org/license/mit/)
