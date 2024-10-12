package payment



// Calculate the fee and deduct it from the collected amount
func CalculateTransactionFee(amount float64) float64 {
    fee := amount * 0.02 // 2% fee
    return fee
}

func CeductFeeFromCollection(amount float64) float64 {
    fee := CalculateTransactionFee(amount)
    return amount - fee // Deduct the fee from the collected amount
}
