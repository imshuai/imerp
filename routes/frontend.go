package routes

import (
	"erp/embedded"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupFrontendRoutes 配置前端路由
func SetupFrontendRoutes(r *gin.Engine) {
	// 获取嵌入的文件系统
	distFS, err := fs.Sub(embedded.DistFS(), "dist")
	if err != nil {
		// 如果没有构建产物，返回提示信息
		r.NoRoute(func(c *gin.Context) {
			c.HTML(200, "", `
			<!DOCTYPE html>
			<html>
			<head>
				<title>前端未构建</title>
				<style>
					body { font-family: Arial, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #f5f5f5; }
					.box { text-align: center; padding: 40px; background: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
					h1 { color: #409EFF; }
					p { color: #606266; }
					code { background: #f5f5f5; padding: 2px 6px; border-radius: 4px; font-family: monospace; }
				</style>
			</head>
			<body>
				<div class="box">
					<h1>前端资源未构建</h1>
					<p>请先构建前端资源：</p>
					<p><code>cd frontend && npm run build</code></p>
				</div>
			</body>
			</html>
			`)
		})
		return
	}

	// 创建HTTP文件系统服务器
	fileServer := http.FileServer(http.FS(distFS))

	// 静态资源服务（处理 /assets 路径）
	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/")
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	// SPA fallback - 使用 NoRoute 处理所有其他路由
	r.NoRoute(func(c *gin.Context) {
		// 跳过API路由
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(404, gin.H{"code": 1, "message": "Not Found"})
			return
		}
		// 对于所有其他前端路由，读取并返回 index.html
		indexData, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			c.String(404, "index.html not found: %v", err)
			return
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Data(200, "text/html; charset=utf-8", indexData)
	})
}
