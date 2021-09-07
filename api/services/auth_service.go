package services

import (
	"os"
	"strconv"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/google/uuid"
)

var (
	AuthService  authServiceInterface = &authService{}
	activeTokens []domain.Token
	currentTime  = time.Now().Local()
)

type authService struct{}

type authServiceInterface interface {
	Create(*domain.Token) (*domain.Token, error_utils.MessageErr)
	Get(int64) (*domain.Token, error_utils.MessageErr)
	List() ([]domain.Token, error_utils.MessageErr)
	Reset(*domain.Token) (*domain.Token, error_utils.MessageErr)
	Delete(int64) error_utils.MessageErr
	ValidateToken(token *domain.Token, requiredScope string) bool
}

func (as authService) Create(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
	if err := token.Validate(); err != nil {
		return nil, err
	}

	generated := uuid.New().String()
	token.Token = generated
	token.CreatedAt = currentTime
	token.Modified = currentTime
	token.ExpiresAt = TokenExpiryDate(token.Validity)

	tk, err := domain.TokenRepo.Create(token)
	if err != nil {
		return nil, err
	}

	updateInMemoryTokens(token, nil)
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
	tokens, err := domain.TokenRepo.List()
	if err != nil {
		return nil, err
	}

	activeTokens = tokens
	return tokens, nil
}

func (as authService) Reset(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
	current, err := domain.TokenRepo.Get(token.Id)
	if err != nil {
		return nil, err
	}

	current.Token = token.Token
	current.Modified = time.Now().Local()
	token.ExpiresAt = TokenExpiryDate(token.Validity)

	updated, err := domain.TokenRepo.Reset(current)

	updateInMemoryTokens(updated, token)

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

	updateInMemoryTokens(nil, tk)

	if err != nil {
		return err
	}

	return nil
}

func (as authService) ValidateToken(token *domain.Token, requiredScope string) bool {
	for _, val := range activeTokens {
		if val.Token == token.Token {
			return containsRequiredScope(&val, requiredScope) && val.ExpiresAt.After(currentTime)
		}
	}

	return false
}

func TokenExpiryDate(validity int) time.Time {
	tvh := os.Getenv("TOKEN_VALIDITY_HOURS")
	dur := validity

	if validity == 0 {
		t, err := strconv.ParseInt(tvh, 10, 0)

		if err != nil {
			dur = 1
		} else {
			dur = int(t)
		}
	}

	return currentTime.Add(time.Hour * time.Duration(dur)).Local()
}

func updateInMemoryTokens(adding *domain.Token, removing *domain.Token) {
	if removing != nil {
		removeExpiredToken(*removing)
	}

	if adding != nil {
		activeTokens = append(activeTokens, *adding)
	}
}

func removeExpiredToken(tk domain.Token) {
	a := make([]domain.Token, 0)

	for _, val := range activeTokens {
		if val.Token != tk.Token {
			a = append(a, val)
		}
	}

	activeTokens = a
}

func containsRequiredScope(token *domain.Token, requiredScope string) bool {
	for _, val := range token.Scopes {
		if val == requiredScope {
			return true
		}
	}

	return false
}
