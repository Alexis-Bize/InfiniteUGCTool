package credentials

import (
	"encoding/json"
	"fmt"
	"infinite-bookmarker/internal"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/modules/crypto"
	"os"
	"path/filepath"
	"strings"
)

const fileName = "credentials.bin"

func GetOrCreateCredentials(defaultCredentials Credentials) (Credentials, error) {
	credentials, err := loadCredentials()
	if err != nil {
		return Credentials{}, err
	}

	if credentials == (Credentials{}) {
		err := saveCredentials(defaultCredentials)
		if err != nil {
			return Credentials{}, err
		}

		return defaultCredentials, nil
	}

	return credentials, nil
}

func loadCredentials() (Credentials, error) {
	var credentials Credentials

	filePath, err := getCredentialsFilePath()
	if err != nil {
		return credentials, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return credentials, nil
		}

		return credentials, fmt.Errorf("%w: %s", errors.ErrCredentialsReadFailure, err.Error())
	}

	decrypt, err := crypto.Decrypt(data, nil)
	if err != nil {
		return credentials, err
	}

	err = json.Unmarshal(decrypt, &credentials)
	if err != nil {
		return credentials, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	return credentials, nil
}

func saveCredentials(credentials Credentials) error {
	filePath, err := getCredentialsFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	encrypt, err := crypto.Encrypt(data, nil)
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("%w: %s", errors.ErrCredentialsDirectoryCreateFailure, err.Error())
	}

	err = os.WriteFile(filePath, encrypt, 0644)
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrCredentialsWriteFailure, err.Error())
	}

	return nil
}

func getCredentialsFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	return filepath.Join(homeDir, strings.ReplaceAll(strings.ToLower(internal.GetConfig().Title), " ", "-"), fileName), nil
}