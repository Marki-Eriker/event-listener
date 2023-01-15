package encoder

type Encoder interface {
	Decode([]byte) []byte
	Encode([]byte) []byte
}
