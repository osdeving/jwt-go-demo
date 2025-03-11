
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

