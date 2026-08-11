package constant

// ------------------login相关------------------
type device struct {
	Web     string // 网页端
	Desktop string // 桌面端（Electron、Tauri）
	Mobile  string // 移动端（App）
	Unknown string // 未知来源
}

// Device 登录设备
var Device = device{
	Web:     "web",
	Desktop: "desktop",
	Mobile:  "mobile",
	Unknown: "unknown",
}

// ------------------app platform相关------------------
type appPlatform struct {
	Android string
	IOS     string
	Windows string
	MacOS   string
	Web     string
}

// AppPlatform 客户端平台
var AppPlatform = appPlatform{
	Android: "android",
	IOS:     "ios",
	Windows: "windows",
	MacOS:   "macos",
	Web:     "web",
}

func (c appPlatform) Validate(v string) bool {
	switch v {
	case c.Android, c.IOS, c.Windows, c.MacOS, c.Web:
		return true
	default:
		return false
	}
}
