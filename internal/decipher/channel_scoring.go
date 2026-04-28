package decipher

// calculateChannelScore computes a score from signals.
func calculateChannelScore(signals []Signal) float64 {
	if len(signals) == 0 {
		return 0.0
	}

	var totalWeight float64
	var verifiedWeight float64
	var plausibleWeight float64
	var speculativeWeight float64

	for _, s := range signals {
		totalWeight += s.Weight

		switch s.Confidence {
		case ConfidenceVerified:
			verifiedWeight += s.Weight
		case ConfidencePlausible:
			plausibleWeight += s.Weight
		case ConfidenceSpeculative:
			speculativeWeight += s.Weight
		}
	}

	// Weight verified signals more heavily
	score := (verifiedWeight * 1.0) + (plausibleWeight * 0.6) + (speculativeWeight * 0.3)
	score /= totalWeight

	return score
}

// itoa converts int to string (simple implementation).
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
