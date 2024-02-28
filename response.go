package apdu

import (
	"fmt"
)

type Response struct {
	Data    []byte
	Trailer Trailer
}

func ParseResponse(data []byte) (Response, error) {
	responseLength := len(data)
	if responseLength < 2 {
		return Response{}, fmt.Errorf("invalid response length: %d", responseLength)
	}

	resp := Response{
		Trailer: NewTrailer(data[responseLength-2], data[responseLength-1]),
	}
	data = data[:responseLength-2]

	resp.Data = data
	return resp, nil
}

func (r Response) String() string {
	return fmt.Sprintf("%2X [%s]", r.Data, r.Trailer)
}

func (r Response) HasMoreData() bool {
	return r.Trailer.SW1() == 0x61
}

func (r Response) HasWrongLength() bool {
	return r.Trailer.SW1() == 0x6C
}
