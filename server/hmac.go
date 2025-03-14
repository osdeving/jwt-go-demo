package main

// =============================
//  Implementação do HMAC-SHA256 Manual
// =============================

func HMACSHA256(message, secret string) string {
	var key []byte
	if len(secret) > 64 {
		hash := SHA256([]byte(secret))
		key = hash[:]
	} else {
		key = make([]byte, 64)
		copy(key, secret)
	}

	opad := make([]byte, 64)
	ipad := make([]byte, 64)
	for i := range opad {
		opad[i] = 0x5c ^ key[i]
		ipad[i] = 0x36 ^ key[i]
	}

	innerHash := SHA256(append(ipad, []byte(message)...))
	outerHash := SHA256(append(opad, innerHash[:]...))

	return ENCB64URL(outerHash[:])
}