package main

import "strings"

// =============================
//  Estruturas do JWT
// =============================

type Jwt struct {
	Header    string
	Payload   string
	Signature string
}

type Token struct {
	AccessToken string `json:"access_token"`
}

func SIGJWT(header, payload, secret string) Jwt {
	jsonHeaderBase64 := ENCB64URL([]byte(header))
	jsonPayloadBase64 := ENCB64URL([]byte(payload))

	// replace +, / and =
	jsonHeaderBase64 = strings.ReplaceAll(jsonHeaderBase64, "+", "-")
	jsonHeaderBase64 = strings.ReplaceAll(jsonHeaderBase64, "/", "_")

	stringToSign := jsonHeaderBase64 + "." + jsonPayloadBase64
	signature := HMACSHA256(stringToSign, secret)

	jwt := Jwt{
		Header:    jsonHeaderBase64,
		Payload:   jsonPayloadBase64,
		Signature: signature,
	}

	return jwt
}
