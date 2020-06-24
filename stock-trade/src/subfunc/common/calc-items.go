package common

var TaxRate = 0.1

func CalculateFeeEst(price float64, q int) (totalAmount, fee, tax float64) {
	var payAmount float64

	// 約定金額に応じて手数料決定
	payAmount = price * float64(q)
	switch {
	case payAmount < 250000:
		fee = 1000
	case payAmount < 500000:
		fee = 5000
	default:
		fee = 10000
	}

	tax = float64(fee) * TaxRate

	totalAmount = payAmount + fee + tax

	return fee, tax, totalAmount

}
