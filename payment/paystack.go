package payment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// PaystackAPI structure to hold Paystack API keys and configurations
type PaystackAPI struct {
	APIKey  string
	BaseURL string
}

// NewPaystackAPI returns a new instance of the Paystack API client
func NewPaystackAPI(apiKey string) *PaystackAPI {
	return &PaystackAPI{
		APIKey:  apiKey,
		BaseURL: "https://api.paystack.co", // Paystack base URL
	}
}

// InitializeTransaction initiates a payment request to Paystack
func (p *PaystackAPI) InitializeTransaction(amount float64, reference string, mobileMoneyPhone string, mobileMoneyProvider string) (string, error) {
	client := resty.New()

	// Paystack requires amount in pesewas (for GHS), which is 100 sub-units of Cedi
	pesewasAmount := int(amount * 100)

	// Create the request body
	requestBody := map[string]interface{}{
		"email":     "mineracail@gmail.com",
		"amount":    pesewasAmount,
		"reference": reference,
		"currency":  "GHS", // Set the currency as needed
	}

	// Debug: Print the API Key
	fmt.Println("API Key:", p.APIKey) // Debug line; remove in production

	  // Include mobile money information if provided
      if mobileMoneyPhone != "" && mobileMoneyProvider != "" {
        requestBody["mobile_money"] = map[string]interface{}{
            "phone":    mobileMoneyPhone,
            "provider": mobileMoneyProvider,
        }
    }

    // Send the request to Paystack API
    resp, err := client.R().
        SetHeader("Authorization", fmt.Sprintf("Bearer %s", p.APIKey)).
        SetHeader("Content-Type", "application/json").
        SetBody(requestBody).
        Post(fmt.Sprintf("%s/transaction/initialize", p.BaseURL))

    if err != nil {
        return "", fmt.Errorf("failed to initiate Paystack request: %v", err)
    }

    if resp.StatusCode() != http.StatusOK {
        return "", fmt.Errorf(" Paystack API error: %v", resp.String())
    }

    // Parse the JSON response body
    var result map[string]interface{}
    err = json.Unmarshal(resp.Body(), &result)
    if err != nil {
        return "", fmt.Errorf("error decoding Paystack response: %v", err)
    }

    // Extract the "authorization_url" from the parsed response
    data, ok := result["data"].(map[string]interface{})
    if !ok {
        return "", fmt.Errorf("missing or invalid 'data' field in response")
    }

    authorizationURL, ok := data["authorization_url"].(string)
    if !ok {
        return "", fmt.Errorf("missing or invalid 'authorization_url' field")
    }

    // Return the authorization URL
    return authorizationURL, nil
}

// VerifyTransaction verifies the transaction status based on the reference
func (p *PaystackAPI) VerifyTransaction(reference string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", p.APIKey)).
		Get(fmt.Sprintf("%s/transaction/verify/%s", p.BaseURL, reference))

	if err != nil {
		return "", fmt.Errorf("failed to verify Paystack transaction: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("error retrieving transaction status: %v", resp.String())
	}

	return resp.String(), nil
}
