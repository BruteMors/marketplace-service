package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	grpcmodels "github.com/BruteMors/marketplace-service/cart/pkg/api/grpc/loms/v1"
	"github.com/BruteMors/marketplace-service/test/e2e/responses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const baseURL = "http://localhost:8082"

func TestCartAndOrderFlow(t *testing.T) {
	userID := 1
	skuID := 1076963
	count := 3

	// Step 1: Add item to cart
	addUrl := baseURL + "/user/" + strconv.Itoa(userID) + "/cart/" + strconv.Itoa(skuID)
	addRequestBody, _ := json.Marshal(map[string]int{"count": count})

	addRequest, err := http.NewRequest("POST", addUrl, bytes.NewBuffer(addRequestBody))
	require.NoError(t, err, "Creating POST request should not fail")

	addResponse, err := http.DefaultClient.Do(addRequest)
	require.NoError(t, err, "POST request should succeed")
	defer addResponse.Body.Close()

	assert.Equal(t, http.StatusOK, addResponse.StatusCode, "Response status code should be 200 OK after adding item")

	// Step 2: Get cart items
	getUrl := baseURL + "/user/" + strconv.Itoa(userID) + "/cart/list"
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
	assert.Equal(t, uint16(count), cartResponse.Items[0].Count, "Item count should match")
	assert.Equal(t, int64(skuID), cartResponse.Items[0].SkuID, "SkuID should match")
	assert.Greater(t, cartResponse.TotalPrice, uint32(0), "Total price should be greater than 0")

	// Step 3: Checkout cart
	checkoutUrl := baseURL + "/user/" + strconv.Itoa(userID) + "/cart/checkout"
	checkoutRequest, err := http.NewRequest("POST", checkoutUrl, nil)
	require.NoError(t, err, "Creating POST request should not fail")

	checkoutResponse, err := http.DefaultClient.Do(checkoutRequest)
	require.NoError(t, err, "POST request should succeed")
	defer checkoutResponse.Body.Close()

	assert.Equal(t, http.StatusOK, checkoutResponse.StatusCode, "Response status code should be 200 OK after checkout")
	var checkoutResponseData responses.Checkout
	err = json.NewDecoder(checkoutResponse.Body).Decode(&checkoutResponseData)
	require.NoError(t, err, "Decoding response JSON should not fail")
	orderID := checkoutResponseData.OrderID

	// Step 4: Pay for the order
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "gRPC connection should not fail")
	defer conn.Close()
	client := grpcmodels.NewOrdersClient(conn)

	payRequest := &grpcmodels.OrderPayRequest{OrderId: orderID}
	_, err = client.OrderPay(context.Background(), payRequest)
	require.NoError(t, err, "gRPC OrderPay request should not fail")
}
