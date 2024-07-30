package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/BruteMors/marketplace-service/test/e2e/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCart(t *testing.T) {
	userID := 1
	addUrl := "http://localhost:8082/user/" + strconv.Itoa(userID) + "/cart/1076963"
	addRequestBody, _ := json.Marshal(map[string]int{"count": 3})

	addRequest, err := http.NewRequest("POST", addUrl, bytes.NewBuffer(addRequestBody))
	require.NoError(t, err, "Creating POST request should not fail")

	addResponse, err := http.DefaultClient.Do(addRequest)
	require.NoError(t, err, "POST request should succeed")
	defer addResponse.Body.Close()

	assert.Equal(t, http.StatusOK, addResponse.StatusCode, "Response status code should be 200 OK after adding item")

	getUrl := "http://localhost:8082/user/" + strconv.Itoa(userID) + "/cart/list"
	getRequest, err := http.NewRequest("GET", getUrl, nil)
	require.NoError(t, err, "Creating GET request should not fail")

	getResponse, err := http.DefaultClient.Do(getRequest)
	require.NoError(t, err, "GET request should succeed")
	defer getResponse.Body.Close()

	assert.Equal(t, http.StatusOK, getResponse.StatusCode, "Response status code should be 200 OK")
	var cartResponse responses.GetCart
	err = json.NewDecoder(getResponse.Body).Decode(&cartResponse)
	require.NoError(t, err, "Decoding response JSON should not fail")

	require.NotEmpty(t, cartResponse.Items, "Cart should not be empty")
	assert.Equal(t, uint16(3), cartResponse.Items[0].Count, "Item count should match")
	assert.Equal(t, int64(1076963), cartResponse.Items[0].SkuID, "SkuID should match")
	assert.Greater(t, cartResponse.TotalPrice, uint32(0), "Total price should be greater than 0")
}
