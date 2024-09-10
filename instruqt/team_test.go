package instruqt

import (
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
