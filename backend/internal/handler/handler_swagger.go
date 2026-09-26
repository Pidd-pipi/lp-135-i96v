package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const swaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "公益捐赠追踪平台 API",
    "description": "公益组织发布项目、审核、展示、捐款、电子凭证、筹款进度、执行动态与排行榜。",
    "version": "1.0.0"
  },
  "basePath": "/api/v1",
  "schemes": ["http"],
  "paths": {
    "/auth/register": { "post": { "summary": "注册", "tags": ["auth"] } },
    "/auth/login": { "post": { "summary": "登录", "tags": ["auth"] } },
    "/auth/me": { "get": { "summary": "当前用户", "tags": ["auth"] } },
    "/auth/profile": { "put": { "summary": "更新资料", "tags": ["auth"] } },
    "/projects": { "get": { "summary": "项目列表", "tags": ["project"] }, "post": { "summary": "发布项目", "tags": ["project"] } },
    "/projects/org/my": { "get": { "summary": "我的项目", "tags": ["project"] } },
    "/projects/{id}": { "get": { "summary": "项目详情", "tags": ["project"] } },
    "/projects/{id}/updates": { "get": { "summary": "项目进展", "tags": ["project"] }, "post": { "summary": "上传进展", "tags": ["project"] } },
    "/donations": { "post": { "summary": "捐款", "tags": ["donation"] } },
    "/donations/my": { "get": { "summary": "我的捐赠", "tags": ["donation"] } },
    "/donations/{id}/certificate": { "get": { "summary": "电子凭证", "tags": ["donation"] } },
    "/ranking/donation": { "get": { "summary": "捐款排行", "tags": ["ranking"] } },
    "/ranking/service": { "get": { "summary": "服务时长排行", "tags": ["ranking"] } },
    "/ranking/stats": { "get": { "summary": "平台统计", "tags": ["ranking"] } },
    "/admin/projects/pending": { "get": { "summary": "待审核项目", "tags": ["admin"] } },
    "/admin/projects/{id}/review": { "post": { "summary": "审核项目", "tags": ["admin"] } },
    "/admin/organizations/pending": { "get": { "summary": "待审核组织", "tags": ["admin"] } },
    "/admin/organizations/{id}/review": { "post": { "summary": "审核组织", "tags": ["admin"] } }
  }
}`

// SwaggerJSON 提供 swagger.json。
func (h *HealthHandler) SwaggerJSON(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, swaggerJSON)
}
