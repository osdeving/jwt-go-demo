# Entendendo Base64

Como funciona uma codificação? Imagine que você tem uma regra que diz o seguinte: caso a entrada for 'a', então o valor será 97, caso a entrada for 'b', então o valor será 98, caso for 'c', então o valor será 99, e assim em diante. Isso é um de-para e é também conhecido como tabela ASCII e envolve outros caracteres além do alfabeto. Descrevi um tipo de codificação. 

Base64 é uma regra de codificação (Encode) que aplicamos na entrada produzindo uma saída. Isso significa que podemos lêr a entrada em base64 e produzir uma saída decodificada ou ler uma entrada decodificada e codificar em base64, respectivamente conhecidos como Encode e Decode. Diferença da tabela ASCII, a entrada pode ser qualquer coisa e não se limita apenas a caracteres individuais. No caso do Base64, a codificação funciona agrupando a entrada em blocos de 3 bytes (24 bits) e dividindo-os em 4 grupos de 6 bits. Cada um desses grupos de 6 bits é então convertido para um caractere correspondente em uma tabela de 64 símbolos, que inclui letras maiúsculas e minúsculas do alfabeto, números e dois caracteres especiais (+ e / no Base64 padrão). Quando a entrada não é um múltiplo exato de 3 bytes, a codificação adiciona um caractere de preenchimento (=) para manter a consistência.

Diferente da tabela ASCII, onde cada caractere corresponde diretamente a um número fixo, no Base64 a conversão ocorre em blocos, garantindo que qualquer sequência de bytes possa ser representada apenas com caracteres seguros para transporte em protocolos como e-mails (MIME), URLs e JSON. Isso torna o Base64 útil para codificar binários, imagens, chaves criptográficas e outros dados que não são diretamente representáveis como texto legível.

O processo inverso, conhecido como decodificação (Decode), pega uma string codificada em Base64 e reconstrói os bytes originais, revertendo a conversão de 6 bits para 8 bits e removendo qualquer padding (=) que tenha sido adicionado na codificação. Assim como podemos converter a → 97 na tabela ASCII e depois reverter 97 → a, no Base64 podemos pegar TWFu e recuperar Man.

No entanto, é importante lembrar que Base64 não é uma criptografia, pois qualquer pessoa pode decodificá-lo facilmente. Ele é apenas um método de representação de dados que facilita o transporte e armazenamento em sistemas que não suportam caracteres binários diretamente.

## Como funciona na prática?

Vamos supor que a entrada sera a string "Rox". Cada caractere em um computador é representado por um número inteiro conforme a tabela ASCII. 

### Passo 1: Obter os valores ASCII

Primeiro, vamos converter os caracteres em seus valores numéricos:

* R → 82
* o → 111
* x → 120

Agora, precisamos representar esses números em binário (base 2), pois a codificação Base64 trabalha diretamente com bits:

### Passo 2: Converter os valores ASCII para binário (8 bits cada)

R (82) → 01010010
o (111) → 01101111
x (120) → 01111000

Agora temos um total de 3 bytes (24 bits), que são agrupados assim:

01010010 01101111 01111000

### Passo 3: Dividir em blocos de 6 bits

O Base64 trabalha com grupos de 6 bits, então precisamos separar nossos 24 bits assim:

010100 100110 111101 111000

Agora, cada grupo de 6 bits será convertido em um decimal cujo qual representa um índice na tabela Base64.

### Passo 4: Converter os grupos de 6 bits para decimal

* 010100 → 20
* 100110 → 38
* 111101 → 61
* 111000 → 56

Agora, usamos a tabela Base64 para converter esses números em caracteres.

### Passo 5: Mapear os valores na tabela Base64

A tabela Base64 contém 64 caracteres indexados de 0 até 63. Essa tabela segue uma ordem específica:

1) Letras maiúsculas (A-Z) → índices 0 a 25
2) Letras minúsculas (a-z) → índices 26 a 51
3) Números (0-9) → índices 52 a 61
4) Caracteres especiais (+ e /) → índices 62 e 63

Representação em Go:

```go
const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
```
Agora, usamos essa tabela para mapear os valores 20, 38, 61 e 56:

* base64Table[20] → "U"
* base64Table[38] → "m"
* base64Table[61] → "9"
* base64Table[56] → "4"

Com isso, o texto "Rox" foi convertido para Base64 como "Um94"

### Resumo

1 - Convertemos os caracteres para seus valores ASCII (R = 82, o = 111, x = 120).

2 - Transformamos esses valores em binário (01010010 01101111 01111000).

3 - Dividimos em blocos de 6 bits (010100 100110 111101 111000).

4 - Convertemos cada grupo de 6 bits em um número decimal (20, 38, 61, 56).

5 - Mapeamos esses números na tabela Base64, resultando na string "Um94".

Bem tranquilo, certo?

## Implementação em Go

O maior desafio que teremos é agrupar os bits em grupos de 6 bits, pois estamos trabalhando com 8 bits (1 byte). Isso implica que, a cada 3 bytes (24 bits no total), teremos 4 grupos de 6 bits, pois 24 / 6 = 4. 


---
<summary>Extra: Entendendo o bitwise AND

<details>
Se tivermos um valor binário qualquer e aplicarmos uma máscara bit a bit (bitwise AND) com 0b111111 (0x3F em hexadecimal), conseguimos extrair exatamente 6 bits da posição desejada.

Isso funciona porque o operador AND (&), mantém apenas os bits onde há 1 nos dois operadores, então podemos isolar porções específicas de um número maior.

