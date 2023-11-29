package intents

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func submitHandler(c *gin.Context) {
	var body Body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Process the valid request
	c.JSON(http.StatusOK, gin.H{"message": "Received successfully"})
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	if err := NewValidator(); err != nil {
		panic(err)
	}

	r.POST("/submit", submitHandler)
	return r
}

func TestSubmitHandler(t *testing.T) {
	// Setup
	router := setupRouter()

	// Define test cases
	const senderAddress = "0x0A7199a96fdf0252E09F76545c1eF2be3692F46b"
	testCases := []struct {
		description string
		payload     Body
		expectCode  int
	}{
		{
			description: "Valid Request",
			payload: Body{
				Sender: senderAddress,
				Intents: []Intent{
					{
						Sender:     senderAddress,
						Kind:       "buy",
						SellToken:  "TokenA",
						BuyToken:   "TokenB",
						SellAmount: 10.0,
						BuyAmount:  5.0,
						Status:     "Received",
					},
				},
			},
			expectCode: http.StatusOK,
		},
		{
			description: "Invalid Request (missing fields)",
			payload: Body{
				Sender: "0xInValidEthAddress",
			},
			expectCode: http.StatusBadRequest,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/submit", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectCode {
				t.Errorf("Expected status code %d, got %d", tc.expectCode, w.Code)
			}
		})
	}
}
