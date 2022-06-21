package models

import (
	"time"
)

// Users [...]
type Users struct {
	UID         int64     `gorm:"primaryKey;column:uid;type:bigint;not null" json:"-"`                                        // 用户唯一id
	Uname       string    `gorm:"unique;column:uname;type:varchar(64);not null" json:"uname"`                                 // 用户名
	Passwd      string    `gorm:"column:passwd;type:varchar(64);not null;default:''" json:"passwd"`                           // 密码
	Nickname    string    `gorm:"column:nickname;type:varchar(64);not null;default:''" json:"nickname"`                       // 昵称
	Roleid      int       `gorm:"column:roleid;type:int;not null;default:1" json:"roleid"`                                    // 角色码
	Gender      int8      `gorm:"column:gender;type:tinyint;not null;default:0" json:"gender"`                                // 性别
	Avatar      string    `gorm:"column:avatar;type:varchar(1024);not null;default:''" json:"avatar"`                         // 头像
	Phone       string    `gorm:"column:phone;type:varchar(32);not null;default:''" json:"phone"`                             // 电话号码
	Email       string    `gorm:"column:email;type:varchar(64);not null;default:''" json:"email"`                             // 电子邮箱
	Stat        int8      `gorm:"column:stat;type:tinyint;not null;default:0" json:"stat"`                                    // 状态码
	CreateAt    time.Time `gorm:"column:create_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"createAt"`          // 创建时间
	LastLoginAt time.Time `gorm:"column:last_login_at;type:datetime;not null;default:1000-01-01 00:00:00" json:"lastLoginAt"` // 最后登录时间
}

// TableName get sql table name.获取数据库表名
func (m *Users) TableName() string {
	return "users"
}

// UsersColumns get sql column name.获取数据库列名
var UsersColumns = struct {
	UID         string
	Uname       string
	Passwd      string
	Nickname    string
	Roleid      string
	Gender      string
	Avatar      string
	Phone       string
	Email       string
	Stat        string
	CreateAt    string
	LastLoginAt string
}{
	UID:         "uid",
	Uname:       "uname",
	Passwd:      "passwd",
	Nickname:    "nickname",
	Roleid:      "roleid",
	Gender:      "gender",
	Avatar:      "avatar",
	Phone:       "phone",
	Email:       "email",
	Stat:        "stat",
	CreateAt:    "create_at",
	LastLoginAt: "last_login_at",
}
