package payment

import (
    "fmt"
    "github.com/go-resty/resty/v2"
    "net/http"
)

// VodafoneCashAPI structure to hold Vodafone Cash API keys and configurations
type VodafoneCashAPI struct {
    APIKey     string
    ClientID   string
    ClientSecret string
    BaseURL    string
}

// NewVodafoneCashAPI returns a new instance of the Vodafone Cash API client
func NewVodafoneCashAPI(apiKey, clientID, clientSecret string) *VodafoneCashAPI {
    return &VodafoneCashAPI{
        APIKey:       apiKey,
        ClientID:     clientID,
        ClientSecret: clientSecret,
        BaseURL:      "https://sandbox.vodafonecash.com.gh/v1/collections", // Replace with the actual Vodafone Cash endpoint
    }
}

// RequestToPay initiates a payment request to Vodafone Cash for a collection transaction
func (v *VodafoneCashAPI) RequestToPay(amount float64, phoneNumber string, transactionID string) error {
    client := resty.New()

    // Create the request body
    requestBody := map[string]interface{}{
        "amount":       fmt.Sprintf("%.2f", amount),
        "currency":     "GHS", // Use Ghanaian Cedis
        "transactionId": transactionID, // Unique transaction reference
        "payer": map[string]string{
            "partyIdType": "MSISDN", // Mobile Station International Subscriber Directory Number
            "partyId":     phoneNumber, // Customer's Vodafone Cash phone number
        },
        "payerMessage": "Payment for product/service",
        "payeeNote":    "Thank you for your business",
    }

    // Add authentication headers, Vodafone Cash might use OAuth for token-based authentication
    token, err := v.getAccessToken(client)
    if err != nil {
        return fmt.Errorf("failed to get access token: %v", err)
    }

    // Send the request to Vodafone Cash API
    resp, err := client.R().
        SetHeader("Authorization", "Bearer "+token).
        SetHeader("Ocp-Apim-Subscription-Key", v.APIKey).
        SetHeader("Content-Type", "application/json").
        SetBody(requestBody).
        Post(v.BaseURL + "/requesttopay")

    if err != nil {
        return fmt.Errorf("failed to initiate Vodafone Cash request: %v", err)
    }

    if resp.StatusCode() != http.StatusAccepted {
        return fmt.Errorf("Vodafone Cash API error: %v", resp.String())
    }

    return nil
}

// getAccessToken is used to retrieve OAuth token from Vodafone Cash (if applicable)
func (v *VodafoneCashAPI) getAccessToken(client *resty.Client) (string, error) {
   // authResp, err := client.R().
    _, err := client.R().
        SetFormData(map[string]string{
            "client_id":     v.ClientID,
            "client_secret": v.ClientSecret,
            "grant_type":    "client_credentials",
        }).
        Post(v.BaseURL + "/oauth/token")

    if err != nil {
        return "", fmt.Errorf("failed to get access token: %v", err)
    }

    var result map[string]interface{}
   // if err := authResp.Unmarshal(&result); err != nil {
    //    return "", fmt.Errorf("failed to parse token response: %v", err)
   // }

    token, ok := result["access_token"].(string)
    if !ok {
        return "", fmt.Errorf("access token not found in response")
    }

    return token, nil
}
