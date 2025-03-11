package main

// =============================
//  Estruturas do JWT
// =============================

type Jwt struct {
	Header    string
	Payload   string
	Signature string
}



var K = [64]uint32{
	0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5,
	0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
	0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
	0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
	0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc,
	0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
	0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7,
	0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
	0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13,
	0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
	0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3,
	0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
	0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5,
	0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
	0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208,
	0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
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


