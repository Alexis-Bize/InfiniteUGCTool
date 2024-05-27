// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
	"runtime"

	"infinite-ugc-tool/pkg/modules/errors"
)

const keyLength = 32

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) == 0 {
		key = getLocalKey()
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) == 0 {
		key = getLocalKey()
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.Format("ciphertext too short", errors.ErrInternal)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// security 1:1, more for user's convenience
func getLocalKey() []byte {
	hostname, err := os.Hostname()
	if err != nil {
		return fixKeyLength([]byte(runtime.GOOS))
	}

	return fixKeyLength([]byte(hostname))
}

func fixKeyLength(key []byte) []byte {
	if len(key) >= keyLength {
		return key[:keyLength]
	}

	paddedKey := make([]byte, keyLength)
	copy(paddedKey, key)
	return paddedKey
}
