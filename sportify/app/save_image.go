package app

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // Add jpeg format for image
	_ "image/png"  // Add png format for image
)

func hashContent(content []byte) (string, error) {
	hash := sha256.New()

	_, err := hash.Write(content)
	if err != nil {
		return "", fmt.Errorf("to write hash: %w", err)
	}

	result := hash.Sum(nil)

	return hex.EncodeToString(result), nil
}

var ErrWrongFormat = errors.New("формат файла должен быть png или jpeg")

func (a *App) SaveImage(ctx context.Context, file []byte) (string, error) {

	_, format, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return "", fmt.Errorf("to decode(format %s): %w", format, err)
		//return "", fmt.Errorf("вы используете формат %s %w", format, ErrWrongFormat)
	}

	fileName, err := hashContent(file)
	if err != nil {
		return "", fmt.Errorf("to hash content: %w", err)
	}

	existenceFile, err := a.fileStorage.Check(ctx, []string{fileName})
	if existenceFile[0] {
		return a.urlPrefixFile + fileName, nil
	}

	err = a.fileStorage.SaveFile(ctx, file, fileName)
	if err != nil {
		return "", fmt.Errorf("to save file: %w", err)
	}

	return a.urlPrefixFile + fileName, nil
}
