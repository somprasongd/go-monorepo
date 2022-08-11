package service

import (
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/common/logger"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/config"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/auth/core/dto"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/auth/core/ports"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/model"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/util"
)

var (
	ErrUserEmailDuplication = common.NewBadRequestError("email already exists")
	ErrHashPassword         = common.NewUnexpectedError("error occurred while hashing password")
	ErrLogin                = common.NewUnauthorizedError("the email or password are incorrect")
	ErrGenerateAccessToken  = common.NewUnexpectedError("error occurred while generating token")
	ErrGenerateRefreshToken = common.NewUnexpectedError("error occurred while generating refresh token")
	ErrValidateToken        = common.NewUnexpectedError("error occurred while validating token")
	ErrNoToken              = common.NewUnauthorizedError("the token is required")
	ErrInvalidToken         = common.NewUnauthorizedError("the token is invalid")
	ErrUserNotfound         = common.NewUnauthorizedError("user not found")
	ErrUserPasswordNotMatch = common.NewBadRequestError("password is not match")
	ErrInvalidRefreshToken  = common.NewUnauthorizedError("the refresh token is invalid or expired")
	ErrUnmarshalPayload     = common.NewUnexpectedError("error occurred while convert user data")
)

type authService struct {
	config    *config.Config
	repo      ports.AuthRepository
	tokenRepo ports.TokenRepository
}

func NewAuthService(config *config.Config, repo ports.AuthRepository, tokenRepo ports.TokenRepository) ports.AuthService {
	return &authService{config, repo, tokenRepo}
}

func (s authService) Register(form dto.RegisterForm, log logger.Interface) error {
	// validate
	if err := common.ValidateDto(form); err != nil {
		return common.NewInvalidError(err.Error())
	}

	u, err := s.repo.FindUserByEmail(form.Email)

	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		log.Error(err.Error())
		return common.ErrDbQuery
	}

	if u != nil {
		return ErrUserEmailDuplication
	}

	auth := model.User{Email: form.Email}
	hashPwd, err := util.HashPasswordArgon(form.Password)
	if err != nil {
		log.Error(err.Error())
		return ErrHashPassword
	}
	auth.Password = hashPwd

	err = s.repo.CreateUser(&auth)
	if err != nil {
		log.Error(err.Error())
		return common.ErrDbInsert
	}

	return nil
}

func (s authService) Login(form dto.LoginForm, log logger.Interface) (*dto.AuthResponse, error) {
	// validate form
	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}
	// ค้นหาจาก email
	user, err := s.repo.FindUserByEmail(form.Email)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrLogin
		}
		return nil, common.ErrDbQuery
	}
	// ตรวจสอบรหัสผ่าน ตรงกันหรือไม่
	match := util.CheckPasswordHashArgon(form.Password, user.Password)

	if !match {
		return nil, ErrLogin
	}

	tokenId, _ := uuid.NewV4()
	payload := map[string]any{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role.String(),
	}

	// สร้าง refresh token
	refreshToken, refreshExpiresAt, err := util.GenerateToken(tokenId.String(), nil, s.config.Token.RefreshSecretKey, s.config.Token.RefreshExpires)

	if err != nil {
		log.Error(err.Error())
		return nil, ErrGenerateRefreshToken
	}

	// บันทึก user ลง redis
	err = s.tokenRepo.SetToken(tokenId.String(), payload, s.config.Token.RefreshExpires)
	if err != nil {
		log.Error(err.Error())
		return nil, ErrGenerateRefreshToken
	}

	// สร้าง access token
	accessToken, accessExpiresAt, err := util.GenerateToken(tokenId.String(), payload, s.config.Token.AccessSecretKey, s.config.Token.AccessExpires)

	if err != nil {
		log.Error(err.Error())
		return nil, ErrGenerateAccessToken
	}

	// ตอบกลับไปพร้อมข้อมูล user
	serialized := dto.AuthResponse{
		User: &dto.UserInfo{
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role.String(),
		},
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}
	return &serialized, nil
}

