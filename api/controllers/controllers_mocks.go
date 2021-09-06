package controllers

import (
	"github.com/RemeJuan/lattr/domain"
	"github.com/RemeJuan/lattr/utils/error_utils"
)

var (
	createTweetService     func(message *domain.Tweet) (*domain.Tweet, error_utils.MessageErr)
	getTweetService        func(msgId int64) (*domain.Tweet, error_utils.MessageErr)
	updateTweetService     func(tweet *domain.Tweet) (*domain.Tweet, error_utils.MessageErr)
	deleteTweetService     func(msgId int64) error_utils.MessageErr
	getAllTweetService     func(userId string) ([]domain.Tweet, error_utils.MessageErr)
	getPendingTweetService func() ([]domain.Tweet, error_utils.MessageErr)
	getLastTweet           func() (*domain.Tweet, error_utils.MessageErr)
	createTokenService     func(token *domain.Token) (*domain.Token, error_utils.MessageErr)
	getTokenService        func(id int64) (*domain.Token, error_utils.MessageErr)
	listTokensService      func() ([]domain.Token, error_utils.MessageErr)
	resetTokensService     func(token *domain.Token) (*domain.Token, error_utils.MessageErr)
	deleteTokensService    func(id int64) error_utils.MessageErr
)

type tweetServiceMock struct {
}

func (sm *tweetServiceMock) Create(tweet *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
	return createTweetService(tweet)
}

func (sm *tweetServiceMock) Get(id int64) (*domain.Tweet, error_utils.MessageErr) {
	return getTweetService(id)
}

func (sm *tweetServiceMock) GetAll(userId string) ([]domain.Tweet, error_utils.MessageErr) {
	return getAllTweetService(userId)
}

func (sm *tweetServiceMock) Update(tweet *domain.Tweet) (*domain.Tweet, error_utils.MessageErr) {
	return updateTweetService(tweet)
}

func (sm *tweetServiceMock) Delete(id int64) error_utils.MessageErr {
	return deleteTweetService(id)
}

func (sm *tweetServiceMock) GetPending() ([]domain.Tweet, error_utils.MessageErr) {
	return getPendingTweetService()
}

func (sm *tweetServiceMock) GetLast() (*domain.Tweet, error_utils.MessageErr) {
	return getLastTweet()
}

type authServiceMock struct{}

func (asm *authServiceMock) Create(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
	return createTokenService(token)
}

func (asm *authServiceMock) Get(id int64) (*domain.Token, error_utils.MessageErr) {
	return getTokenService(id)
}

func (asm *authServiceMock) List() ([]domain.Token, error_utils.MessageErr) {
	return listTokensService()
}

func (asm *authServiceMock) Reset(token *domain.Token) (*domain.Token, error_utils.MessageErr) {
	return resetTokensService(token)
}

func (asm *authServiceMock) Delete(id int64) error_utils.MessageErr {
	return deleteTokensService(id)
}

func (asm *authServiceMock) ValidateToken(token *domain.Token, requiredScope string) bool {
	return true
}
