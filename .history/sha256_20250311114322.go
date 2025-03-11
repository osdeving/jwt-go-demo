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
    // ... (restante das 64 constantes)
    0xc67178f2,
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


