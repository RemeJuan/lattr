package services

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/error_utils"
	"github.com/stretchr/testify/assert"
)

var (
	createTokenDomain func(token *domain.Token) (*domain.Token, error_utils.MessageErr)
	getTokenDomain    func(tokenId int64) (*domain.Token, error_utils.MessageErr)
	getTokensDomain   func() ([]domain.Token, error_utils.MessageErr)
	resetTokenDomain  func(*domain.Token) (*domain.Token, error_utils.MessageErr)
	deleteTokenDomain func(tokenId int64) error_utils.MessageErr
)

type authDBMock struct {
	domain.TokenRepoInterface
}

func (m *authDBMock) Create(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
	return createTokenDomain(token)
}

func (m *authDBMock) Get(id int64) (*domain.Token, error_utils.MessageErr) {
	return getTokenDomain(id)
}

func (m *authDBMock) List() ([]domain.Token, error_utils.MessageErr) {
	return getTokensDomain()
}

func (m *authDBMock) Reset(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
	return resetTokenDomain(token)
}

func (m *authDBMock) Delete(id int64) error_utils.MessageErr {
	return deleteTokenDomain(id)
}

const mockTokenId int64 = 1
const mockTokenName = "token name"
const mockToken = "mock-token"

func TestAuthService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		createTokenDomain = func(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
			return &domain.Token{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				Scopes:    []string{"token:create"},
				CreatedAt: tm,
				Modified:  tm,
			}, nil
		}

		request := &domain.Token{
			Name:      mockTokenName,
			Token:     mockToken,
			Scopes:    []string{"token:create"},
			CreatedAt: tm,
			Modified:  tm,
		}

		result, err := AuthService.Create(request)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockTokenId, result.Id)
		assert.Equal(t, mockToken, result.Token)
		assert.Equal(t, mockTokenName, result.Name)
		assert.Equal(t, tm, result.CreatedAt)
		assert.Equal(t, tm, result.Modified)
	})

	t.Run("Validation failed", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		createTokenDomain = func(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
			return &domain.Token{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				CreatedAt: tm,
				Modified:  tm,
			}, nil
		}

		request := &domain.Token{
			Name:      "",
			Token:     mockToken,
			CreatedAt: tm,
			Modified:  tm,
		}

		result, err := AuthService.Create(request)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusUnprocessableEntity, err.Status())
		assert.EqualValues(t, "invalid_request", err.Error())
		assert.EqualValues(t, "Name cannot be empty", err.Message())
	})

	t.Run("Create failed", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		createTokenDomain = func(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
			return nil, error_utils.InternalServerError("Unknown error occurred")
		}

		request := &domain.Token{
			Name:      mockTokenName,
			Token:     mockToken,
			Scopes:    []string{"token:create"},
			CreatedAt: tm,
			Modified:  tm,
		}

		result, err := AuthService.Create(request)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
		assert.EqualValues(t, "server_error", err.Error())
		assert.EqualValues(t, "Unknown error occurred", err.Message())
	})
}

func TestAuthService_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		getTokenDomain = func(id int64) (*domain.Token, error_utils.MessageErr) {
			return &domain.Token{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				CreatedAt: tm,
				Modified:  tm,
			}, nil
		}

		result, err := AuthService.Get(mockTokenId)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockTokenId, result.Id)
		assert.Equal(t, mockTokenName, result.Name)
		assert.Equal(t, mockToken, result.Token)
		assert.Equal(t, tm, result.CreatedAt)
		assert.Equal(t, tm, result.Modified)
	})

	t.Run("Not found", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}
		const errMessage = "No token found for ID"

		getTokenDomain = func(id int64) (*domain.Token, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError(errMessage)
		}

		result, err := AuthService.Get(mockTokenId)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Status())
		assert.Equal(t, "not_found", err.Error())
		assert.Equal(t, errMessage, err.Message())
	})
}

func TestAuthService_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		getTokensDomain = func() ([]domain.Token, error_utils.MessageErr) {
			return []domain.Token{
				{
					Id:        mockTokenId,
					Name:      mockTokenName,
					Token:     mockToken,
					CreatedAt: tm,
					Modified:  tm,
				},
				{
					Id:        mockTokenId,
					Name:      mockTokenName,
					Token:     mockToken,
					CreatedAt: tm,
					Modified:  tm,
				},
			}, nil
		}

		result, err := AuthService.List()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		first := result[0]
		second := result[1]

		assert.Equal(t, mockTokenId, first.Id)
		assert.Equal(t, mockTokenName, first.Name)
		assert.Equal(t, mockToken, first.Token)
		assert.Equal(t, tm, first.CreatedAt)
		assert.Equal(t, tm, first.Modified)

		assert.Equal(t, mockTokenId, second.Id)
		assert.Equal(t, mockTokenName, second.Name)
		assert.Equal(t, mockToken, second.Token)
		assert.Equal(t, tm, second.CreatedAt)
		assert.Equal(t, tm, second.Modified)
	})

	t.Run("Not found", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		const errMessage = "No records found"

		getTokensDomain = func() ([]domain.Token, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError(errMessage)
		}

		result, err := AuthService.List()

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Status())
		assert.Equal(t, "not_found", err.Error())
		assert.Equal(t, errMessage, err.Message())
	})
}

