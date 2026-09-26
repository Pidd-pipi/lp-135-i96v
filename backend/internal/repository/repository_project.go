package repository

import (
	"errors"
	"fmt"

	"github.com/givetrack/givetrack/internal/model"
	"gorm.io/gorm"
)

// ProjectRepository 项目数据访问。
type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(p *model.Project) error {
	if err := r.db.Create(p).Error; err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

func (r *ProjectRepository) FindByID(id uint) (*model.Project, error) {
	var p model.Project
	err := r.db.Preload("Organization").Preload("Organization.User").First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find project by id: %w", err)
	}
	return &p, nil
}

func (r *ProjectRepository) Update(p *model.Project) error {
	if err := r.db.Save(p).Error; err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	return nil
}

// List 分页查询项目，支持分类/状态筛选。
func (r *ProjectRepository) List(category, status string, page, pageSize int) ([]model.Project, int64, error) {
	var list []model.Project
	var total int64
	q := r.db.Model(&model.Project{}).Preload("Organization")
	if category != "" && category != "all" {
		q = q.Where("category = ?", category)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count projects: %w", err)
	}
	if err := q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list projects: %w", err)
	}
	return list, total, nil
}

// ListByOrg 组织自己的项目。
func (r *ProjectRepository) ListByOrg(orgID uint) ([]model.Project, error) {
	var list []model.Project
	if err := r.db.Where("organization_id = ?", orgID).
		Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list projects by org: %w", err)
	}
	return list, nil
}

// FindPending 待审核项目列表。
func (r *ProjectRepository) FindPending() ([]model.Project, error) {
	var list []model.Project
	if err := r.db.Preload("Organization").Where("status = ?", "pending").
		Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("find pending projects: %w", err)
	}
	return list, nil
}

// ProjectUpdateRepository 项目进展数据访问。
type ProjectUpdateRepository struct {
	db *gorm.DB
}

func NewProjectUpdateRepository(db *gorm.DB) *ProjectUpdateRepository {
	return &ProjectUpdateRepository{db: db}
}

func (r *ProjectUpdateRepository) Create(u *model.ProjectUpdate) error {
	if err := r.db.Create(u).Error; err != nil {
		return fmt.Errorf("create project update: %w", err)
	}
	return nil
}

func (r *ProjectUpdateRepository) ListByProject(projectID uint) ([]model.ProjectUpdate, error) {
	var list []model.ProjectUpdate
	if err := r.db.Where("project_id = ?", projectID).
		Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list project updates: %w", err)
	}
	return list, nil
}
