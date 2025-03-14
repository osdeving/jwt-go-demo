package main

import (
	"fmt"
	"encoding/base64"
)

const b64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func Encb64(in []byte) string {
	
    var out []byte
	len := len(in)
	rem := len % 3

    for i := 0; i < len - rem; i += 3 {
    
        blk := int64(in[i    ]) << 16 | 
               int64(in[i + 1]) << 8  | 
               int64(in[i + 2])

        out = append(out,
            b64[(blk >> 18) & 0b00111111],
            b64[(blk >> 12) & 0b00111111],
            b64[(blk >>  6) & 0b00111111],
            b64[ blk        & 0b00111111],
		)
	}

	if rem == 0 {
		return string(out)
	}
    
	var blk int64 = 0

 	if rem == 1 {
        blk = int64(in[len - rem]) << 16
        out = append(out,
            b64[(blk >> 18) & 0b00111111],
            b64[(blk >> 12) & 0b00111111],
            '=',
            '=',
        )
    } else if rem == 2 {
        blk = int64(in[len - rem]) << 16 | int64(in[len - rem + 1]) << 8
        out = append(out,
            b64[(blk >> 18) & 0b00111111],
            b64[(blk >> 12) & 0b00111111],
            b64[(blk >> 6)  & 0b00111111],
            '=',
        )
    }

	return string(out)
}

func main() {
    testCases := []string{"Manaed", "1234567890"}
	
    for _, test := range testCases {
		encoded := Encb64([]byte(test))
		encondedGo := base64.StdEncoding.EncodeToString([]byte(test))
		fmt.Println("Base64 de ", test, ": ", encoded)
		fmt.Println("Base64 Go : ", test, ": ", encondedGo)
		
	}
}