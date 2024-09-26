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
	"crypto/x509"
	"encoding/pem"
	"net/url"
	"testing"

	"github.com/shurcooL/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTPGPublicKey(t *testing.T) {
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}

	expectedPublicKey := "mocked-public-key"
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		query := args.Get(1).(*teamQuery)
		query.Team.TPGPublicKey = graphql.String(expectedPublicKey)
	}).Return(nil)

	publicKey, err := client.GetTPGPublicKey()

	assert.NoError(t, err)
	assert.Equal(t, expectedPublicKey, publicKey)
	mockClient.AssertExpectations(t)
}

func TestEncryptPII(t *testing.T) {
	// Generate a temporary RSA key pair for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	// Extract the public key and encode it in PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Create the mock client
	mockClient := new(MockGraphQLClient)
	client := &Client{
		GraphQLClient: mockClient,
	}
	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		query := args.Get(1).(*teamQuery)
		query.Team.TPGPublicKey = graphql.String(publicKeyPEM)
	}).Return(nil)

	// Create the PII data to be encrypted
	data := url.Values{}
	data.Set("email", "test@example.com")
	data.Set("first_name", "John")
	data.Set("last_name", "Doe")

	// Call the EncryptPII function
	encryptedPII, err := client.EncryptPII(data.Encode())
	assert.NoError(t, err)

	// Ensure the encrypted PII is a non-empty string
	assert.NotEmpty(t, encryptedPII)
}
