package service

import (
	"context"
	"fmt"
	"io"
	"messenger/backend/internal/domain"
	"messenger/backend/internal/repository"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	users     repository.UserRepository
	avatarDir string
}

func NewUserService(users repository.UserRepository, avatarDir string) *UserService {
	return &UserService{users: users, avatarDir: avatarDir}
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) GetPublicByID(ctx context.Context, id string) (*domain.UserPublic, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	pub := u.ToPublic()
	return &pub, nil
}

type UpdateUserInput struct {
	Username string
	Bio      string
}

func (s *UserService) UpdateMe(ctx context.Context, userID string, in UpdateUserInput) (*domain.User, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("user not found")
	}
	u.Username = in.Username
	u.Bio = in.Bio
	u.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) UploadAvatar(ctx context.Context, userID string, r io.Reader, ext string) (string, error) {
	if err := os.MkdirAll(s.avatarDir, 0755); err != nil {
		return "", fmt.Errorf("create avatar dir: %w", err)
	}
	filename := uuid.New().String() + ext
	path := filepath.Join(s.avatarDir, filename)
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create avatar file: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return "", fmt.Errorf("write avatar: %w", err)
	}

	avatarURL := "/avatars/" + filename
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	u.AvatarURL = avatarURL
	u.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, u); err != nil {
		return "", err
	}
	return avatarURL, nil
}

func (s *UserService) Search(ctx context.Context, query string) ([]domain.UserPublic, error) {
	if len(query) < 1 {
		return nil, nil
	}
	return s.users.Search(ctx, query, 20)
}
