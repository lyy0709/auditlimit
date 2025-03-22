package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	PORT = 8080
	ForbiddenWords = []string{} // 禁止词
	OAIKEY    = "" // OAIKEY
	OAIKEYLOG = "" // OAIKEYLOG 隐藏
	MODERATION = "https://gateway.ai.cloudflare.com/v1/040ac2002b4dd67637e97c628feb3484/xyhelper/openai/moderations"
)

func init() {
	ctx := gctx.GetInitCtx()
	port := g.Cfg().MustGetWithEnv(ctx, "PORT").Int()
	if port > 0 {
		PORT = port
	}
	g.Log().Info(ctx, "PORT:", PORT)
	oaikey := g.Cfg().MustGetWithEnv(ctx, "OAIKEY").String()
	// oaikey 不为空
	if oaikey != "" {
		OAIKEY = oaikey
		// 日志隐藏 oaikey，有 * 代表有值
		OAIKEYLOG = "******"
	}
	g.Log().Info(ctx, "OAIKEY:", OAIKEYLOG)
	moderation := g.Cfg().MustGetWithEnv(ctx, "MODERATION").String()
	if moderation != "" {
		MODERATION = moderation
	}
	g.Log().Info(ctx, "MODERATION:", MODERATION)
}