func TestAuthService_Reset(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		getTokenDomain = func(id int64) (*domain.Token, error_utils.MessageErr) {
			return &domain.Token{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				CreatedAt: tm,
				Modified:  tm,
			}, nil
		}

		resetTokenDomain = func(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
			return &domain.Token{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				CreatedAt: tm,
				Modified:  tm,
			}, nil
		}

		request := &domain.Token{
			Id:       mockTokenId,
			Token:    mockToken,
			Modified: tm,
		}

		result, err := AuthService.Reset(request)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockTokenId, result.Id)
		assert.Equal(t, mockToken, result.Token)
		assert.Equal(t, tm, result.Modified)
	})

	t.Run("Not found", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		const errMessage = "No matching record found"
		getTokenDomain = func(id int64) (*domain.Token, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError(errMessage)
		}

		request := &domain.Token{
			Id:       mockTokenId,
			Token:    mockToken,
			Modified: tm,
		}

		result, err := AuthService.Reset(request)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "not_found", err.Error())
		assert.EqualValues(t, errMessage, err.Message())
	})

	t.Run("Reset failed", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		const errMessage = "Unable to reset token"

		getTokenDomain = func(id int64) (*domain.Token, error_utils.MessageErr) {
			return &domain.Token{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				CreatedAt: tm,
				Modified:  tm,
			}, nil
		}

		resetTokenDomain = func(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
			return nil, error_utils.InternalServerError(errMessage)
		}

		request := &domain.Token{
			Id:       mockTokenId,
			Token:    mockToken,
			Modified: tm,
		}

		result, err := AuthService.Reset(request)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
		assert.EqualValues(t, "server_error", err.Error())
		assert.EqualValues(t, errMessage, err.Message())
	})
}

func TestAuthService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		getTokenDomain = func(id int64) (*domain.Token, error_utils.MessageErr) {
			return &domain.Token{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				CreatedAt: tm,
				Modified:  tm,
			}, nil
		}

		deleteTokenDomain = func(id int64) error_utils.MessageErr {
			return nil
		}

		err := AuthService.Delete(mockTokenId)

		assert.Nil(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		const errMessage = "No matching record found"
		getTokenDomain = func(id int64) (*domain.Token, error_utils.MessageErr) {
			return nil, error_utils.NotFoundError(errMessage)
		}

		err := AuthService.Delete(mockTokenId)

		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusNotFound, err.Status())
		assert.EqualValues(t, "not_found", err.Error())
		assert.EqualValues(t, errMessage, err.Message())
	})

	t.Run("Delete failed", func(t *testing.T) {
		domain.TokenRepo = &authDBMock{}

		const errMessage = "Unable to reset token"

		getTokenDomain = func(id int64) (*domain.Token, error_utils.MessageErr) {
			return &domain.Token{
				Id:        mockTokenId,
				Name:      mockTokenName,
				Token:     mockToken,
				CreatedAt: tm,
				Modified:  tm,
			}, nil
		}

		deleteTokenDomain = func(id int64) error_utils.MessageErr {
			return error_utils.InternalServerError(errMessage)
		}

		err := AuthService.Delete(mockTokenId)

		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusInternalServerError, err.Status())
		assert.EqualValues(t, "server_error", err.Error())
		assert.EqualValues(t, errMessage, err.Message())
	})
}

func TestValidateToken(t *testing.T) {
	currentTime = time.Date(2021, 8, 20, 13, 10, 0, 0, time.Local)

	activeTokens = []domain.Token{
		{
			Token:     "1",
			Scopes:    []string{"token:read", "token:create"},
			ExpiresAt: currentTime.Add(time.Hour * time.Duration(1)),
		},
		{
			Token:     "2",
			Scopes:    []string{"token:read", "token:create"},
			ExpiresAt: currentTime.Add(time.Hour * time.Duration(1)),
		},
	}

	t.Run("Exists", func(t *testing.T) {
		token := domain.Token{Token: "1"}

		result := AuthService.ValidateToken(&token, "token:read")

		assert.Equal(t, true, result)
	})

	t.Run("!Exists", func(t *testing.T) {
		token := domain.Token{Token: "3"}

		result := AuthService.ValidateToken(&token, "token:read")

		assert.Equal(t, false, result)
	})
}

func TestTokenExpiryDate(t *testing.T) {
	currentTime = time.Date(2021, 8, 20, 13, 10, 0, 0, time.Local)

	t.Run("Default", func(t *testing.T) {
		_ = os.Setenv("TOKEN_VALIDITY_HOURS", "1")

		expiry := TokenExpiryDate(0)

		result := time.Date(2021, 8, 20, 14, 10, 0, 0, time.Local)

		assert.NotNil(t, expiry)
		assert.Equal(t, result, expiry)
	})

	t.Run("Specified validity", func(t *testing.T) {
		_ = os.Setenv("TOKEN_VALIDITY_HOURS", "1")

		expiry := TokenExpiryDate(2)

		result := time.Date(2021, 8, 20, 15, 10, 0, 0, time.Local)

		assert.NotNil(t, expiry)
		assert.Equal(t, result, expiry)
	})

	t.Run("Invalid env", func(t *testing.T) {
		_ = os.Setenv("TOKEN_VALIDITY_HOURS", "")

		expiry := TokenExpiryDate(0)

		result := time.Date(2021, 8, 20, 14, 10, 0, 0, time.Local)

		assert.NotNil(t, expiry)
		assert.Equal(t, result, expiry)
	})
}
