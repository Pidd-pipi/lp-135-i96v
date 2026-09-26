package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/service"
	"github.com/givetrack/givetrack/internal/util"
)

// DonationHandler 捐赠处理器。
type DonationHandler struct {
	donationSvc *service.DonationService
}

func NewDonationHandler(donationSvc *service.DonationService) *DonationHandler {
	return &DonationHandler{donationSvc: donationSvc}
}

// Create 捐款。
func (h *DonationHandler) Create(c *gin.Context) {
	var req service.CreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Fail(c, http.StatusBadRequest, 42200, err.Error())
		return
	}
	d, err := h.donationSvc.Create(c.GetUint("user_id"), req)
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.Created(c, gin.H{
		"message": "捐款成功",
		"donation": gin.H{
			"id":            d.ID,
			"amount":        d.Amount,
			"certificateNo": d.CertificateNo,
			"createdAt":     d.CreatedAt,
		},
	})
}

// My 我的捐赠记录。
func (h *DonationHandler) My(c *gin.Context) {
	page, ps := util.NormalizePage(atoi(c.Query("page")), atoi(c.Query("page_size")))
	if c.Query("limit") != "" {
		ps = atoi(c.Query("limit"))
	}
	list, total, err := h.donationSvc.MyDonations(c.GetUint("user_id"), page, ps)
	if err != nil {
		util.FailError(c, err)
		return
	}
	util.OK(c, gin.H{"donations": list, "total": total, "page": page, "limit": ps})
}

// Certificate 电子凭证。
func (h *DonationHandler) Certificate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, 40000, "invalid donation id")
		return
	}
	d, err := h.donationSvc.Certificate(c.GetUint("user_id"), uint(id))
	if err != nil {
		util.FailError(c, err)
		return
	}
	donorName := d.User.RealName
	if donorName == "" {
		donorName = d.User.Username
	}
	util.OK(c, gin.H{"certificate": gin.H{
		"id":            d.ID,
		"certificateNo": d.CertificateNo,
		"amount":        d.Amount,
		"projectTitle":  d.Project.Title,
		"donorName":     donorName,
		"createdAt":     d.CreatedAt,
	}})
}
