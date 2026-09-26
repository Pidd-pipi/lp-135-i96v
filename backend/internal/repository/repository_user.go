package repository

import (
	"errors"
	"fmt"

	"github.com/givetrack/givetrack/internal/model"
	"gorm.io/gorm"
)

// ErrNotFound 哨兵错误。
var ErrNotFound = errors.New("record not found")

// UserRepository 用户数据访问。
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *model.User) error {
	if err := r.db.Create(u).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByUsernameOrEmail(username, email string) (*model.User, error) {
	var u model.User
	err := r.db.Where("username = ? OR email = ?", username, email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var u model.User
	err := r.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) Update(u *model.User) error {
	if err := r.db.Save(u).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// RankingTopDonation 按累计捐款金额排行。
func (r *UserRepository) RankingTopDonation(limit int) ([]model.User, error) {
	var list []model.User
	if err := r.db.Where("role = ? AND total_donation > 0", "user").
		Order("total_donation DESC").Limit(limit).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("ranking donation: %w", err)
	}
	return list, nil
}

// RankingTopService 按服务时长排行。
func (r *UserRepository) RankingTopService(limit int) ([]model.User, error) {
	var list []model.User
	if err := r.db.Where("role = ? AND service_hours > 0", "user").
		Order("service_hours DESC").Limit(limit).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("ranking service: %w", err)
	}
	return list, nil
}

// Stats 平台统计。
func (r *UserRepository) Stats() (totalUsers int64, totalDonation float64, totalServiceHours float64, err error) {
	var u struct {
		TotalUsers  int64
		TotalDon    float64
		TotalHours  float64
	}
	if err := r.db.Model(&model.User{}).Where("role = ?", "user").
		Select("COUNT(*) AS total_users, COALESCE(SUM(total_donation),0) AS total_don, COALESCE(SUM(service_hours),0) AS total_hours").
		Scan(&u).Error; err != nil {
		return 0, 0, 0, fmt.Errorf("stats: %w", err)
	}
	return u.TotalUsers, u.TotalDon, u.TotalHours, nil
}

// OrganizationRepository 组织数据访问。
type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(o *model.Organization) error {
	if err := r.db.Create(o).Error; err != nil {
		return fmt.Errorf("create organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) FindByUserID(userID uint) (*model.Organization, error) {
	var o model.Organization
	err := r.db.Where("user_id = ?", userID).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find organization by user id: %w", err)
	}
	return &o, nil
}

func (r *OrganizationRepository) FindByID(id uint) (*model.Organization, error) {
	var o model.Organization
	err := r.db.First(&o, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find organization by id: %w", err)
	}
	return &o, nil
}

func (r *OrganizationRepository) Update(o *model.Organization) error {
	if err := r.db.Save(o).Error; err != nil {
		return fmt.Errorf("update organization: %w", err)
	}
	return nil
}

// FindPending 待审核组织列表。
func (r *OrganizationRepository) FindPending() ([]model.Organization, error) {
	var list []model.Organization
	if err := r.db.Preload("User").Where("status = ?", "pending").
		Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("find pending organizations: %w", err)
	}
	return list, nil
}
