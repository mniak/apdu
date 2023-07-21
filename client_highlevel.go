package apdu

import (
	"errors"

	"github.com/mniak/krypton/encoding/tlv"
)

type HighLevelCommands interface {
	SelectByName(dfname []byte) (FileControlInformation, error)
	ReadRecord(sfi, recordNumber int) (RecordTemplate, error)

	// ReadAllRecords tries to read the records of a file starting from record 1. When a
	// status 6A83 (Record Not Found) is returned, it considers that the sequence ended
	// and returns the previous records returned.
	ReadAllRecords(sfi int) ([]RecordTemplate, error)
	GetPSE(contactless bool) ([]RecordTemplate, error)
	GetProcessingOptions(pdolData []byte) (EMVResponseMessageTemplateFormat2, error)
	GenerateARQC(cdolData []byte) (GenerateACResponse, error)
}

type _HighLevelClient struct {
	Low LowLevelCommands
}

func (c _HighLevelClient) GetProcessingOptions(pdolData []byte) (EMVResponseMessageTemplateFormat2, error) {
	return unmarshal[EMVResponseMessageTemplateFormat2](
		c.Low.GetProcessingOptions(pdolData),
	)
}

func (c _HighLevelClient) SelectByName(dfname []byte) (FileControlInformation, error) {
	return unmarshal[FileControlInformation](
		c.Low.SelectByName(dfname),
	)
}

func (c _HighLevelClient) ReadRecord(sfi, recordNumber int) (RecordTemplate, error) {
	return unmarshal[RecordTemplate](
		c.Low.ReadRecord(sfi, recordNumber),
	)
}

func (c _HighLevelClient) ReadAllRecords(sfi int) ([]RecordTemplate, error) {
	var result []RecordTemplate
	recordNumber := 1
	for {
		record, err := c.ReadRecord(sfi, recordNumber)
		if errors.Is(err, ErrRecordNotFound) {
			break
		}
		if err != nil {
			return nil, err
		}

		recordNumber++
		result = append(result, record)
	}
	return result, nil
}

func (c _HighLevelClient) GetPSE(contactless bool) ([]RecordTemplate, error) {
	const PSE1 = "1PAY.SYS.DDF01"
	const PSE2 = "2PAY.SYS.DDF01"

	dfname := []byte(PSE1)
	if contactless {
		dfname = []byte(PSE2)
	}

	fci, err := c.SelectByName(dfname)
	if err != nil {
		return nil, err
	}

	if fci.FCITemplate == nil || fci.FCITemplate.ProprietaryInformation == nil || fci.FCITemplate.ProprietaryInformation.ShortFileIdentifier == nil {
		return nil, nil
	}

	records, err := c.ReadAllRecords(*fci.FCITemplate.ProprietaryInformation.ShortFileIdentifier)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (c _HighLevelClient) GenerateARQC(transactionData []byte) (GenerateACResponse, error) {
	return unmarshal[GenerateACResponse](
		c.Low.GenerateAC(ARQC, transactionData),
	)
}

func unmarshal[T any](data []byte, err error) (T, error) {
	var result T
	if err != nil {
		return result, err
	}

	err = tlv.UnmarshalBER(data, &result)
	if err != nil {
		return result, err
	}

	return result, err
}
