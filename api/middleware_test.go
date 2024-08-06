package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Yelsnik/trackinginventory/token"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	ownerID int64,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(ownerID, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
