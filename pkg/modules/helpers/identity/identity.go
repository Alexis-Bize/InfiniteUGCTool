package identity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"infinite-ugc-tool/configs"
	"infinite-ugc-tool/pkg/modules/crypto"
	"infinite-ugc-tool/pkg/modules/errors"
)

type UserCredentials struct {
	Email		string `json:"email,omitempty"`
	Password	string `json:"password,omitempty"`
}

type SpartanTokenDetails struct {
	Value		string `json:"value,omitempty"`
	Expiration	string `json:"expiration,omitempty"`
}

type XboxNetworkIdentity struct {
	Xuid		string `json:"xuid,omitempty"`
	Gamertag	string `json:"gamertag,omitempty"`
}

type Identity struct {
	User			UserCredentials		`json:"user,omitempty"`
	SpartanToken	SpartanTokenDetails `json:"spartan_token,omitempty"`
	XboxNetwork		XboxNetworkIdentity	`json:"xbox_network,omitempty"`
}

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
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	encrypt, err := crypto.Encrypt(data, nil)
	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return errors.Format(err.Error(), errors.ErrIdentityDirectoryCreateFailure)
	}

	err = os.WriteFile(filePath, encrypt, 0644)
	if err != nil {
		return errors.Format(err.Error(), errors.ErrIdentityWriteFailure)
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

		return identity, errors.Format(err.Error(), errors.ErrIdentityReadFailure)
	}


	decrypt, err := crypto.Decrypt(data, nil)
	if err != nil {
		return identity, err
	}

	err = json.Unmarshal(decrypt, &identity)
	if err != nil {
		return identity, errors.Format(err.Error(), errors.ErrInternal)
	}

	return identity, nil
}

func getIdentityFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}

	return filepath.Join(homeDir, strings.ReplaceAll(configs.GetConfig().Name, " ", "-"), fileName), nil
}
