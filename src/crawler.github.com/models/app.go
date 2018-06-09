package models

import "time"

// APP 应用信息
type APP struct {
	ID       uint   `json:"id,omitempty"`       // ID
	Name     string `json:"name,omitempty"`     // 应用名称
	PkgName  string `json:"pkg_name,omitempty"` // 包名
	APPID    string `json:"appid,omitempty"`    // appid
	Category string `json:"category,omitempty"` // 分类
	Size     string `json:"size,omitempty"`     // 软件大小
	Version  string `json:"version,omitempty"`  // 版本号
	Company  string `json:"company,omitempty"`  // 开发商

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
