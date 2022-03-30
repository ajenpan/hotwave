package account

import (
	"time"
)

// User [...]
type User struct {
	ID       int64     `gorm:"primaryKey;column:id;type:bigint;not null" json:"-"`
	Name     string    `gorm:"unique;column:name;type:varchar(32);not null" json:"name"`
	Passwd   string    `gorm:"column:passwd;type:varchar(64);not null" json:"passwd"`
	Nickname string    `gorm:"column:nickname;type:varchar(64);not null;default:''" json:"nickname"`
	Role     int       `gorm:"column:role;type:int;not null;default:1" json:"role"`
	Gender   string    `gorm:"column:gender;type:enum('M','F','X');not null;default:X" json:"gender"`
	Avatar   string    `gorm:"column:avatar;type:varchar(1024);not null;default:''" json:"avatar"`
	Phone    string    `gorm:"column:phone;type:varchar(32);not null;default:''" json:"phone"`
	Email    string    `gorm:"column:email;type:varchar(128);not null;default:''" json:"email"`
	Stat     int8      `gorm:"column:stat;type:tinyint;not null;default:0" json:"stat"`
	CreateAt time.Time `gorm:"column:create_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createAt"`
}

// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "user"
}

// UserColumns get sql column name.获取数据库列名
var UserColumns = struct {
	ID       string
	Name     string
	Passwd   string
	Nickname string
	Role     string
	Gender   string
	Avatar   string
	Phone    string
	Email    string
	Stat     string
	CreateAt string
}{
	ID:       "id",
	Name:     "name",
	Passwd:   "passwd",
	Nickname: "nickname",
	Role:     "role",
	Gender:   "gender",
	Avatar:   "avatar",
	Phone:    "phone",
	Email:    "email",
	Stat:     "stat",
	CreateAt: "create_at",
}
