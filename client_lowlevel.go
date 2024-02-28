package apdu

type LowLevelCommands interface {
	SelectByName(dfname []byte) ([]byte, error)
	ReadRecord(sfi, recordNumber int) ([]byte, error)
	GetProcessingOptions(pdolData []byte) ([]byte, error)
	GenerateAC(cryptogramType ApplicationCryptogramType, transactionData []byte) ([]byte, error)
}

type _LowLevelClient struct {
	RawClient
}

func (c _LowLevelClient) SelectByName(dfname []byte) ([]byte, error) {
	resp, err := c.SendCommand(Command{
		Class:       0x00,
		Instruction: InstructionA4_Select,
		Parameters: Parameters{
			P1: 0x04, // Select by DF name
			P2: 0x00, // First or only occurrence
		},
		Data: []byte(dfname),
	})
	if err != nil {
		return nil, err
	}

	return resp.Data, resp.Trailer.GetError()
}

func (c _LowLevelClient) ReadRecord(sfi, recordNumber int) ([]byte, error) {
	cmd := Command{
		Class:       0x00,
		Instruction: InstructionB2_ReadRecords,
		Parameters: Parameters{
			P1: byte(recordNumber),   // Record number
			P2: byte(sfi<<3) | 0b100, // Short EF ID | P1 is the record number
		},
		Data: nil,
	}
	resp, err := c.SendCommand(cmd)
	if err != nil {
		return nil, err
	}

	return resp.Data, resp.Trailer.GetError()
}

func (c _LowLevelClient) GetProcessingOptions(pdolData []byte) ([]byte, error) {
	cmd := Command{
		Class:       0x80,
		Instruction: EMVInstructionA8_GetProcessingOptions,
		Parameters: Parameters{
			P1: 0x00,
			P2: 0x00,
		},
		Data: pdolData,
	}
	resp, err := c.SendCommand(cmd)
	if err != nil {
		return nil, err
	}

	return resp.Data, resp.Trailer.GetError()
}

type ApplicationCryptogramType byte

const (
	AAC  ApplicationCryptogramType = 0b00
	TC   ApplicationCryptogramType = 0b01
	ARQC ApplicationCryptogramType = 0b10
)

func (c _LowLevelClient) GenerateAC(cryptogramType ApplicationCryptogramType, transactionData []byte) ([]byte, error) {
	cmd := Command{
		Class:       0x80,
		Instruction: EMVInstructionA8_GenerateAC,
		Parameters: Parameters{
			P1: byte(cryptogramType << 6),
			P2: 0x00,
		},
		Data: transactionData,
	}
	resp, err := c.SendCommand(cmd)
	if err != nil {
		return nil, err
	}
	return resp.Data, resp.Trailer.GetError()
}
