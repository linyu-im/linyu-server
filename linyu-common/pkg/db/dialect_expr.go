package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DateSubDaysExpr 生成「基准时间减去 daysColumn 天」表达式。
// mysql: DATE_SUB(?, INTERVAL col DAY)
// sqlite: datetime(?, '-' || col || ' days')
// 调用方 Where 中使用 "?placeholder" 作为 time 参数占位，daysColumn 为列名如 ms.expire_days。
func DateSubDaysExpr(gdb *gorm.DB, daysColumn string) string {
	if isSQLite(gdb) {
		return "datetime(?, '-' || " + daysColumn + " || ' days')"
	}
	return "DATE_SUB(?, INTERVAL " + daysColumn + " DAY)"
}

// PrefixPathExpr 用 newPath 替换 path 列前缀（从 oldPathLen+1 起保留后缀）。
// mysql: CONCAT(?, SUBSTRING(path, ?))
// sqlite: ? || substr(path, ?)
func PrefixPathExpr(gdb *gorm.DB, newPath string, oldPathLen int) clause.Expr {
	start := oldPathLen + 1
	if isSQLite(gdb) {
		return gorm.Expr("? || substr(path, ?)", newPath, start)
	}
	return gorm.Expr("CONCAT(?, SUBSTRING(path, ?))", newPath, start)
}

func isSQLite(gdb *gorm.DB) bool {
	if gdb == nil || gdb.Dialector == nil {
		return false
	}
	name := gdb.Dialector.Name()
	return name == "sqlite" || name == "sqlite3"
}