func (s authService) Profile(email string, log logger.Interface) (*dto.UserInfo, error) {
	// validate
	if email == "" {
		return nil, ErrUserNotfound
	}

	user, err := s.repo.FindUserByEmail(email)

	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotfound
		}
		log.Error(err.Error())
		return nil, common.ErrDbQuery
	}

	serialized := dto.UserInfo{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role.String(),
	}

	return &serialized, nil
}

func (s authService) UpdateProfile(email string, form dto.UpdateProfileForm, log logger.Interface) (*dto.UserInfo, error) {
	// validate
	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotfound
		}
		log.Error(err.Error())
		return nil, common.ErrDbQuery
	}

	match := util.CheckPasswordHash(form.PasswordOld, user.Password)

	if !match {
		return nil, ErrUserPasswordNotMatch
	}

	hashPwd, err := util.HashPassword(form.PasswordNew)

	if err != nil {
		log.Error(err.Error())
		return nil, ErrHashPassword
	}

	user.Password = hashPwd

	err = s.repo.SaveProfile(user)
	if err != nil {
		log.Error(err.Error())
		return nil, common.ErrDbUpdate
	}

	serialized := dto.UserInfo{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role.String(),
	}

	return &serialized, nil
}

func (s authService) RefreshToken(form dto.RefreshForm, log logger.Interface) (*dto.AuthResponse, error) {
	// validate form
	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	// ตรวจสอบ refresh token ว่ายัง valid หรือไม่
	cliams, err := util.ValidateToken(form.Token, s.config.Token.RefreshSecretKey)

	if err != nil {
		log.Error(err.Error())
		return nil, ErrInvalidRefreshToken
	}

	// เอา token id ไปหาใน redis
	tokenId := cliams["sub"].(string)
	encodedUser, err := s.tokenRepo.GetToken(tokenId)

	if err != nil {
		log.Error(err.Error())
		return nil, ErrInvalidRefreshToken
	}

	// ถ้าอ่านค่าได้เป็นค่าว่าง แสดงว่าหมดอายุแล้ว
	if encodedUser == "" {
		return nil, ErrInvalidRefreshToken
	}

	// อ่านค่า user ออกมาไว้ใน payload
	payload := map[string]any{}
	err = json.Unmarshal([]byte(encodedUser), &payload)
	if err != nil {
		log.Error(err.Error())
		return nil, ErrUnmarshalPayload
	}

	// สร้าง tokenId ใหม่
	newTkId, _ := uuid.NewV4()

	// สร้าง refresh token ใหม่
	refreshToken, refreshExpiresAt, err := util.GenerateToken(newTkId.String(), nil, s.config.Token.RefreshSecretKey, s.config.Token.RefreshExpires)

	if err != nil {
		log.Error(err.Error())
		return nil, ErrGenerateRefreshToken
	}

	// บันทึก user ลง redis
	err = s.tokenRepo.SetToken(newTkId.String(), payload, s.config.Token.RefreshExpires)
	if err != nil {
		log.Error(err.Error())
		return nil, ErrGenerateRefreshToken
	}
	// ลบ tokenId เดิม ป้องกันใช้ซ้ำ
	s.tokenRepo.DeleteToken(tokenId)

	// สร้าง access token ใหม่
	accessToken, accessExpiresAt, err := util.GenerateToken(newTkId.String(), payload, s.config.Token.AccessSecretKey, s.config.Token.AccessExpires)

	if err != nil {
		log.Error(err.Error())
		return nil, ErrGenerateAccessToken
	}

	// ส่ง token ใหม่กลับไป
	serialized := dto.AuthResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}
	return &serialized, nil
}

func (s authService) RevokeToken(form dto.RefreshForm, log logger.Interface) error {
	// validate form
	err := common.ValidateDto(form)
	if err != nil {
		return common.NewInvalidError(err.Error())
	}

	// ตรวจสอบ refresh token ว่ายัง valid หรือไม่
	cliams, err := util.ValidateToken(form.Token, s.config.Token.RefreshSecretKey)

	if err != nil {
		log.Error(err.Error())
		return ErrInvalidRefreshToken
	}

	// เอา token id ไปหาใน redis
	tokenId := cliams["sub"].(string)
	s.tokenRepo.DeleteToken(tokenId)

	return nil
}
