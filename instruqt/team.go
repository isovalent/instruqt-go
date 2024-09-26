// Copyright 2024 Cisco Systems, Inc. and its affiliates

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package instruqt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"

	"github.com/shurcooL/graphql"
)

// teamQuery represents the GraphQL query structure for retrieving the TPG public key
// associated with a specific team identified by its slug.
type teamQuery struct {
	Team struct {
		TPGPublicKey graphql.String `graphql:"tpgPublicKey"` // The TPG public key of the team.
	} `graphql:"team(teamSlug: $teamSlug)"`
}

// GetTPGPublicKey retrieves the TPG public key for the team associated with the client.
//
// Returns:
//   - string: The TPG public key of the team.
//   - error: Any error encountered while retrieving the TPG public key.
func (c *Client) GetTPGPublicKey() (string, error) {
	var q teamQuery
	variables := map[string]interface{}{
		"teamSlug": graphql.String(c.TeamSlug),
	}

	if err := c.GraphQLClient.Query(c.Context, &q, variables); err != nil {
		return "", fmt.Errorf("failed to retrieve TPG Public Key: %v", err)
	}

	return string(q.Team.TPGPublicKey), nil
}

// EncryptPII encrypts PII using the public key fetched from the GetTPGPublicKey function.
// It takes a string representing the PII data, encodes it, and then encrypts it using RSA.
func (c *Client) EncryptPII(encodedPII string) (string, error) {
	// Fetch the public key using the GetTPGPublicKey function
	publicKeyPEM, err := c.GetTPGPublicKey()
	if err != nil {
		return "", fmt.Errorf("failed to get public key: %v", err)
	}

	// Decode the PEM public key
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Parse the public key
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse DER encoded public key: %v", err)
	}

	// Assert the public key is of type *rsa.PublicKey
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("not an RSA public key")
	}

	// Encrypt the PII
	hash := sha256.New()
	encryptedPII, err := rsa.EncryptOAEP(hash, rand.Reader, rsaPublicKey, []byte(encodedPII), nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt PII: %v", err)
	}

	// Encode the encrypted data to base64
	encryptedPIIBase64 := base64.StdEncoding.EncodeToString(encryptedPII)
	return encryptedPIIBase64, nil
}

// EncryptUserPII creates PII data (first name, last name, and email) and encrypts it using the public key.
func (c *Client) EncryptUserPII(firstName, lastName, email string) (string, error) {
	// Prepare the PII data
	piiData := url.Values{
		"fn": {firstName},
		"ln": {lastName},
		"e":  {email},
	}

	// Encrypt the PII data
	encryptedPII, err := c.EncryptPII(piiData.Encode())
	if err != nil {
		return "", err
	}

	return encryptedPII, nil
}