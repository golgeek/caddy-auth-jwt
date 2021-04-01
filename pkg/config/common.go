// Copyright 2020 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"crypto/ecdsa"
	"crypto/rsa"
	jwterrors "github.com/greenpau/caddy-auth-jwt/pkg/errors"
)

var (
	defaultKeyID             = "0"
	defaultTokenName         = "access_token"
	defaultTokenLifetime int = 900
)

// EnvTokenRSADir the env variable used to indicate a directory
const EnvTokenRSADir = "JWT_RSA_DIR"

// EnvTokenRSAFile then env variable (or prefix) used to indicate a file containing a RS key
const EnvTokenRSAFile = "JWT_RSA_FILE"

// EnvTokenRSAKey the env variable (or prefix) used to indicte a RS key
const EnvTokenRSAKey = "JWT_RSA_KEY"

// EnvTokenECDSADir the env variable used to indicate a directory containing ECDSA keys.
const EnvTokenECDSADir = "JWT_ECDSA_DIR"

// EnvTokenECDSAFile then env variable (or prefix) used to indicate a file containing ECDSA key.
const EnvTokenECDSAFile = "JWT_ECDSA_FILE"

// EnvTokenECDSAKey the env variable (or prefix) used to indicate ECDSA key.
const EnvTokenECDSAKey = "JWT_ECDSA_KEY"

// EnvTokenSecret the env variable used to indicate shared secret key.
const EnvTokenSecret = "JWT_TOKEN_SECRET"

// EnvTokenLifetime the env variable used to set default token lifetime.
const EnvTokenLifetime = "JWT_TOKEN_LIFETIME"

// EnvTokenName the env variable used to set default token name.
const EnvTokenName = "JWT_TOKEN_NAME"

// CommonTokenConfig is common token-related configuration settings.
// The setting are used by TokenProvider and TokenValidator.
type CommonTokenConfig struct {
	TokenSignMethod string `json:"token_sign_method,omitempty" xml:"token_sign_method,omitempty" yaml:"token_sign_method,omitempty"`
	TokenName       string `json:"token_name,omitempty" xml:"token_name" yaml:"token_name"`
	TokenOrigin     string `json:"token_origin,omitempty" xml:"token_origin" yaml:"token_origin"`
	// The expiration time of a token in seconds
	TokenLifetime int      `json:"token_lifetime,omitempty" xml:"token_lifetime" yaml:"token_lifetime"`
	EvalExpr      []string `json:"token_eval_expr,omitempty" xml:"token_eval_expr" yaml:"token_eval_expr"`

	HMACSignMethodConfig
	RSASignMethodConfig
	ECDSASignMethodConfig

	// The source of token configuration, config or environment variables.
	tokenOrigin string
	tokenType   string
	// The map containing key material, e.g. *rsa.PrivateKey, *rsa.PublicKey,
	// *ecdsa.PrivateKey, etc.
	tokenKeys map[string]interface{}
}

// HMACSignMethodConfig holds configuration for signing messages by means of a shared key.
type HMACSignMethodConfig struct {
	TokenSecret string `json:"token_secret,omitempty" xml:"token_secret" yaml:"token_secret"`
}

// RSASignMethodConfig holds data for RSA keys that can be used to sign and verify JWT tokens
// TokenRSDirectory is a directory that is like:
//
// <kid>'s can only contain ascii letters/numbers and underscores. (otherwise they are not loaded)
//
// <dirname>
//    +-- <kid_1>
//          +-- private.key
//    +-- <kid_2>
//          +-- public.key
//    +-- kid_3.key
//    +-- kid_4.key
//    +-- kid.5.key
// The above directory will result in a TokenRSKeys that looks like:
//
// TokenRSKeys{
//     "kid_1_private": "---- RSA PRIVATE KEY ---- ...",
//     "kid_2_public": "---- RSA PUBLIC KEY ---- ...",
//     "kid_3": "---- RSA PRIVATE KEY ---- ...",
//     "kid_4": "---- RSA PUBLIC KEY ---- ...",
//     // there is no "kid.5" becuase the "." is invalid.
// }
//
// There only needs to be public keys loaded for verification. If you're using the Grantor method then
// you need to load a PrivateKey so that keys can be signed.
//
// The TokenRS fields translate to the following config values:
//
// "token_rs_dir": "<path to dir>"
// "token_rs_files": {"<kid>": "<path to file>", ...}
// "token_rs_keys": {"<kid>": "<key PEM value>", ...}
//
// there are two special config values:
//
// "token_rs_file": "<path to file>"
// "token_rs_key": "<key PEM value>"
//
// The above two variables map to a <kid> of "0", these are always evaluated first so they can be overwritten if
// a <kid> of "0" is used explictly
//
// The TokenRS fields translate to the following enviornment variables:
//
// JWT_RS_DIR="<path to dir>"
// JWT_RS_FILE_<kid>="<path to file>"
// JWT_RS_KEY_<kid>="<key PEM value>"
//
// there are two special environment variables:
//
// JWT_RS_FILE="<path to file>"
// JWT_RS_KEY="<key PEM value>"
//
// The above two variables map to a <kid> of "0", these are always evaluated first so they can be overwritten if
// a <kid> of "0" is used explictly
//
// Enviroment variable KID's get lowercased. All other KID's are left untouched.

