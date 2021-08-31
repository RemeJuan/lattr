package services

import (
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/error_utils"
)

var (
	AuthService authServiceInterface = &authService{}
)

type authService struct{}

type authServiceInterface interface {
	Create(*domain.Token) (*domain.Token, error_utils.MessageErr)
	Get(int64) (*domain.Token, error_utils.MessageErr)
	List() ([]domain.Token, error_utils.MessageErr)
	Reset(*domain.Token) (*domain.Token, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
}

func (as authService) Create(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
	if err := token.Validate(); err != nil {
		return nil, err
	}

	token.CreatedAt = time.Now().Local()
	token.Modified = time.Now().Local()

	tk, err := domain.TokenRepo.Create(token)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

func (as authService) Get(id int64) (*domain.Token, error_utils.MessageErr) {
	token, err := domain.TokenRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (as authService) List() ([]domain.Token, error_utils.MessageErr) {
	token, err := domain.TokenRepo.List()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (as authService) Reset(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
	current, err := domain.TokenRepo.Get(token.Id)
	if err != nil {
		return nil, err
	}

	current.Token = token.Token
	current.Modified = time.Now().Local()

	updated, err := domain.TokenRepo.Reset(current)

	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (as authService) Delete(i int64) error_utils.MessageErr {
	tk, err := domain.TokenRepo.Get(i)
	if err != nil {
		return err
	}
	err = domain.TokenRepo.Delete(tk.Id)
	if err != nil {
		return err
	}

	return nil
}