#### Exemplo prático:

```
BYTE qualquer:            01010010  (82 em decimal)
Máscara de 6 bits:        00111111  (0x3F em hexadecimal)
Resultado após AND:       00010010  (18 em decimal)
```

Note que AND tem o poder de desligar o bit ou mantê-lo sem modificação. 

1 AND 1 = 1
1 AND 0 = 0
0 AND 0 = 0

Se você aplicar a máscara b00000000, vai apagar tudo

Se você aplicar a máscara b00000001 vai apagar tudo e, para o primeiro bit, vai depender se o outro valor tem 1 ou 0: se tiver 1, ele será mantido, se tiver 0, ele apaga.

Outro exemplo: pegar apenas os bits 3 e 4 de um byte qualquer

```
BYTE qualquer:            10110110  (182 em decimal)
Máscara de 3 bits:        00011000  (24 em decimal)
Resultado após AND:       00010000  (16 em decimal)
```
</details>
</summary>

---

Podemos trabalhar diretamente com 6 bytes em um int64 já que poderíamos percorrer a mensagem de entrada de 6 em 6 bytes e para cada passagem jogar os 6 bytes no int64 e extrair os valores usando máscara de bits.

```go
for i := 0; i < len(input); i += 6 {
    
    block := int64(input[i])     << 40 | 
             int64(input[i + 1]) << 32 |
             int64(input[i + 2]) << 24 |
             int64(input[i + 3]) << 16 | 
             int64(input[i + 4]) << 8  | 
             int64(input[i + 5])

    output = append(output,
        base64Table[(block >> 42) & 0b00111111],
        base64Table[(block >> 36) & 0b00111111],
        base64Table[(block >> 30) & 0b00111111],
        base64Table[(block >> 24) & 0b00111111],
        base64Table[(block >> 18) & 0b00111111],
        base64Table[(block >> 12) & 0b00111111],
        base64Table[(block >>  6) & 0b00111111],
        base64Table[ block        & 0b00111111],
    )
}
```

Agora só precisamos lidar com casos onde o input tem menos que 6 bytes ou o input não é múltiplo de 6 bytes

```go
len := len(input)
rem := len % 6

for i := 0; i < len - rem; i += 6 {
    ...
}

if rem > 0 {
    var block int64 = 0 |
        (int64(input[len -rem     ]) << 40) * int64(rem >= 1) |
        (int64(input[len - rem + 1]) << 32) * int64(rem >= 2) |
        (int64(input[len - rem + 2]) << 24) * int64(rem >= 3) |
        (int64(input[len - rem + 3]) << 16) * int64(rem >= 4) |
        (int64(input[len - rem + 4]) <<  8) * int64(rem >= 5)

    
    base64Chunk := [4]byte{
        base64Table[(block >> 42) & 0b00111111],
        base64Table[(block >> 36) & 0b00111111],
        base64Table[(block >> 30) & 0b00111111],
        base64Table[(block >> 24) & 0b00111111],
    }

    pad := 4 - ((rem * 4) / 3)
    base64Chunk[3] = '=' * byte(pad>>0 & 1) // Se pad >= 1, substitui o último
    base64Chunk[2] = '=' * byte(pad>>1 & 1) // Se pad == 2, substitui o penúltimo

    output = append(output, base64Chunk[:]...)
}


```

O código completo fica:

```go
package main

const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func EncodeBase64(input []byte) string {
	
    var output []byte
	len := len(input)
	rem := len % 6

    for i := 0; i < len(input); i += 6 {
    
        block := int64(input[i])     << 40 | 
                 int64(input[i + 1]) << 32 |
                 int64(input[i + 2]) << 24 |
                 int64(input[i + 3]) << 16 | 
                 int64(input[i + 4]) << 8  | 
                 int64(input[i + 5])

        output = append(output,
            base64Table[(block >> 42) & 0b00111111],
            base64Table[(block >> 36) & 0b00111111],
            base64Table[(block >> 30) & 0b00111111],
            base64Table[(block >> 24) & 0b00111111],
            base64Table[(block >> 18) & 0b00111111],
            base64Table[(block >> 12) & 0b00111111],
            base64Table[(block >>  6) & 0b00111111],
            base64Table[ block        & 0b00111111],
    )

    
    var block int64 = 0 |
        (int64(input[len - rem    ]) << 40) * int64(rem >= 1) |
        (int64(input[len - rem + 1]) << 32) * int64(rem >= 2) |
        (int64(input[len - rem + 2]) << 24) * int64(rem >= 3) |
        (int64(input[len - rem + 3]) << 16) * int64(rem >= 4) |
        (int64(input[len - rem + 4]) <<  8) * int64(rem >= 5)


    base64Chunk := [4]byte{
        base64Table[(block >> 42) & 0b00111111],
        base64Table[(block >> 36) & 0b00111111],
        base64Table[(block >> 30) & 0b00111111],
        base64Table[(block >> 24) & 0b00111111],
    }

    pad := 4 - ((rem * 4) / 3)
    base64Chunk[3] = '=' * byte(pad >> 0 & 1)
    base64Chunk[2] = '=' * byte(pad >> 1 & 1)

    output = append(output, base64Chunk[:]...)
    
    return string(output)
}

func main() {
    testCases := []string{"Man", "Hello", "T", "Testing123", "Base64!", "abcdefg"}
	
    for _, test := range testCases {
		encoded := EncodeBase64([]byte(test))
		println("Base64 de ", test, ": ", encoded)
	}
}
