// RSASignMethodConfig defines configuration unique to RSA keys.
type RSASignMethodConfig struct {
	// TokenRSDir holds the absolute path to where a nested directory of key paths are, otherwise the name of the file
	// is used as the kid and the values are parse into TokenRSKeys
	TokenRSADir string `json:"token_rsa_dir,omitempty" xml:"token_rsa_dir" yaml:"token_rsa_dir"`

	// TokenRSFiles holds a map of <kid> to filename. These files should hold the public or private key. They are parsed to TokenRSKeys values
	TokenRSAFiles map[string]string `json:"token_rsa_files,omitempty" xml:"token_rsa_files" yaml:"token_rsa_files"`

	// TokenRSKeys holds a map of <kid> to the key PEM value
	TokenRSAKeys map[string]string `json:"token_rsa_keys,omitempty" xml:"token_rsa_keys" yaml:"token_rsa_keys"`

	// Special (see the comment above to see how they work)

	TokenRSAFile string `json:"token_rsa_file,omitempty" xml:"token_rsa_file" yaml:"token_rsa_file"`
	TokenRSAKey  string `json:"token_rsa_key,omitempty" xml:"token_rsa_key" yaml:"token_rsa_key"`
}

// ECDSASignMethodConfig defines configuration unique to ECDSA keys.
type ECDSASignMethodConfig struct {
	TokenECDSADir   string            `json:"token_ecdsa_dir,omitempty" xml:"token_ecdsa_dir" yaml:"token_ecdsa_dir"`
	TokenECDSAFiles map[string]string `json:"token_ecdsa_files,omitempty" xml:"token_ecdsa_files" yaml:"token_ecdsa_files"`
	TokenECDSAKeys  map[string]string `json:"token_ecdsa_keys,omitempty" xml:"token_ecdsa_keys" yaml:"token_ecdsa_keys"`
	TokenECDSAFile  string            `json:"token_ecdsa_file,omitempty" xml:"token_ecdsa_file" yaml:"token_ecdsa_file"`
	TokenECDSAKey   string            `json:"token_ecdsa_key,omitempty" xml:"token_ecdsa_key" yaml:"token_ecdsa_key"`
}

// HasRSAKeys returns true if the configuration has RSA encryption keys and files
func (c *CommonTokenConfig) HasRSAKeys() bool {
	if c.TokenRSADir != "" {
		return true
	}
	if c.TokenRSAFile != "" {
		return true
	}
	if c.TokenRSAKey != "" {
		return true
	}
	if c.TokenRSAFiles != nil {
		return true
	}
	if c.TokenRSAKeys != nil {
		return true
	}
	return false
}

// HasECDSAKeys returns true if the configuration has ECDSA encryption keys and files
func (c *CommonTokenConfig) HasECDSAKeys() bool {
	if c.TokenECDSADir != "" {
		return true
	}
	if c.TokenECDSAFile != "" {
		return true
	}
	if c.TokenECDSAKey != "" {
		return true
	}
	if c.TokenECDSAFiles != nil {
		return true
	}
	if c.TokenECDSAKeys != nil {
		return true
	}
	return false
}

// NewCommonTokenConfig returns an instance of CommonTokenConfig.
func NewCommonTokenConfig() *CommonTokenConfig {
	return &CommonTokenConfig{}
}

