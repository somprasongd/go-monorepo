package ports

import (
	"time"

	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/auth/core/dto"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/model"
)

// interface สำหรับ output port
type AuthRepository interface {
	FindUserByEmail(email string) (*model.User, error)
	CreateUser(*model.User) error
	SaveProfile(m *model.User) error
}

type TokenRepository interface {
	SetToken(tokenId string, data map[string]any, duration time.Duration) error
	GetToken(tokenId string) (string, error)
	DeleteToken(tokenId string) (int64, error)
}

// interface สำหรับ input port
type AuthService interface {
	Register(form dto.RegisterForm, reqId string) error
	Login(form dto.LoginForm, reqId string) (*dto.AuthResponse, error)
	Profile(email string, reqId string) (*dto.UserInfo, error)
	UpdateProfile(email string, form dto.UpdateProfileForm, reqId string) (*dto.UserInfo, error)
	RefreshToken(form dto.RefreshForm, reqId string) (*dto.AuthResponse, error)
	RevokeToken(form dto.RefreshForm, reqId string) error
}
