# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AuditLimit is a rate limiting service for AI model API proxies. It enforces per-user, per-model rate limits across ChatGPT, Claude, Grok, and Gemini endpoints, with optional content moderation via OpenAI's moderation API and local keyword filtering.

## Build & Run

```bash
# Run locally (requires Go 1.18+)
go run main.go

# Run tests
go test ./api/...

# Build for linux/amd64 (requires GoFrame CLI `gf`)
gf build main.go -a amd64 -s linux -p ./temp

# Build and push Docker image
bash release.sh

# Deploy with Docker
docker-compose up -d
```

## Architecture

**Framework**: GoFrame v2 (`github.com/gogf/gf/v2`) — provides HTTP server, config, logging, JSON utilities.

**Entry point**: `main.go` registers four POST endpoints and loads keywords from `data/keywords.txt` on startup.

### Packages

- **`config/`** — Global config vars (`PORT`, `OAIKEY`, `MODERATION`, `ForbiddenWords`). Reads from env vars and `config/config.yaml` via GoFrame's config system.
- **`api/`** — All request handling and rate limiting logic:
  - `audit_limit.go` — Four handler functions (`GPTAuditLimit`, `ClaudeAuditLimit`, `GrokAuditLimit`, `GeminiAuditLimit`), each tailored to the provider's request/response format.
  - `limit.go` — Token-bucket rate limiter using `golang.org/x/time/rate`. Limiters are stored in an in-memory map keyed by `token|MODEL`. Weekly cleanup of stale entries.
  - `result.go` — Shared response message constants.

### Request Flow (all four handlers follow this pattern)

1. Extract Bearer token from `Authorization` header (user identity)
2. Parse model name from request JSON (field name varies by provider)
3. Check prompt against forbidden keywords → 400
4. Call OpenAI moderation API if `OAIKEY` is set → 400
5. Look up or create rate limiter for `{PREFIX}-{model}` (e.g., `CHATGPT-GPT-4O`)
6. If limit=0 in config → 403 (model disabled)
7. If tokens exhausted → 429 with retry delay
8. Otherwise → 200, consume one token

### Rate Limit Configuration

Format: `"{LIMIT}/{DURATION}"` (e.g., `"20/3h"`, `"7/24h"`, `"0/1h"` to disable). Default: `"40/3h"`.

Config keys use prefix convention: `CHATGPT-`, `CLAUDE-`, `GROK-`, `GEMINI-` + uppercase model name. Set via environment variables or `config/config.yaml`.

### Provider-Specific Differences

| Provider | Endpoint | Model field | Prompt field | Model normalization |
|----------|----------|------------|-------------|-------------------|
| ChatGPT | `/chatgpt/audit_limit` | `model` | `messages.0.content.parts.0` | `system_hints` can override to `research`/`agent` |
| Claude | `/claude/audit_limit` | `model` | `prompt` | Normalized to `sonnet`/`opus`/`haiku` |
| Grok | `/grok/audit_limit` | `modelName` | `message` | Used as-is |
| Gemini | `/gemini/audit_limit` | `model` (header) | N/A (moderation disabled) | Used as-is |

### Error Response Formats

Each provider uses a different JSON error structure to match what the upstream client expects. ChatGPT uses `{"error": "..."}`, Claude uses `{"type":"error","error":{"type":"...","message":"..."}}`, Grok uses `{"error":{"code":13,"message":"...","detail":[...]}}`.
