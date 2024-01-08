package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int32          `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"column:add_time"`
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_time"`
	IsDeleted bool
}

/*
1. 密文 2. 密文不可反解
	1. 对称加密
	2. 非对称加密
	3. md5 信息摘要算法（md5(SHA-128)已不安全，推荐使用SHA-2(SHA-256、SHA-512)散列函数）
	密码如果不可以反解，用户如何找回密码
*/
type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
}
