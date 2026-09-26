package model

import "time"

// Project 公益项目。
type Project struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	OrganizationID uint         `gorm:"index;not null" json:"organizationId"`
	Title          string       `gorm:"size:200;not null" json:"title"`
	Description    string       `gorm:"type:text" json:"description"`
	Category       string       `gorm:"size:32;index;not null" json:"category"`
	TargetAmount   float64      `gorm:"type:decimal(14,2);not null" json:"targetAmount"`
	CurrentAmount  float64      `gorm:"type:decimal(14,2);default:0" json:"currentAmount"`
	ExecutionPlan  string       `gorm:"type:text" json:"executionPlan"`
	CoverImage     string       `gorm:"size:255" json:"coverImage"`
	Status         string       `gorm:"size:20;index;default:pending" json:"status"`
	StartDate      *time.Time   `json:"startDate"`
	EndDate        *time.Time   `json:"endDate"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
}

// ProjectUpdate 项目执行进展。
type ProjectUpdate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"index;not null" json:"projectId"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Images    string    `gorm:"type:text" json:"images"`
	CreatedAt time.Time `json:"createdAt"`
}

// Donation 捐赠记录。
type Donation struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"index;not null" json:"userId"`
	ProjectID      uint      `gorm:"index;not null" json:"projectId"`
	Amount         float64   `gorm:"type:decimal(14,2);not null" json:"amount"`
	PaymentMethod  string    `gorm:"size:20" json:"paymentMethod"`
	PaymentStatus  string    `gorm:"size:20;index;default:success" json:"paymentStatus"`
	TransactionID  string    `gorm:"size:64" json:"transactionId"`
	CertificateNo  string    `gorm:"size:64" json:"certificateNo"`
	IsAnonymous    bool      `gorm:"default:false" json:"isAnonymous"`
	Message        string    `gorm:"size:255" json:"message"`
	CreatedAt      time.Time `json:"createdAt"`
	User           *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Project        *Project  `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// AdminReview 审核记录。
type AdminReview struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProjectID      uint      `gorm:"index" json:"projectId"`
	OrganizationID uint      `gorm:"index" json:"organizationId"`
	ReviewerID     uint      `json:"reviewerId"`
	Status         string    `gorm:"size:20" json:"status"`
	Comment        string    `gorm:"size:255" json:"comment"`
	CreatedAt      time.Time `json:"createdAt"`
}

// VolunteerService 志愿服务时长记录（用于服务时长排行榜）。
type VolunteerService struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"userId"`
	ProjectID uint      `gorm:"index" json:"projectId"`
	Hours     float64   `gorm:"type:decimal(8,2);not null" json:"hours"`
	ServiceAt time.Time `json:"serviceAt"`
	CreatedAt time.Time `json:"createdAt"`
}
