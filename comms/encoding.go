package comms

import "io"

type Decoder interface {
	Decode(io.Reader, any) error
}
