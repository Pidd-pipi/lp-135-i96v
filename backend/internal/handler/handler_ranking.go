package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/service"
	"github.com/givetrack/givetrack/internal/util"
)

// RankingHandler 排行榜处理器。
type RankingHandler struct {
	rankingSvc *service.RankingService
}

func NewRankingHandler(rankingSvc *service.RankingService) *RankingHandler {
	return &RankingHandler{rankingSvc: rankingSvc}
}

// DonationRanking 捐款金额排行。
func (h *RankingHandler) DonationRanking(c *gin.Context) {
	list, err := h.rankingSvc.DonationRanking(atoi(c.Query("limit")))
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"rankings": list})
}

// ServiceRanking 服务时长排行。
func (h *RankingHandler) ServiceRanking(c *gin.Context) {
	list, err := h.rankingSvc.ServiceRanking(atoi(c.Query("limit")))
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"rankings": list})
}

// Stats 平台统计。
func (h *RankingHandler) Stats(c *gin.Context) {
	stats, err := h.rankingSvc.Stats()
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"stats": stats})
}
