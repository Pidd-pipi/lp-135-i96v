package util

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/constants"
	"github.com/givetrack/givetrack/internal/repository"
)

// Response 统一响应体。
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK 成功响应。
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: constants.CodeOK, Message: "ok", Data: data})
}

// Created 创建成功响应。
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Code: constants.CodeOK, Message: "ok", Data: data})
}

// Fail 业务失败响应。
func Fail(c *gin.Context, status, code int, message string) {
	c.JSON(status, Response{Code: code, Message: message})
}

// FailError 根据错误类型转换为统一响应。
func FailError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		Fail(c, http.StatusNotFound, constants.CodeNotFound, err.Error())
	case errors.Is(err, ErrUnauthorized):
		Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, err.Error())
	case errors.Is(err, ErrForbidden):
		Fail(c, http.StatusForbidden, constants.CodeForbidden, err.Error())
	default:
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, err.Error())
	}
}

// ErrUnauthorized 未认证错误。
var ErrUnauthorized = errors.New("unauthorized")

// ErrForbidden 无权限错误。
var ErrForbidden = errors.New("forbidden")

// Page 分页参数。
type Page struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	Limit    int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// NormalizePage 规范化分页参数。
func NormalizePage(p, ps int) (int, int) {
	if p <= 0 {
		p = constants.PageDefault
	}
	if ps <= 0 {
		ps = constants.PageSizeDefault
	}
	if ps > constants.PageSizeMax {
		ps = constants.PageSizeMax
	}
	return p, ps
}
