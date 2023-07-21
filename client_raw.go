package apdu

import (
	"log"

	"github.com/mniak/apdu/internal/noop"
	"github.com/mniak/apdu/internal/utils"
	"github.com/mniak/krypton/encoding/tlv"
)

type RawClient interface {
	SendCommand(cmd Command) (Response, error)
}

type _RawClient struct {
	driver        Driver
	lengthEncoder tlv.LengthEncoder
	logger        *log.Logger
}

func NewRawClient(driver Driver) RawClient {
	return _RawClient{
		driver:        driver,
		lengthEncoder: tlv.ShortLengthEncoder,
		logger:        noop.Logger(),
	}
}

func (c _RawClient) SendCommand(cmd Command) (Response, error) {
	bytes, err := cmd.Bytes(c.lengthEncoder)
	if err != nil {
		return Response{}, err
	}
	c.logger.Printf("APDU sent: %2X\n%s\n", bytes, utils.IndentString(cmd.StringPretty(), "  "))
	responseBytes, err := c.driver.SendBytes(bytes)
	if err != nil {
		return Response{}, err
	}
	resp, err := ParseResponse(responseBytes)
	if err != nil {
		return resp, err
	}

	// var buffer bytes.Buffer
	// for {
	// 	if has, length := resp.HasMoreData(); has {
	// 		moreDataCmd := Command{
	// 			Instruction:      0xC0,
	// 			MaxReponseLength: uint32(length),
	// 		}
	// 		resp2Bytes, err := a.Driver.SendBytes(moreDataCmd.Bytes())
	// 		if err != nil {
	// 			return Response{}, err
	// 		}
	// 		resp, err = ParseResponse(resp2Bytes)
	// 	} else {
	// 		break
	// 	}
	// }
	// return Response{
	// 	Trailer: resp.Trailer,
	// 	Data:    Data(buffer.Bytes()),
	// }, err

	if has, length := resp.HasMoreData(); has {
		moreDataCmd := Command{
			Instruction:      0xC0,
			MaxReponseLength: byte(length),
		}

		var moreDataCmdBytes []byte
		moreDataCmdBytes, err = moreDataCmd.Bytes(c.lengthEncoder)
		if err != nil {
			return Response{}, err
		}

		var resp2Bytes []byte
		resp2Bytes, err = c.driver.SendBytes(moreDataCmdBytes)
		if err != nil {
			return Response{}, err
		}

		resp, err = ParseResponse(resp2Bytes)
	}
	c.logger.Printf("APDU received: [%2X] [%02X %02X]\n", resp.Data, resp.Trailer.SW1(), resp.Trailer.SW2())
	return resp, err
}
