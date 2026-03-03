package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/yujiawei/nexus-mm/internal/model"
	"github.com/yujiawei/nexus-mm/internal/store/postgres"
	"github.com/yujiawei/nexus-mm/internal/wkim"
)

type UserService struct {
	store      *postgres.UserStore
	wk         *wkim.Client
	jwtSecret  []byte
	jwtExpireH int
}

func NewUserService(store *postgres.UserStore, wk *wkim.Client, jwtSecret string, jwtExpireH int) *UserService {
	return &UserService{
		store:      store,
		wk:         wk,
		jwtSecret:  []byte(jwtSecret),
		jwtExpireH: jwtExpireH,
	}
}

func (s *UserService) Register(ctx context.Context, req *model.UserRegister) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// First registered user becomes system admin.
	role := "member"
	count, err := s.store.Count(ctx)
	if err == nil && count == 0 {
		role = "admin"
	}

	now := time.Now().UTC()
	user := &model.User{
		ID:           ulid.Make().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Nickname:     req.Nickname,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.store.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Register user in WuKongIM with a random token.
	wkToken := randomHex(16)
	if err := s.wk.RegisterUser(ctx, user.ID, wkToken); err != nil {
		// Log but don't fail registration - WuKongIM might be down.
		fmt.Printf("warn: register user in wukongim: %v\n", err)
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req *model.UserLogin) (*model.LoginResponse, error) {
	user, err := s.store.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &model.LoginResponse{Token: token, User: user}, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	return s.store.GetByID(ctx, id)
}

func (s *UserService) UpdateRole(ctx context.Context, userID, role string) error {
	return s.store.UpdateRole(ctx, userID, role)
}

func (s *UserService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(s.jwtExpireH) * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
