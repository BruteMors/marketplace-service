package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteItem(t *testing.T) {
	userID := 1
	skuID := 1076963
	url := "http://localhost:8082" + "/user/" + strconv.Itoa(userID) + "/cart/" + strconv.Itoa(skuID)
	requestBody, _ := json.Marshal(map[string]int{
		"count": 3,
	})

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	require.NoError(t, err, "Creating POST request should not fail")

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err, "POST request should succeed")
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode, "Response status code should be 200 OK")

	request, err = http.NewRequest("DELETE", url, nil)
	require.NoError(t, err, "Creating DELETE request should not fail")

	response, err = http.DefaultClient.Do(request)
	require.NoError(t, err, "DELETE request should succeed")
	defer response.Body.Close()

	assert.Equal(t, http.StatusNoContent, response.StatusCode, "Response status code should be 204 No Content")
}
