package constant

type device struct {
	Web     string // 网页端
	Desktop string // 桌面端（Electron、Tauri）
	Mobile  string // 移动端（App）
	Unknown string // 未知来源
}

var Device = device{
	Web:     "web",
	Desktop: "desktop",
	Mobile:  "mobile",
	Unknown: "unknown",
}
