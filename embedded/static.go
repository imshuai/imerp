package embedded

import "embed"

// StaticFS 嵌入的静态文件系统
//
//go:embed dist
var StaticFS embed.FS

// DistFS 返回嵌入的文件系统
func DistFS() *embed.FS {
	return &StaticFS
}
