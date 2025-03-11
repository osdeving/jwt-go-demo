//go:build ignore
package main

import (
    "fmt"
)

// Valores iniciais do hash (definidos na especificação)
var h = [8]uint32{
    0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
    0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
}

// Constantes SHA-256
var k = [64]uint32{
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


func padMessage(msg []byte) []byte {
    origLen := uint64(len(msg) * 8) // Tamanho original em bits
    msg = append(msg, 0x80)         // Adiciona bit 1
    for len(msg)%64 != 56 {
        msg = append(msg, 0x00)     // Preenche com 0 até ter 56 bytes restantes
    }
    for i := 0; i < 8; i++ {
        msg = append(msg, byte(origLen>>(8*(7-i)))) // Adiciona o tamanho original
    }
    return msg
}

func expandWords(block []byte) [64]uint32 {
    var w [64]uint32
    for i := 0; i < 16; i++ {
        w[i] = uint32(block[i*4])<<24 | uint32(block[i*4+1])<<16 |
               uint32(block[i*4+2])<<8 | uint32(block[i*4+3])
    }
    for i := 16; i < 64; i++ {
        s0 := (w[i-15]>>7 | w[i-15]<<25) ^ (w[i-15]>>18 | w[i-15]<<14) ^ (w[i-15]>>3)
        s1 := (w[i-2]>>17 | w[i-2]<<15) ^ (w[i-2]>>19 | w[i-2]<<13) ^ (w[i-2]>>10)
        w[i] = w[i-16] + s0 + w[i-7] + s1
    }
    return w
}


func processBlock(w [64]uint32, h *[8]uint32) {
    a, b, c, d, e, f, g, hh := h[0], h[1], h[2], h[3], h[4], h[5], h[6], h[7]
    for i := 0; i < 64; i++ {
        s1 := (e>>6 | e<<26) ^ (e>>11 | e<<21) ^ (e>>25 | e<<7)
        ch := (e & f) ^ (^e & g)
        temp1 := hh + s1 + ch + k[i] + w[i]
        s0 := (a>>2 | a<<30) ^ (a>>13 | a<<19) ^ (a>>22 | a<<10)
        maj := (a & b) ^ (a & c) ^ (b & c)
        temp2 := s0 + maj
        hh, g, f, e, d, c, b, a = g, f, e, d+temp1, c, b, a, temp1+temp2
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

func sha256(msg []byte) [32]byte {
    msg = padMessage(msg)
    for i := 0; i < len(msg); i += 64 {
        w := expandWords(msg[i : i+64])
        processBlock(w, &h)
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

func main() {
	msg := []byte("Uma mensagem")
	hash := "1fceeb0e740fc91e820655ea5d49535ee37e757f674fabf647ab90cb53b3ea76"
	digest := sha256(msg)
	fmt.Printf("%x\n", digest)

	if fmt.Sprintf("%x", digest) == hash {
		fmt.Println("Hashes iguais")
	} else {
		fmt.Println("Hashes diferentes")
	}
}

