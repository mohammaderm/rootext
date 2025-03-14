package userService

import (
	"context"
	"errors"

	"github.com/mohammaderm/rootext/entity"
	"github.com/mohammaderm/rootext/params"
	"github.com/mohammaderm/rootext/repository/postgres/userRepository"
	"github.com/mohammaderm/rootext/service/authService"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo    userRepository.UserRepo
	authSvc authService.AuthService
}

func New(userRepo userRepository.UserRepo, authSvc authService.AuthService) Service {
	return Service{
		repo:    userRepo,
		authSvc: authSvc,
	}
}

func (s Service) Login(ctx context.Context, req params.LoginRequest) (params.LoginResponse, error) {

	user, exist, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return params.LoginResponse{}, err
	}
	if !exist {
		return params.LoginResponse{}, errors.New("user is not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return params.LoginResponse{}, errors.New("username or password is not valid")
	}

	accessToken, err := s.authSvc.CreateAccessToken(user)
	if err != nil {
		return params.LoginResponse{}, err
	}
	refreshToken, err := s.authSvc.CreateRefreshToken(user)
	if err != nil {
		return params.LoginResponse{}, err
	}

	return params.LoginResponse{
		User: user,
		Tokens: params.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (s Service) Register(ctx context.Context, req params.RegisterRequest) (params.RegisterResponse, error) {

	if !isusernamevalid(req.Username) {
		return params.RegisterResponse{}, errors.New("user is not valid")
	}
	if !ispasswordvalid(req.Password) {
		return params.RegisterResponse{}, errors.New("password is not valid")
	}
	if ok, err := s.repo.IsUserUnique(ctx, req.Username); err != nil || !ok {
		if err != nil {
			return params.RegisterResponse{}, err
		}
		if !ok {
			return params.RegisterResponse{}, errors.New("user is not unique")
		}
	}
	hashedPassword, err := hashpassword(req.Password)
	if err != nil {
		return params.RegisterResponse{}, err
	}
	user := entity.User{
		Username: req.Username,
		Password: hashedPassword,
	}
	createdUser, err := s.repo.Register(ctx, user)
	if err != nil {
		return params.RegisterResponse{}, err
	}
	return params.RegisterResponse{
		User: createdUser,
	}, nil
}

func (s Service) TokenRenew(ctx context.Context, req params.TokenRenewReq) (params.TokenRenewRes, error) {

	claims, err := s.authSvc.ParseToken(req.RefreshToken)
	if err != nil {
		return params.TokenRenewRes{}, err
	}
	user, ok, err := s.repo.IsUserExistsById(ctx, claims.UserID)

	if err != nil {
		return params.TokenRenewRes{}, err
	}
	if !ok {
		return params.TokenRenewRes{}, errors.New("user is not found")
	}

	accessToken, err := s.authSvc.CreateAccessToken(user)
	if err != nil {
		return params.TokenRenewRes{}, err
	}
	refreshToken, err := s.authSvc.CreateRefreshToken(user)
	if err != nil {
		return params.TokenRenewRes{}, err
	}

	return params.TokenRenewRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func isusernamevalid(username string) bool {
	return len(username) >= 5
}

func ispasswordvalid(password string) bool {
	return len(password) >= 7
}

func hashpassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
