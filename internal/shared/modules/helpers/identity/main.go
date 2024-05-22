package identity

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

const fileName = "identity.bin"

func GetOrCreateIdentity(defaultIdentity Identity) (Identity, error) {
	identity, err := loadIdentity()
	if err != nil {
		return Identity{}, err
	}

	if identity == (Identity{}) {
		err := SaveIdentity(defaultIdentity)
		if err != nil {
			return Identity{}, err
		}

		return defaultIdentity, nil
	}

	return identity, nil
}

func SaveIdentity(identity Identity) error {
	filePath, err := getIdentityFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(identity, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	encrypt, err := crypto.Encrypt(data, nil)
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("%w: %s", errors.ErrIdentityDirectoryCreateFailure, err.Error())
	}

	err = os.WriteFile(filePath, encrypt, 0644)
	if err != nil {
		return fmt.Errorf("%w: %s", errors.ErrIdentityWriteFailure, err.Error())
	}

	return nil
}

func loadIdentity() (Identity, error) {
	var identity Identity

	filePath, err := getIdentityFilePath()
	if err != nil {
		return identity, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return identity, nil
		}

		return identity, fmt.Errorf("%w: %s", errors.ErrIdentityReadFailure, err.Error())
	}

	decrypt, err := crypto.Decrypt(data, nil)
	if err != nil {
		return identity, err
	}

	err = json.Unmarshal(decrypt, &identity)
	if err != nil {
		return identity, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	return identity, nil
}

func getIdentityFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	return filepath.Join(homeDir, strings.ReplaceAll(strings.ToLower(internal.GetConfig().Title), " ", "-"), fileName), nil
}