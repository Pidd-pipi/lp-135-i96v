package model

import "time"

// User 用户（个人/公益组织/管理员）。
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"size:128;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:20;index;not null" json:"role"`
	RealName     string    `gorm:"size:64" json:"realName"`
	Avatar       string    `gorm:"size:255" json:"avatar"`
	Phone        string    `gorm:"size:32" json:"phone"`
	TotalDonation float64  `gorm:"type:decimal(14,2);default:0" json:"totalDonation"`
	ServiceHours float64   `gorm:"type:decimal(10,2);default:0" json:"serviceHours"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Organization 公益组织档案。
type Organization struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"uniqueIndex;not null" json:"userId"`
	Name          string    `gorm:"size:128;not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	LicenseNumber string    `gorm:"size:64" json:"licenseNumber"`
	ContactPerson string    `gorm:"size:64" json:"contactPerson"`
	ContactPhone  string    `gorm:"size:32" json:"contactPhone"`
	Address       string    `gorm:"size:255" json:"address"`
	Status        string    `gorm:"size:20;index;default:pending" json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	User          *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
