package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/givetrack/givetrack/internal/model"
	"github.com/givetrack/givetrack/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务。
type AuthService struct {
	userRepo  *repository.UserRepository
	orgRepo   *repository.OrganizationRepository
	jwtSecret []byte
	expire    time.Duration
	logger    *slog.Logger
}

func NewAuthService(userRepo *repository.UserRepository, orgRepo *repository.OrganizationRepository, jwtSecret string, expireHours int, logger *slog.Logger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		orgRepo:   orgRepo,
		jwtSecret: []byte(jwtSecret),
		expire:    time.Duration(expireHours) * time.Hour,
		logger:    logger,
	}
}

// Claims JWT 载荷。
type Claims struct {
	UserID uint   `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RegisterInput 注册入参。
type RegisterInput struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=72"`
	Role     string `json:"role" validate:"omitempty,oneof=user org"`
	RealName string `json:"realName"`
	Phone    string `json:"phone"`
}

// Register 注册用户（org 角色自动创建待审核组织档案）。
func (s *AuthService) Register(in RegisterInput) (*model.User, *model.Organization, string, error) {
	existing, err := s.userRepo.FindByUsernameOrEmail(in.Username, in.Email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, nil, "", err
	}
	if existing != nil {
		return nil, nil, "", fmt.Errorf("username or email already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, "", fmt.Errorf("hash password: %w", err)
	}
	role := in.Role
	if role == "" {
		role = "user"
	}
	user := &model.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: string(hash),
		Role:         role,
		RealName:     in.RealName,
		Phone:        in.Phone,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, "", err
	}
	var org *model.Organization
	if role == "org" {
		name := in.RealName
		if name == "" {
			name = in.Username
		}
		org = &model.Organization{UserID: user.ID, Name: name, Status: "pending"}
		if err := s.orgRepo.Create(org); err != nil {
			return nil, nil, "", err
		}
	}
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, nil, "", err
	}
	s.logger.Info("user registered", "username", user.Username, "role", user.Role)
	return user, org, token, nil
}

// Login 登录。
func (s *AuthService) Login(username, password string) (*model.User, *model.Organization, string, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(username, username)
	if err != nil {
		return nil, nil, "", fmt.Errorf("invalid username or password")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, nil, "", fmt.Errorf("invalid username or password")
	}
	var org *model.Organization
	if user.Role == "org" {
		org, _ = s.orgRepo.FindByUserID(user.ID)
	}
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, nil, "", err
	}
	s.logger.Info("user logged in", "username", user.Username)
	return user, org, token, nil
}

// GenerateToken 签发 JWT。
func (s *AuthService) GenerateToken(user *model.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expire)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

// ParseToken 解析 JWT。
func (s *AuthService) ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return claims, nil
}

// UpdateProfile 更新个人资料。
func (s *AuthService) UpdateProfile(userID uint, realName, phone, avatar string) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if realName != "" {
		user.RealName = realName
	}
	if phone != "" {
		user.Phone = phone
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}
