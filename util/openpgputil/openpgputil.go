// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package openpgputil

import (
	"compress/flate"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

func Config() *packet.Config {
	return &packet.Config{
		DefaultCompressionAlgo: packet.CompressionZLIB,
		CompressionConfig: &packet.CompressionConfig{
			Level: flate.DefaultCompression,
		},
	}
}

func KeyGenConfig(rsaBits int) *packet.Config {
	cfg := Config()
	cfg.RSABits = rsaBits
	return cfg
}

func ReadArmoredKeyRing(armored string) (openpgp.EntityList, error) {
	entities, err := openpgp.ReadArmoredKeyRing(strings.NewReader(armored))
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, errors.New("no keys found")
	}
	return entities, nil
}

func ReadArmoredEntity(armored string) (*openpgp.Entity, error) {
	entities, err := ReadArmoredKeyRing(armored)
	if err != nil {
		return nil, err
	}
	if len(entities) != 1 {
		return nil, fmt.Errorf("expected one key, got %d", len(entities))
	}
	return entities[0], nil
}

func ReadPrivateKeyRing(privateKey, passphrase string) (openpgp.EntityList, error) {
	entities, err := ReadArmoredKeyRing(privateKey)
	if err != nil {
		return nil, err
	}
	for _, entity := range entities {
		if err := Unlock(entity, passphrase); err != nil {
			return nil, err
		}
	}
	return entities, nil
}

func Unlock(entity *openpgp.Entity, passphrase string) error {
	pass := []byte(passphrase)
	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		if err := entity.PrivateKey.Decrypt(pass); err != nil {
			return err
		}
	}
	for i := range entity.Subkeys {
		pk := entity.Subkeys[i].PrivateKey
		if pk != nil && pk.Encrypted {
			if err := pk.Decrypt(pass); err != nil {
				return err
			}
		}
	}
	return nil
}

func Encrypt(dst io.Writer, entities openpgp.EntityList) (io.WriteCloser, error) {
	return openpgp.Encrypt(dst, []*openpgp.Entity(entities), nil, nil, Config())
}

func EncryptArmored(publicKey string, dst io.Writer) (io.WriteCloser, error) {
	publicKeyRing, err := ReadArmoredKeyRing(publicKey)
	if err != nil {
		return nil, err
	}
	return Encrypt(dst, publicKeyRing)
}

func DecryptArmored(privateKey, passphrase string, src io.Reader) (io.Reader, error) {
	privateKeyRing, err := ReadPrivateKeyRing(privateKey, passphrase)
	if err != nil {
		return nil, err
	}
	md, err := openpgp.ReadMessage(src, privateKeyRing, nil, Config())
	if err != nil {
		return nil, err
	}
	return md.UnverifiedBody, nil
}

func GenerateArmoredKeyPair(name, email, passphrase string,
	rsaBits int) (publicKey, privateKey string, err error) {
	entity, err := openpgp.NewEntity(name, "", email, KeyGenConfig(rsaBits))
	if err != nil {
		return "", "", err
	}
	if passphrase != "" {
		if entity.PrivateKey != nil {
			if err := entity.PrivateKey.Encrypt([]byte(passphrase)); err != nil {
				return "", "", err
			}
		}
		for i := range entity.Subkeys {
			if pk := entity.Subkeys[i].PrivateKey; pk != nil {
				if err := pk.Encrypt([]byte(passphrase)); err != nil {
					return "", "", err
				}
			}
		}
	}
	publicKey, err = serializePublic(entity)
	if err != nil {
		return "", "", err
	}
	privateKey, err = serializePrivate(entity)
	if err != nil {
		return "", "", err
	}
	return publicKey, privateKey, nil
}

func KeyIDHex(entity *openpgp.Entity) string {
	if entity == nil || entity.PrimaryKey == nil {
		return ""
	}
	return fmt.Sprintf("%016X", entity.PrimaryKey.KeyId)
}

func FirstIdentity(entity *openpgp.Entity) string {
	if entity == nil || len(entity.Identities) == 0 {
		return ""
	}
	names := make([]string, 0, len(entity.Identities))
	for name := range entity.Identities {
		names = append(names, name)
	}
	sort.Strings(names)
	return names[0]
}

func serializePublic(entity *openpgp.Entity) (string, error) {
	var s strings.Builder
	w, err := armor.Encode(&s, openpgp.PublicKeyType, nil)
	if err != nil {
		return "", err
	}
	if err := entity.Serialize(w); err != nil {
		w.Close()
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return s.String(), nil
}

func serializePrivate(entity *openpgp.Entity) (string, error) {
	var s strings.Builder
	w, err := armor.Encode(&s, openpgp.PrivateKeyType, nil)
	if err != nil {
		return "", err
	}
	if err := entity.SerializePrivate(w, Config()); err != nil {
		w.Close()
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return s.String(), nil
}
