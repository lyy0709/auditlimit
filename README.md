# AuditLimit

用于内容审核的限流器,本代码更多的是为了演示如果使用限流器,大家可以根据自己的需求进行修改。

如果不想部署,可以使用公共限流器: `https://auditlimit.closeai.biz/audit_limit`

## 部署方法

创建`docker-compose.yml`文件

```yml
version: '3'
services:
  auditlimit:
    image: lyy0709/auditlimit
    restart: always
    ports:
      - "127.0.0.1:9611:8080"
    volumes:
      - ./data:/app/data
    environment:
      PORT: 9611
      OAIKEY: "" # OpenAI API key 用于内容审核
      # ChatGPT 模型限速 (前缀 CHATGPT-, model字段直接拼接)
      CHATGPT-AUTO: "200/3h"
      CHATGPT-TEXT-DAVINCI-002-RENDER-SHA: "200/3h"
      CHATGPT-GPT-4O-MINI: "200/3h"
      CHATGPT-GPT-4O: "60/3h"
      CHATGPT-GPT-4: "20/3h"
      CHATGPT-GPT-4O-CANMORE: "30/3h"
      CHATGPT-O1-PREVIEW: "7/24h"
      CHATGPT-O1-MINI: "50/24h"
      CHATGPT-GPT-5: "60/3h"
      CHATGPT-RESEARCH: "2/24h"
      CHATGPT-AGENT: "10/24h"
      # Claude 模型限速 (代码将所有模型归一化为 sonnet/opus/haiku 三类)
      CLAUDE-SONNET: "20/3h"
      CLAUDE-OPUS: "20/3h"
      CLAUDE-HAIKU: "20/3h"
      # Grok 模型限速 (前缀 GROK-, modelName字段直接拼接)
      GROK-GROK-3: "20/3h"
      GROK-GROK-4: "20/3h"
      GROK-GROK-4-AUTO: "20/3h"
      GROK-REASONING: "20/24h"
      GROK-DEEPSEARCH: "20/24h"
      # Gemini 模型限速 (前缀 GEMINI-, model header直接拼接)
      GEMINI-PRO: "20/3h"
      GEMINI-FAST: "20/3h"
      GEMINI-THINKING: "20/3h"
      
```

然后执行

```bash
docker-compose up -d
```

限流器接口地址为: `http://ip:9611/audit_limit`

## 超速返回格式

状态码: 429

```json
{
  "detail": {
    "clears_in": 252,
    "code": "model_cap_exceeded",
    "message": "You have sent too many messages to the model. Please try again later."
  }
}
````

## 通用提示

状态码: 400

```json
{
  "detail": "别闹了"
}
```

## 正常返回

状态码: 200


## 内容审核

配置环境变量

OAIKEY: "sk-xxxxxx"  # api.openai.com可用的key或sess

将启用moderations接口内容审核