package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"

	"github.com/givetrack/givetrack/internal/constants"
	"github.com/givetrack/givetrack/internal/model"
	"github.com/givetrack/givetrack/internal/repository"
	"gorm.io/gorm"
)

// DonationService 捐赠服务。
type DonationService struct {
	db           *gorm.DB
	donationRepo *repository.DonationRepository
	projectRepo  *repository.ProjectRepository
	userRepo     *repository.UserRepository
	logger       *slog.Logger
}

func NewDonationService(
	db *gorm.DB,
	donationRepo *repository.DonationRepository,
	projectRepo *repository.ProjectRepository,
	userRepo *repository.UserRepository,
	logger *slog.Logger,
) *DonationService {
	return &DonationService{
		db:           db,
		donationRepo: donationRepo,
		projectRepo:  projectRepo,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// CreateInput 捐款入参。
type CreateInput struct {
	ProjectID     uint    `json:"projectId" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"paymentMethod"`
	IsAnonymous   bool    `json:"isAnonymous"`
	Message       string  `json:"message"`
}

// Create 捐款并生成电子凭证。
// 使用数据库事务保证捐赠流水、项目已筹金额、用户累计捐赠一致写入。
func (s *DonationService) Create(userID uint, in CreateInput) (*model.Donation, error) {
	if in.ProjectID == 0 {
		return nil, fmt.Errorf("projectId is required")
	}
	if in.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	certNo, err := newReferenceNo("CERT")
	if err != nil {
		return nil, err
	}
	txnNo, err := newReferenceNo("TXN")
	if err != nil {
		return nil, err
	}

	var donation *model.Donation
	err = s.db.Transaction(func(tx *gorm.DB) error {
		projectRepo := repository.NewProjectRepository(tx)
		donationRepo := repository.NewDonationRepository(tx)
		userRepo := repository.NewUserRepository(tx)

		project, err := projectRepo.FindByID(in.ProjectID)
		if err != nil {
			return err
		}
		if project.Status != constants.ProjectApproved {
			return fmt.Errorf("project not approved for donation")
		}

		method := in.PaymentMethod
		if method == "" {
			method = "wechat"
		}
		donation = &model.Donation{
			UserID:        userID,
			ProjectID:     in.ProjectID,
			Amount:        in.Amount,
			PaymentMethod: method,
			PaymentStatus: constants.PaymentSuccess,
			IsAnonymous:   in.IsAnonymous,
			Message:       in.Message,
			CertificateNo: certNo,
			TransactionID: txnNo,
		}
		if err := donationRepo.Create(donation); err != nil {
			return err
		}

		project.CurrentAmount += in.Amount
		if project.CurrentAmount >= project.TargetAmount {
			project.Status = constants.ProjectCompleted
		}
		if err := projectRepo.Update(project); err != nil {
			return err
		}

		user, err := userRepo.FindByID(userID)
		if err != nil {
			return err
		}
		user.TotalDonation += in.Amount
		if err := userRepo.Update(user); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info("donation created", "donationId", donation.ID, "amount", in.Amount, "certNo", donation.CertificateNo)
	return donation, nil
}

// MyDonations 我的捐赠记录。
func (s *DonationService) MyDonations(userID uint, page, pageSize int) ([]model.Donation, int64, error) {
	return s.donationRepo.ListByUser(userID, page, pageSize)
}

// Certificate 查询电子凭证。
func (s *DonationService) Certificate(userID, donationID uint) (*model.Donation, error) {
	d, err := s.donationRepo.FindByID(donationID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if d.UserID != userID {
		return nil, fmt.Errorf("forbidden: certificate belongs to another user")
	}
	return d, nil
}

func newReferenceNo(prefix string) (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate reference number: %w", err)
	}
	return prefix + hex.EncodeToString(b), nil
}
