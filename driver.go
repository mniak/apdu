package apdu

type Driver interface {
	SendBytes(bytes []byte) ([]byte, error)
}
