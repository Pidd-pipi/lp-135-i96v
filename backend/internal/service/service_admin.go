package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/givetrack/givetrack/internal/constants"
	"github.com/givetrack/givetrack/internal/model"
	"github.com/givetrack/givetrack/internal/repository"
)

// AdminService 平台审核服务。
type AdminService struct {
	projectRepo *repository.ProjectRepository
	orgRepo     *repository.OrganizationRepository
	reviewRepo  *repository.AdminReviewRepository
	logger      *slog.Logger
}

func NewAdminService(projectRepo *repository.ProjectRepository, orgRepo *repository.OrganizationRepository, reviewRepo *repository.AdminReviewRepository, logger *slog.Logger) *AdminService {
	return &AdminService{projectRepo: projectRepo, orgRepo: orgRepo, reviewRepo: reviewRepo, logger: logger}
}

// PendingProjects 待审核项目列表。
func (s *AdminService) PendingProjects() ([]model.Project, error) {
	return s.projectRepo.FindPending()
}

// ReviewProject 审核项目。
func (s *AdminService) ReviewProject(adminID, projectID uint, status, comment string) (*model.Project, error) {
	if status != constants.ProjectApproved && status != constants.ProjectRejected {
		return nil, fmt.Errorf("invalid review status")
	}
	p, err := s.projectRepo.FindByID(projectID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	p.Status = status
	if err := s.projectRepo.Update(p); err != nil {
		return nil, err
	}
	if err := s.reviewRepo.Create(&model.AdminReview{
		ProjectID:  projectID,
		ReviewerID: adminID,
		Status:     status,
		Comment:    comment,
	}); err != nil {
		return nil, err
	}
	s.logger.Info("project reviewed", "projectId", projectID, "status", status)
	return p, nil
}

// PendingOrganizations 待审核组织列表。
func (s *AdminService) PendingOrganizations() ([]model.Organization, error) {
	return s.orgRepo.FindPending()
}

// ReviewOrganization 审核组织。
func (s *AdminService) ReviewOrganization(adminID, orgID uint, status, comment string) (*model.Organization, error) {
	if status != constants.OrgApproved && status != constants.OrgRejected {
		return nil, fmt.Errorf("invalid review status")
	}
	o, err := s.orgRepo.FindByID(orgID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	o.Status = status
	if err := s.orgRepo.Update(o); err != nil {
		return nil, err
	}
	if err := s.reviewRepo.Create(&model.AdminReview{
		OrganizationID: orgID,
		ReviewerID:     adminID,
		Status:         status,
		Comment:        comment,
	}); err != nil {
		return nil, err
	}
	s.logger.Info("organization reviewed", "orgId", orgID, "status", status)
	return o, nil
}
