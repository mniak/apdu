package apdu

import (
	"errors"

	"github.com/mniak/krypton/encoding/tlv"
)

type HighLevelClient struct {
	Commander
	LowLevel LowLevelClient
}

func NewHighLevelClient(driver Driver) HighLevelClient {
	lowLevelClient := NewLowLevelClient(driver)
	return HighLevelClient{
		Commander: lowLevelClient,
		LowLevel:  lowLevelClient,
	}
}

func (c HighLevelClient) SelectByName(dfname []byte) (FileControlInformation, error) {
	resp, err := c.LowLevel.SelectByName(dfname)
	if err != nil {
		return FileControlInformation{}, err
	}

	var fci FileControlInformation
	if err := tlv.UnmarshalBER(resp, &fci); err != nil {
		return FileControlInformation{}, err
	}
	return fci, nil
}

func (c HighLevelClient) ReadRecord(recordNumber, fileID int) (*RecordTemplate, error) {
	data, err := c.LowLevel.ReadRecord(recordNumber, fileID)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, ErrRecordNotFound
	}

	var record RecordTemplate
	if err := tlv.UnmarshalBER(data, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (c HighLevelClient) ReadAllRecords(fileID int) ([]RecordTemplate, error) {
	var result []RecordTemplate
	recordNumber := 1
	for {
		record, err := c.ReadRecord(recordNumber, fileID)
		if errors.Is(err, ErrRecordNotFound) {
			break
		}
		if err != nil {
			return nil, err
		}

		recordNumber++
		result = append(result, *record)
	}
	return result, nil
}

func (c HighLevelClient) GetPSE(contactless bool) ([]RecordTemplate, error) {
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
