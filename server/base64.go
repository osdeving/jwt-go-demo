package main

const b64Std = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
const b64URL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

func encodeBase64(in []byte, table string, pad bool) string {
    
    var out []byte
    len := len(in)
    rem := len % 3

    for i := 0; i < len - rem; i += 3 {
        
        blk := int64(in[i    ]) << 16 | 
               int64(in[i + 1]) <<  8 | 
               int64(in[i + 2])

        out = append(out,
            table[(blk >> 18) & 0b00111111],
            table[(blk >> 12) & 0b00111111],
            table[(blk >>  6) & 0b00111111],
            table[ blk        & 0b00111111],
        )
    }

    if rem == 0 {
        return string(out)
    }
    
    var blk int64 = 0

    if rem == 1 {
        blk = int64(in[len - rem]) << 16
        out = append(out,
            table[(blk >> 18) & 0b00111111],
            table[(blk >> 12) & 0b00111111],
        )
        if pad {
            out = append(out, '=', '=')
        }
    } else if rem == 2 {
        blk = int64(in[len - rem]) << 16 | int64(in[len - rem + 1]) << 8
        out = append(out,
            table[(blk >> 18) & 0b00111111],
            table[(blk >> 12) & 0b00111111],
            table[(blk >>  6) & 0b00111111],
        )
        if pad {
            out = append(out, '=')
        }
    }

    return string(out)
}

func ENCB64(in []byte) string {
    return encodeBase64(in, b64Std, true)
}

func ENCB64URL(in []byte) string {
    return encodeBase64(in, b64URL, false)
}