// GetOrigin returns the origin of the token, i.e. config or env.
func (c *CommonTokenConfig) GetOrigin() string {
	if c.tokenOrigin == "" {
		return "unknown"
	}
	return c.tokenOrigin
}

// SetOrigin sets token origin, i.e. config or env.
func (c *CommonTokenConfig) SetOrigin(name string) error {
	switch name {
	case "config", "env":
	case "empty":
		return jwterrors.ErrEmptyTokenConfigOrigin
	default:
		return jwterrors.ErrUnsupportedTokenConfigOrigin.WithArgs(name)
	}
	c.tokenOrigin = name
	return nil
}

// GetKeys returns a map with keys.
func (c *CommonTokenConfig) GetKeys() (string, map[string]interface{}) {
	return c.tokenType, c.tokenKeys
}

// AddPublicKey adds RSA public key to the map of RSA keys.
func (c *CommonTokenConfig) AddPublicKey(keyID string, keyMaterial interface{}) error {
	if keyID == "" {
		return jwterrors.ErrKeyIDNotFound
	}

	if c.tokenKeys == nil {
		c.tokenKeys = make(map[string]interface{})
	}

	switch kt := keyMaterial.(type) {
	case *rsa.PrivateKey:
		privkey := keyMaterial.(*rsa.PrivateKey)
		c.tokenKeys[keyID] = &privkey.PublicKey
		if _, exists := c.tokenKeys[defaultKeyID]; !exists {
			c.tokenKeys[defaultKeyID] = &privkey.PublicKey
		}
	case *ecdsa.PrivateKey:
		privkey := keyMaterial.(*ecdsa.PrivateKey)
		c.tokenKeys[keyID] = &privkey.PublicKey
		if _, exists := c.tokenKeys[defaultKeyID]; !exists {
			c.tokenKeys[defaultKeyID] = &privkey.PublicKey
		}
	case *rsa.PublicKey, *ecdsa.PublicKey:
		c.tokenKeys[keyID] = keyMaterial
	default:
		return jwterrors.ErrUnsupportedKeyType.WithArgs(kt, keyID)
	}
	return nil
}

// GetPrivateKey returns the first RSA private key it finds.
func (c *CommonTokenConfig) GetPrivateKey() (interface{}, string, error) {
	if c.tokenKeys == nil {
		return nil, "", jwterrors.ErrRSAKeysNotFound
	}
	for keyID, k := range c.tokenKeys {
		if keyID == defaultKeyID {
			continue
		}
		switch k.(type) {
		case *rsa.PrivateKey:
			return k, keyID, nil
		case *ecdsa.PrivateKey:
			return k, keyID, nil
		}
	}
	switch c.tokenType {
	case "ecdsa":
		return nil, "", jwterrors.ErrECDSAKeysNotFound
	}
	return nil, "", jwterrors.ErrRSAKeysNotFound
}

// AddKey adds token key.
func (c *CommonTokenConfig) AddKey(k string, pk interface{}) error {
	if c.tokenKeys == nil {
		c.tokenKeys = make(map[string]interface{})
	}
	keyType, err := c.getKeyType(pk)
	if err != nil {
		return err
	}
	if c.tokenType == "" {
		c.tokenType = keyType
	}
	if c.tokenType != keyType {
		return jwterrors.ErrMixedConfigKeyType.WithArgs(c.tokenType, keyType)
	}
	c.tokenKeys[k] = pk
	return nil
}

func (c *CommonTokenConfig) getKeyType(k interface{}) (string, error) {
	var kt string
	switch k.(type) {
	case string:
		kt = "secret"
	case *rsa.PrivateKey:
		kt = "rsa"
	case *rsa.PublicKey:
		kt = "rsa"
	case *ecdsa.PrivateKey:
		kt = "ecdsa"
	case *ecdsa.PublicKey:
		kt = "ecdsa"
	default:
		return "", jwterrors.ErrUnsupportedConfigKeyType.WithArgs(k)
	}
	return kt, nil
}

// LoadKeys loads key material.
func (c *CommonTokenConfig) LoadKeys() error {
	if len(c.tokenKeys) > 0 {
		// return jwterrors.ErrTokenAlreadyConfigured
		return nil
	}
	return c.load()
}
