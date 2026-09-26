package service

import (
	"log/slog"

	"github.com/givetrack/givetrack/internal/model"
	"github.com/givetrack/givetrack/internal/repository"
)

// RankingService 排行榜服务。
type RankingService struct {
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

func NewRankingService(userRepo *repository.UserRepository, logger *slog.Logger) *RankingService {
	return &RankingService{userRepo: userRepo, logger: logger}
}

// RankItem 排行条目。
type RankItem struct {
	Rank         int     `json:"rank"`
	UserID       uint    `json:"userId"`
	Username     string  `json:"username"`
	RealName     string  `json:"realName"`
	Avatar       string  `json:"avatar"`
	TotalDonation float64 `json:"totalDonation"`
	ServiceHours float64 `json:"serviceHours"`
}

// DonationRanking 捐款金额排行。
func (s *RankingService) DonationRanking(limit int) ([]RankItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	users, err := s.userRepo.RankingTopDonation(limit)
	if err != nil {
		return nil, err
	}
	return s.toItems(users), nil
}

// ServiceRanking 服务时长排行。
func (s *RankingService) ServiceRanking(limit int) ([]RankItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	users, err := s.userRepo.RankingTopService(limit)
	if err != nil {
		return nil, err
	}
	return s.toItems(users), nil
}

func (s *RankingService) toItems(users []model.User) []RankItem {
	out := make([]RankItem, 0, len(users))
	for i, u := range users {
		out = append(out, RankItem{
			Rank:          i + 1,
			UserID:        u.ID,
			Username:      u.Username,
			RealName:      u.RealName,
			Avatar:        u.Avatar,
			TotalDonation: u.TotalDonation,
			ServiceHours:  u.ServiceHours,
		})
	}
	return out
}

// Stats 平台统计。
func (s *RankingService) Stats() (map[string]interface{}, error) {
	totalUsers, totalDonation, totalHours, err := s.userRepo.Stats()
	if err != nil {
		return nil, err
	}
	s.logger.Info("ranking stats", "users", totalUsers)
	return map[string]interface{}{
		"totalUsers":          totalUsers,
		"totalDonationAmount": totalDonation,
		"totalServiceHours":   totalHours,
	}, nil
}
