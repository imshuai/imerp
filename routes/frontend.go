package routes

import (
	"erp/embedded"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupFrontendRoutes 配置前端路由
func SetupFrontendRoutes(r *gin.Engine) {
	// 获取嵌入的文件系统
	distFS, err := fs.Sub(embedded.DistFS(), "dist")
	if err != nil {
		// 如果没有构建产物，返回提示信息
		r.GET("/*filepath", func(c *gin.Context) {
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
		c.Request.URL.Path = c.Request.URL.Path[1:] // 去掉前导 /
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	// SPA fallback - 所有其他路由返回 index.html
	r.GET("/*filepath", func(c *gin.Context) {
		c.Request.URL.Path = "/index.html"
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
