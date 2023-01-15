package encoder

type Mock struct{}

func NewMock() *Mock {
	return &Mock{}
}

func (m Mock) Decode(data []byte) []byte {
	return []byte("mock decoded data")
}

func (m Mock) Encode(data []byte) []byte {
	return []byte("mock encoded data")
}
