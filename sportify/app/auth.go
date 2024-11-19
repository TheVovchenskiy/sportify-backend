package app

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/hashing"

	"github.com/go-pkgz/auth/provider"
	"github.com/google/uuid"
)

type AuthStorage interface {
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	GetPasswordByUsername(ctx context.Context, username string) (string, error)
	CreateUser(ctx context.Context, id uuid.UUID, username, password string, tgUserID *int64) (models.ResponseSuccessLogin, error)
}

var _ AuthStorage = (*db.PostgresStorage)(nil)

func (a *App) NewCredCheckFunc(ctx context.Context) provider.CredCheckerFunc {
	return func(username, plainPassword string) (bool, error) {
		passHash, err := a.authStorage.GetPasswordByUsername(ctx, username)
		if err != nil {
			if errors.Is(err, db.ErrUserNotFound) {
				return false, nil
			}

			return false, fmt.Errorf("get password by username: %w", err)
		}

		return hashing.ComparePassAndHash(passHash, plainPassword)
	}
}

var (
	regexpUsername      = regexp.MustCompile(`^[\p{L}\d\s_]{1,128}$`)
	ErrNotValidUsername = errors.New(`1. Ваш username может содержать только буквы, цифры, пробелы и символ "_"
2. Его длина должна быть от 1 до 128 символов.
3. Он не может состоять только из пробелов.`)
	ErrNotUniqueUsername = errors.New("Такой username уже занят. Выберите другой")

	regexpPassword      = regexp.MustCompile(`^[\p{L}\d\s_]{8,256}$`)
	ErrNotValidPassword = errors.New(`1. Ваш пароль может содержать только буквы, цифры, пробелы и символ "_"
2. Его длина должна быть от 8 до 256 символов.
3. Он должен содержать хотя бы одну букву и одну цифру.
4. Он не может состоять только из пробелов.`)
)

func (a *App) ValidateUsernameAndPassword(username, password string) (string, string, error) {
	username = strings.TrimSpace(username)

	isUsernameValid := regexpUsername.Match([]byte(username))
	if !isUsernameValid {
		return "", "", ErrNotValidUsername
	}

	isPasswordValid := strings.ContainsFunc(password, unicode.IsNumber)
	if !isPasswordValid {
		return "", "", ErrNotValidPassword
	}

	isPasswordValid = strings.ContainsFunc(password, unicode.IsLetter)
	if !isPasswordValid {
		return "", "", ErrNotValidPassword
	}

	isPasswordValid = regexpPassword.Match([]byte(password))
	if !isPasswordValid {
		return "", "", ErrNotValidPassword
	}

	return username, password, nil
}

func (a *App) CreateUser(ctx context.Context, username, password string, tgUserID *int64) (models.ResponseSuccessLogin, error) {
	isUsernameExists, err := a.authStorage.CheckUsernameExists(ctx, username)
	if err != nil {
		return models.ResponseSuccessLogin{}, fmt.Errorf("to check username exists: %w", err)
	}
	if isUsernameExists {
		return models.ResponseSuccessLogin{}, ErrNotUniqueUsername
	}

	hashPass, err := hashing.HashPass(password)
	if err != nil {
		return models.ResponseSuccessLogin{}, fmt.Errorf("to hash pass: %w", err)
	}

	responseSuccessRegister, err := a.authStorage.CreateUser(ctx, uuid.New(), username, hashPass, tgUserID)
	if err != nil {
		return models.ResponseSuccessLogin{}, fmt.Errorf("to register: %w", err)
	}

	return responseSuccessRegister, nil
}
