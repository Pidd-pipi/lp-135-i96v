package database

import (
	"fmt"
	"log/slog"

	"github.com/givetrack/givetrack/internal/model"
	"gorm.io/gorm"
)

// Seed 幂等种子数据。
func Seed(db *gorm.DB) error {
	var adminCount int64
	if err := db.Model(&model.User{}).Where("role = ?", "admin").Count(&adminCount).Error; err != nil {
		return fmt.Errorf("count admin: %w", err)
	}
	if adminCount == 0 {
		admin := &model.User{Username: "admin", Email: "admin@givetrack.cn", PasswordHash: hashPwd("admin123"), Role: "admin", RealName: "平台管理员"}
		if err := db.Create(admin).Error; err != nil {
			return fmt.Errorf("seed admin: %w", err)
		}
		slog.Info("seeded admin")
	}
	var orgCount int64
	if err := db.Model(&model.User{}).Where("role = ?", "org").Count(&orgCount).Error; err != nil {
		return fmt.Errorf("count org: %w", err)
	}
	if orgCount == 0 {
		orgUser := &model.User{Username: "careorg", Email: "org@givetrack.cn", PasswordHash: hashPwd("org123"), Role: "org", RealName: "阳光公益服务中心"}
		if err := db.Create(orgUser).Error; err != nil {
			return fmt.Errorf("seed org user: %w", err)
		}
		org := &model.Organization{UserID: orgUser.ID, Name: "阳光公益服务中心", Description: "致力于助学、助老与救灾的公益组织", Status: "approved"}
		if err := db.Create(org).Error; err != nil {
			return fmt.Errorf("seed org: %w", err)
		}
		orgUser2 := &model.User{Username: "neworg", Email: "neworg@givetrack.cn", PasswordHash: hashPwd("org123"), Role: "org", RealName: "新起点志愿者协会"}
		if err := db.Create(orgUser2).Error; err != nil {
			return fmt.Errorf("seed org2 user: %w", err)
		}
		org2 := &model.Organization{UserID: orgUser2.ID, Name: "新起点志愿者协会", Description: "待审核组织", Status: "pending"}
		if err := db.Create(org2).Error; err != nil {
			return fmt.Errorf("seed org2: %w", err)
		}
		projects := []*model.Project{
			{OrganizationID: org.ID, Title: "山区小学图书角建设", Description: "为 30 所山区小学建设图书角，配置课外读物与阅读桌椅。", Category: "education", TargetAmount: 100000, CurrentAmount: 65000, ExecutionPlan: "第一期采购图书 5000 册；第二期配送安装书架。", Status: "approved", CoverImage: ""},
			{OrganizationID: org.ID, Title: "独居老人暖心陪伴计划", Description: "组织志愿者定期探访独居老人，提供生活照料与陪伴。", Category: "elderly", TargetAmount: 50000, CurrentAmount: 50000, ExecutionPlan: "每月两次探访，覆盖 200 位老人。", Status: "completed", CoverImage: ""},
			{OrganizationID: org.ID, Title: "乡村儿童助学金项目", Description: "为乡村困境儿童提供学年助学金。", Category: "education", TargetAmount: 200000, CurrentAmount: 30000, ExecutionPlan: "按学年发放助学金。", Status: "approved", CoverImage: ""},
			{OrganizationID: org.ID, Title: "灾害应急物资储备", Description: "储备应急救灾物资，快速响应灾害。", Category: "disaster", TargetAmount: 150000, CurrentAmount: 0, ExecutionPlan: "采购应急物资并建立仓储。", Status: "pending", CoverImage: ""},
		}
		for _, p := range projects {
			if err := db.Create(p).Error; err != nil {
				return fmt.Errorf("seed project: %w", err)
			}
		}
		users := []*model.User{
			{Username: "donor1", Email: "donor1@givetrack.cn", PasswordHash: hashPwd("user123"), Role: "user", RealName: "李明", TotalDonation: 12000, ServiceHours: 36},
			{Username: "donor2", Email: "donor2@givetrack.cn", PasswordHash: hashPwd("user123"), Role: "user", RealName: "王芳", TotalDonation: 8800, ServiceHours: 52},
			{Username: "donor3", Email: "donor3@givetrack.cn", PasswordHash: hashPwd("user123"), Role: "user", RealName: "张伟", TotalDonation: 5600, ServiceHours: 20},
		}
		var donorIDs []uint
		for _, u := range users {
			if err := db.Create(u).Error; err != nil {
				return fmt.Errorf("seed user: %w", err)
			}
			donorIDs = append(donorIDs, u.ID)
		}
		donationSpec := []struct {
			userIdx int
			projID  uint
			amount  float64
			method  string
		}{
			{0, 1, 5000, "wechat"}, {1, 1, 3000, "alipay"}, {2, 2, 2000, "wechat"},
			{0, 2, 1000, "bank"}, {1, 3, 1500, "alipay"},
		}
		for i, ds := range donationSpec {
			d := &model.Donation{
				UserID: donorIDs[ds.userIdx], ProjectID: ds.projID, Amount: ds.amount,
				PaymentMethod: ds.method, PaymentStatus: "success",
				CertificateNo: fmt.Sprintf("CERT%06d", 100001+i),
				TransactionID: fmt.Sprintf("TXN%06d", 100001+i),
			}
			if err := db.Create(d).Error; err != nil {
				return fmt.Errorf("seed donation: %w", err)
			}
		}
		updates := []*model.ProjectUpdate{
			{ProjectID: 1, Title: "首批图书采购完成", Content: "已完成 3000 册图书采购，进入配送阶段。", Images: ""},
			{ProjectID: 2, Title: "5 月探访活动顺利开展", Content: "本月完成 60 位老人探访。", Images: ""},
		}
		for _, u := range updates {
			if err := db.Create(u).Error; err != nil {
				return fmt.Errorf("seed update: %w", err)
			}
		}
		slog.Info("seeded demo data")
	}
	return nil
}

func hashPwd(pwd string) string {
	h, err := bcryptHash(pwd)
	if err != nil {
		slog.Error("hash pwd", "err", err)
		return ""
	}
	return h
}
