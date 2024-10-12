    package payment

    import (
        "fmt"
        "net/http"

        "github.com/go-resty/resty/v2"
        "github.com/google/uuid"
    )

    // MoMoAPI structure to hold MoMo API keys and configurations
    type MoMoAPI struct {
        APIKey          string
        SubscriptionKey string
        UserID          string
        BaseURL         string
        Environment     string 
    }

    // NewMoMoAPI returns a new instance of the MoMo API client
    func NewMoMoAPI(apiKey, subscriptionKey, environment string) *MoMoAPI {
        baseURL := "https://sandbox.momodeveloper.mtn.com/collection/v1_0/requesttopay"
    // SubscriptionKey := "https://sandbox.momodeveloper.mtn.com/collection/v1_0/requesttopay"
        if environment == "production" {
            baseURL = "https://momodeveloper.mtn.com/collection/v1_0/requesttopay"
        }

        return &MoMoAPI{
            APIKey:          apiKey,
            SubscriptionKey: subscriptionKey,
            Environment:     environment,
            BaseURL:         baseURL,
        }
    }




    func DeductFeeFromCollection(amount float64) float64 {
        fee := CalculateTransactionFee(amount)
        return amount - fee // Deduct the fee from the collected amount
    }

    // RequestToPay initiates a payment request to MoMo for a collection transaction
    func (m *MoMoAPI) RequestToPay(amount float64, phoneNumber string, externalID string) (string, error) {
        client := resty.New()

        // Create a new UUID for the X-Reference-Id header
        referenceID := uuid.New().String()

        // Deduct fee from the original amount
        netAmount := DeductFeeFromCollection(amount)

        // Create the request body
        requestBody := map[string]interface{}{
            "amount":       fmt.Sprintf("%.2f", netAmount),
            "currency":     "GHS", // Use Ghanaian Cedis for Ghana
            "externalId":   externalID,
            "payer": map[string]string{
                "partyIdType": "MSISDN",
                "partyId":     phoneNumber, // Customer's phone number
            },
            "payerMessage": "Payment for product/service",
            "payeeNote":    "Thank you for your business",
        }

        // Send the request to MoMo API
        resp, err := client.R().
            SetHeader("Ocp-Apim-Subscription-Key", m.SubscriptionKey).
            SetHeader("X-Reference-Id", referenceID). // Use the generated UUID
            SetHeader("X-Target-Environment", m.Environment). // Use the environment variable
            SetBody(requestBody).
            Post(m.BaseURL)

        if err != nil {
            return "", fmt.Errorf("failed to initiate MoMo request: %v", err)
        }

        if resp.StatusCode() != http.StatusAccepted {
            return "", fmt.Errorf("MoMo API error: %v", resp.String())
        }

        return referenceID, nil // Return the reference ID for tracking
    }

    // CheckPaymentStatus retrieves the payment status based on the reference ID
    func (m *MoMoAPI) CheckPaymentStatus(referenceID string) (string, error) {
        client := resty.New()
        resp, err := client.R().
            SetHeader("Ocp-Apim-Subscription-Key", m.SubscriptionKey).
            SetHeader("X-Target-Environment", m.Environment).
            Get(fmt.Sprintf("%s/%s", m.BaseURL, referenceID))

        if err != nil {
            return "", fmt.Errorf("failed to check payment status: %v", err)
        }

        if resp.StatusCode() != http.StatusOK {
            return "", fmt.Errorf("error retrieving payment status: %v", resp.String())
        }

        return resp.String(), nil
    }
