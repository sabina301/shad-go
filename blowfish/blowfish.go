//go:build !solution

package blowfish

// #cgo pkg-config: libcrypto
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <openssl/blowfish.h>
import "C"
import "unsafe"

type Blowfish struct {
	key []byte
}

func (b Blowfish) BlockSize() int {
	return 8
}

func (b Blowfish) Encrypt(bytes []byte, bytes2 []byte) {
	var bfKey C.BF_KEY
	var keyLen = len(b.key)

	C.BF_set_key(
		&bfKey,
		C.int(keyLen),
		(*C.uchar)(unsafe.Pointer(&b.key[0])),
	)

	C.BF_ecb_encrypt(
		(*C.uchar)(unsafe.Pointer(&bytes2[0])),
		(*C.uchar)(unsafe.Pointer(&bytes[0])),
		&bfKey,
		C.BF_ENCRYPT,
	)
}

func (b Blowfish) Decrypt(bytes []byte, bytes2 []byte) {
	var bfKey C.BF_KEY
	var keyLen = len(b.key)
	C.BF_set_key(
		&bfKey,
		C.int(keyLen),
		(*C.uchar)(unsafe.Pointer(&b.key[0])),
	)

	C.BF_ecb_encrypt(
		(*C.uchar)(unsafe.Pointer(&bytes2[0])),
		(*C.uchar)(unsafe.Pointer(&bytes[0])),
		&bfKey,
		C.BF_DECRYPT,
	)
}

func New(key []byte) *Blowfish {
	return &Blowfish{
		key: key,
	}
}
