package apdu

import (
	"io"
	"log"

	"github.com/mniak/apdu/internal/noop"
	"github.com/mniak/apdu/internal/utils"
	"github.com/mniak/krypton/encoding/tlv"
)

//go:generate mockgen -package=apdu -destination=client_raw_mock.go -source=client_raw.go

type RawClient interface {
	SendCommand(cmd Command) (Response, error)
}

type _RawClient struct {
	driver        Driver
	lengthEncoder tlv.LengthEncoder
	logger        *log.Logger
}

func (d *_RawClient) LoggingTo(w io.Writer) *_RawClient {
	d.logger = log.New(w, "[Client] ", 0)
	return d
}

func NewRawClient(driver Driver) RawClient {
	result := &_RawClient{
		driver:        driver,
		lengthEncoder: tlv.ShortLengthEncoder,
		logger:        noop.Logger(),
	}
	return result
}

func (c _RawClient) internalSendCommand(cmd Command) (Response, error) {
	bytes, err := cmd.Bytes(c.lengthEncoder)
	if err != nil {
		return Response{}, err
	}
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

	if resp.HasWrongLength() {
		cmd.MaxReponseLength = resp.Trailer.SW2()
		resp, err = c.internalSendCommand(cmd)
	}

	for resp.HasMoreData() {
		moreDataCmd := Command{
			Instruction:      0xC0,
			MaxReponseLength: resp.Trailer.SW2(),
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

	return resp, err
}

func (c _RawClient) SendCommand(cmd Command) (Response, error) {
	cmdbytes, err := cmd.Bytes(c.lengthEncoder)
	if err != nil {
		return Response{}, err
	}
	c.logger.Printf("APDU sent: %2X\n%s\n", cmdbytes, utils.IndentString(cmd.StringPretty(), "  "))

	resp, err := c.internalSendCommand(cmd)
	c.logger.Printf("APDU received: [%2X] [%02X %02X]\n", resp.Data, resp.Trailer.SW1(), resp.Trailer.SW2())
	return resp, err
}
