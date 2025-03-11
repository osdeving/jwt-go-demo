package main

// =============================
//  Estruturas do JWT
// =============================

type Jwt struct {
	Header    string
	Payload   string
	Signature string
}

// =============================
//  Função Principal
// =============================

func main() {
	// Criando o cabeçalho (Header)
	header := "{\"alg\":\"HS256\",\"typ\":\"JWT\"}"

	// Criando o payload (Payload)
	payload := "{\"sub\":\"willams\",\"iat\":1516239022}"


	// Codificação Base64-URL
	jsonHeaderBase64 := encodeBase64URL([]byte(header))
	jsonPayloadBase64 := encodeBase64URL([]byte(payload))

	// Criamos a string que será assinada
	stringToSign := jsonHeaderBase64 + "." + jsonPayloadBase64
	secretKey := "123"

	// Geramos a assinatura HMAC-SHA256 manualmente
	signature := signHMACSHA256(stringToSign, secretKey)

	// Construímos o JWT final
	jwt := Jwt{
		Header:    jsonHeaderBase64,
		Payload:   jsonPayloadBase64,
		Signature: signature,
	}

	// Exibe o JWT final
	println("JWT: " + jwt.Header + "." + jwt.Payload + "." + jwt.Signature)
}

// =============================
//  Implementação do HMAC-SHA256 Manual
// =============================

func signHMACSHA256(message, secret string) string {
    // Create outer and inner padding arrays
    opad := make([]byte, 64)
    ipad := make([]byte, 64)

    // Fill with values
    for i := range opad {
        opad[i] = 0x5c
        ipad[i] = 0x36
    }

    // Process the key
    var key []byte
    if len(secret) > 64 {
        key = sha256([]byte(secret))[:]
    } else {
        key = []byte(secret)
    }

    // XOR the paddings with key
    for i := 0; i < len(key); i++ {
        opad[i] ^= key[i]
        ipad[i] ^= key[i]
    }

    // Inner hash
    inner := append(ipad, []byte(message)...)
    innerHash := sha256(inner)

    // Outer hash
    outer := append(opad, innerHash[:]...)
    outerHash := sha256(outer)

    // Encode final hash to base64url
    return encodeBase64URL(outerHash[:])
}



// =============================
//  Implementação do SHA-256 Manual
// =============================

func sha256(msg []byte) [32]byte {
	// Constantes SHA-256
	h := [8]uint32{
		0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
		0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
	}

	// Padding da mensagem
	msg = padSHA256(msg)

	// Processamento em blocos de 512 bits
	for i := 0; i < len(msg); i += 64 {
		w := make([]uint32, 64)
		for j := 0; j < 16; j++ {
			w[j] = uint32(msg[i+j*4])<<24 | uint32(msg[i+j*4+1])<<16 |
				uint32(msg[i+j*4+2])<<8 | uint32(msg[i+j*4+3])
		}

		for j := 16; j < 64; j++ {
			s0 := (w[j-15]>>7 | w[j-15]<<25) ^ (w[j-15]>>18 | w[j-15]<<14) ^ (w[j-15] >> 3)
			s1 := (w[j-2]>>17 | w[j-2]<<15) ^ (w[j-2]>>19 | w[j-2]<<13) ^ (w[j-2] >> 10)
			w[j] = w[j-16] + s0 + w[j-7] + s1
		}

		a, b, c, d, e, f, g, hh := h[0], h[1], h[2], h[3], h[4], h[5], h[6], h[7]
		for j := 0; j < 64; j++ {
			t1 := hh + ((e>>6 | e<<26) ^ (e>>11 | e<<21) ^ (e>>25 | e<<7)) + ((e & f) ^ (^e & g)) + w[j]
			t2 := ((a>>2 | a<<30) ^ (a>>13 | a<<19) ^ (a>>22 | a<<10)) + ((a & b) ^ (a & c) ^ (b & c))
			hh, g, f, e, d, c, b, a = g, f, e, d+t1, c, b, a, t1+t2
		}

		h[0] += a
		h[1] += b
		h[2] += c
		h[3] += d
		h[4] += e
		h[5] += f
		h[6] += g
		h[7] += hh
	}

	var digest [32]byte
	for i, v := range h {
		digest[i*4] = byte(v >> 24)
		digest[i*4+1] = byte(v >> 16)
		digest[i*4+2] = byte(v >> 8)
		digest[i*4+3] = byte(v)
	}
	return digest
}

// =============================
//  Implementação do Padding SHA-256
// =============================

func padSHA256(msg []byte) []byte {
	origLen := len(msg) * 8
	msg = append(msg, 0x80) // Adiciona bit 1
	for len(msg)%64 != 56 {
		msg = append(msg, 0x00) // Preenche com zeros
	}

	lenBytes := make([]byte, 8)
	for i := uint(0); i < 8; i++ {
		lenBytes[7-i] = byte(origLen >> (i * 8))
	}
	msg = append(msg, lenBytes...)
	return msg
}

// =============================
//  Implementação correta do Base64-URL
// =============================
func encodeBase64URL(data []byte) string {
	const base64URLChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

	var encoded string
	length := len(data)

	for i := 0; i < length; i += 3 {
		val := int(data[i]) << 16
		if i+1 < length {
			val |= int(data[i+1]) << 8
		}
		if i+2 < length {
			val |= int(data[i+2])
		}

		encoded += string(base64URLChars[(val>>18)&0x3F])
		encoded += string(base64URLChars[(val>>12)&0x3F])

		if i+1 < length {
			encoded += string(base64URLChars[(val>>6)&0x3F])
		}
		if i+2 < length {
			encoded += string(base64URLChars[val&0x3F])
		}
	}

	return encoded
}


