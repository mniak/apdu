package apdu

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/mniak/apdu/internal/utils"
	"github.com/mniak/tlv"
)

type PSEFile struct {
	ADFName                      []byte  `tlv:"4F"`
	ApplicationLabel             *string `tlv:"50"`
	ApplicationPriorityIndicator *int    `tlv:"87"`
	ApplicationPreferredName     *string `tlv:"9F12"`
}

type FileControlInformation struct {
	FCPTemplate []byte       `tlv:"62"`
	FMDTemplate []byte       `tlv:"64"`
	FCITemplate *FCITemplate `tlv:"6F"`

	Raw6F string `tlv:"6F,hex"`
}

type FCITemplate struct {
	DFName                 *string                 `tlv:"84,hex"`
	ProprietaryInformation *ProprietaryFCITemplate `tlv:"A5"`

	RawA5 string `tlv:"A5,hex"`
}

type ProprietaryFCITemplate struct {
	ApplicationLabel                 *string           `tlv:"50"`
	ApplicationPriorityIndicator     *int              `tlv:"87"`
	ShortFileIdentifier              *int              `tlv:"88"`
	LanguagePreference               *string           `tlv:"5f2d"`
	ApplicationPreferredNameEncoding *int              `tlv:"9f11"`
	ApplicationPreferredName         *string           `tlv:"9f12"`
	PDOL                             tlv.TL            `tlv:"9f38"`
	IssuerDiscretionaryData          map[string][]byte `tlv:"bf0c"`
}

type RecordTemplate struct {
	EMVProprietaryTemplates []EMVProprietaryTemplate `tlv:"70"`

	// RawData tlv.TLV `tlv:"raw"`
}

type EMVProprietaryTemplate struct {
	ApplicationTemplates []ApplicationTemplate `tlv:"61"`

	Track1DiscretionaryData    string  `tlv:"9f1f"`
	Track2EquivalentData       []byte  `tlv:"57"`
	CardholderName             string  `tlv:"5f20"`
	PAN                        string  `tlv:"5a,hex"`
	PANSequenceNumber          string  `tlv:"5f34,hex"`
	ExpirationDate             string  `tlv:"5f24,hex"`
	UsageControl               string  `tlv:"9f07,hex"`
	IssuerCountryCode          string  `tlv:"5f28,hex"`
	EffectiveDate              string  `tlv:"5f25,hex"`
	ServiceCode                string  `tlv:"5f30,hex"`
	IssuerActionCodeDenial     string  `tlv:"9f0e,hex"`
	IssuerActionCodeOnline     string  `tlv:"9f0f,hex"`
	IssuerActionCodeDefault    string  `tlv:"9f0d,hex"`
	CAPublicKeyIndex1          string  `tlv:"8f,hex"`
	IssuerPublicKeyExponent    string  `tlv:"9f32,hex"`
	IssuerPublicKeyCertificate string  `tlv:"90,hex"`
	CurrencyCode               string  `tlv:"9f42,hex"`
	CurrencyExponent           string  `tlv:"9f44,hex"`
	CDOL1                      tlv.TL  `tlv:"8c"`
	CDOL1Hex                   string  `tlv:"8c,hex"`
	CDOL2                      tlv.TL  `tlv:"8d,hex"`
	CDOL2Hex                   string  `tlv:"8d,hex"`
	VersionNumber1             string  `tlv:"9f08,hex"`
	ICCPublicKeyCertificate    string  `tlv:"9f46,hex"`
	ICCPublicKeyExponent       string  `tlv:"9f47,hex"`
	DDOL                       string  `tlv:"9f49,hex"`
	CVMList                    CVMList `tlv:"8e"`
	CVMListBytes               []byte  `tlv:"8e"`

	UnknownTag9F69 []byte  `tlv:"9f69"`
	RawTLV         tlv.TLV `tlv:"raw"`
}

type (
	CVMList struct {
		Amount       int
		SecondAmount int
		CVRules      []CVRule
	}
	CVRule struct {
		CVMCode       byte
		ConditionCode byte
	}
)

func (l *CVMList) Unmarshal(data []byte) error {
	if len(data) < 4 {
		return errors.New("failed to parse CVM list: does not contain amount")
	}
	l.Amount = int(binary.BigEndian.Uint32(data[:4]))
	data = data[4:]

	if len(data) < 4 {
		return errors.New("failed to parse CVM list: does not contain second amount")
	}
	l.SecondAmount = int(binary.BigEndian.Uint32(data[:4]))
	data = data[4:]

	for len(data) >= 2 {
		l.CVRules = append(l.CVRules, CVRule{
			CVMCode:       data[0],
			ConditionCode: data[1],
		})
		data = data[2:]
	}

	if len(data) > 0 {
		return errors.New("failed to parse CVM list: there are bytes remaining after the CV rules")
	}
	return nil
}

func (cv CVRule) FailIfUncessful() bool {
	return cv.CVMCode>>6&1 == 1
}

func (cv CVRule) Description() string {
	mode := cv.CVMCode & 0b111111
	switch mode {
	case 0b000000:
		return "Fail CVM processing"
	case 0b000001:
		return "Plaintext PIN verification performed by ICC"
	case 0b000010:
		return "Enciphered PIN verified online"
	case 0b000011:
		return "Plaintext PIN verification performed by ICC and signature"
	case 0b000100:
		return "Enciphered PIN verification performed by ICC"
	case 0b000101:
		return "Enciphered PIN verification performed by ICC and signature"
	case 0b000110:
		return "Facial biometric verified offline (by ICC)"
	case 0b000111:
		return "Facial biometric verified online"
	case 0b001000:
		return "Finger biometric verified offline (by ICC)"
	case 0b001001:
		return "Finger biometric verified online"
	case 0b001010:
		return "Palm biometric verified offline (by ICC)"
	case 0b001011:
		return "Palm biometric verified online"
	case 0b001100:
		return "Iris biometric verified offline (by ICC)"
	case 0b001101:
		return "Iris biometric verified online"
	case 0b001110:
		return "Voice biometric verified offline (by ICC)"
	case 0b001111:
		return "Voice biometric verified online"
	case 0b011110:
		return "Signature"
	case 0b011111:
		return "No CVM required"

	}
	if mode >= 0b100000 && mode <= 0b101111 {
		return "Reserved for use by the individual payment systems"
	}
	if mode >= 0b110000 && mode <= 0b111110 {
		return "Reserved for use by the issuer"
	}
	if mode >= 0b010000 && mode <= 0b011101 {
		return "Reserved for use by this specification"
	}
	return "This value is not available for use"
}

func (cv CVRule) ConditionString(amount1, amount2 int) string {
	switch cv.ConditionCode {
	case 0x00:
		return "Always"
	case 0x01:
		return "If unattended cash"
	case 0x02:
		return "If not unattended cash and not manual cash and not purchase with cashback"
	case 0x03:
		return "If terminal supports the CVM"
	case 0x04:
		return "If manual cash"
	case 0x05:
		return "If purchase with cashback"
	case 0x06:
		return fmt.Sprintf("If transaction is in the application currency and is under %d", amount1)
	case 0x07:
		return fmt.Sprintf("If transaction is in the application currency and is over %d", amount1)
	case 0x08:
		return fmt.Sprintf("If transaction is in the application currency and is under %d", amount2)
	case 0x09:
		return fmt.Sprintf("If transaction is in the application currency and is over %d", amount2)
	}
	if cv.ConditionCode >= 0x0A && cv.ConditionCode <= 0x7F {
		return "RFU"
	}
	if cv.ConditionCode >= 0x80 && cv.ConditionCode <= 0xFF {
		return "Reserved for use by individual payment systems"
	}
	return "Invalid"
}

func (et EMVProprietaryTemplate) Merge(other EMVProprietaryTemplate) EMVProprietaryTemplate {
	if len(et.Track2EquivalentData) == 0 {
		et.Track2EquivalentData = other.Track2EquivalentData
	}
	et.CardholderName = utils.CoalesceString(et.CardholderName, other.CardholderName)
	et.PAN = utils.CoalesceString(et.PAN, other.PAN)
	et.PANSequenceNumber = utils.CoalesceString(et.PANSequenceNumber, other.PANSequenceNumber)
	et.ExpirationDate = utils.CoalesceString(et.ExpirationDate, other.ExpirationDate)
	et.UsageControl = utils.CoalesceString(et.UsageControl, other.UsageControl)
	et.IssuerCountryCode = utils.CoalesceString(et.IssuerCountryCode, other.IssuerCountryCode)
	et.EffectiveDate = utils.CoalesceString(et.EffectiveDate, other.EffectiveDate)
	et.IssuerActionCodeDenial = utils.CoalesceString(et.IssuerActionCodeDenial, other.IssuerActionCodeDenial)
	et.IssuerActionCodeOnline = utils.CoalesceString(et.IssuerActionCodeOnline, other.IssuerActionCodeOnline)
	et.IssuerActionCodeDefault = utils.CoalesceString(et.IssuerActionCodeDefault, other.IssuerActionCodeDefault)
	et.CAPublicKeyIndex1 = utils.CoalesceString(et.CAPublicKeyIndex1, other.CAPublicKeyIndex1)
	et.IssuerPublicKeyExponent = utils.CoalesceString(et.IssuerPublicKeyExponent, other.IssuerPublicKeyExponent)
	et.IssuerPublicKeyCertificate = utils.CoalesceString(et.IssuerPublicKeyCertificate, other.IssuerPublicKeyCertificate)
	et.CurrencyCode = utils.CoalesceString(et.CurrencyCode, other.CurrencyCode)
	et.CurrencyExponent = utils.CoalesceString(et.CurrencyExponent, other.CurrencyExponent)
	et.CDOL1Hex = utils.CoalesceString(et.CDOL1Hex, other.CDOL1Hex)
	if len(et.CDOL1) == 0 {
		et.CDOL1 = other.CDOL1
	}
	if len(et.CDOL2) == 0 {
		et.CDOL2 = other.CDOL2
	}
	et.VersionNumber1 = utils.CoalesceString(et.VersionNumber1, other.VersionNumber1)
	et.ICCPublicKeyCertificate = utils.CoalesceString(et.ICCPublicKeyCertificate, other.ICCPublicKeyCertificate)
	et.ICCPublicKeyExponent = utils.CoalesceString(et.ICCPublicKeyExponent, other.ICCPublicKeyExponent)
	et.DDOL = utils.CoalesceString(et.DDOL, other.DDOL)
	et.CVMListBytes = utils.CoalesceBytes(et.CVMListBytes, other.CVMListBytes)
	if len(et.CVMList.CVRules) == 0 {
		et.CVMList = other.CVMList
	}
	if len(et.UnknownTag9F69) == 0 {
		et.UnknownTag9F69 = other.UnknownTag9F69
	}

	return et
}

type ApplicationTemplate struct {
	ID                []byte  `tlv:"4F"`
	Label             *string `tlv:"50"`
	PriorityIndicator *int    `tlv:"87"`
	PreferredName     *string `tlv:"9F12"`

	// RawData tlv.TLV `tlv:"raw"`
}

type GetProcessingOptionsResponse struct {
	Format1 []byte                               `tlv:"80"`
	Format2 *GetProcessingOptionsResponseFormat2 `tlv:"77"`
}

func (resp GetProcessingOptionsResponse) InterchangeProfile() AIP {
	if len(resp.Format1) >= 2 {
		return resp.Format1[:2]
	}
	if resp.Format2 != nil {
		return resp.Format2.InterchangeProfile
	}
	return nil
}

func (resp GetProcessingOptionsResponse) FileLocator() AFL {
	if len(resp.Format1) >= 2 {
		return resp.Format1[2:]
	}
	if resp.Format2 != nil {
		return resp.Format2.FileLocator
	}
	return nil
}

type GetProcessingOptionsResponseFormat2 struct {
	InterchangeProfile AIP `tlv:"82"`
	FileLocator        AFL `tlv:"94"`
}

type AIP []byte

type AFL []byte

type AFLEntry struct {
	SFI               int
	FirstRecord       int
	LastRecord        int
	RecordsInDataAuth int
}

func (afl AFL) GetEntries() ([]AFLEntry, error) {
	var results []AFLEntry
	for len(afl) > 0 {
		if len(afl) < 4 {
			return results, fmt.Errorf("AFL entry is too short (%d bytes)", len(afl))
		}
		entryBytes := afl[:4]
		afl = afl[4:]
		results = append(results, AFLEntry{
			SFI:               int(entryBytes[0]) >> 3,
			FirstRecord:       int(entryBytes[1]),
			LastRecord:        int(entryBytes[2]),
			RecordsInDataAuth: int(entryBytes[3]),
		})
	}
	return results, nil
}

func (afl AFL) GoString() string {
	entries, err := afl.GetEntries()
	if err != nil {
		return "invalid: " + err.Error()
	}
	if len(entries) == 0 {
		return "[]"
	}

	var sb strings.Builder
	sb.WriteString("[\n")
	for _, e := range entries {
		fmt.Fprintf(&sb, "  - SFI: %d\n", e.SFI)
		fmt.Fprintf(&sb, "    Records: %d-%d\n", e.FirstRecord, e.LastRecord)
		fmt.Fprintf(&sb, "    RecordsInDataAuth: %d\n", e.RecordsInDataAuth)
	}
	sb.WriteString("]")
	return sb.String()
}

type GenerateACResponse struct {
	Format1 GenerateACResponseFormat1 `tlv:"80"`
	Format2 struct {
		CryptogramInformationData     CryptogramInformationData `tlv:"9f27"`
		ApplicationTransactionCounter string                    `tlv:"9f36,hex"`
		ApplicationCryptogram         string                    `tlv:"9f26,hex"`
		IssuerApplicationData         []byte                    `tlv:"9f10,hex"`
		SignedDynamicApplicationData  string                    `tlv:"9f4b,hex"`
		Raw                           tlv.TLV                   `tlv:"raw"`
	} `tlv:"77"`
	Raw tlv.TLV `tlv:"raw"`
}

func (resp GenerateACResponse) CryptogramInformationData() CryptogramInformationData {
	if len(resp.Format1) > 0 {
		return resp.Format1.CryptogramInformationData()
	}
	return resp.Format2.CryptogramInformationData
}

func (resp GenerateACResponse) ApplicationTransactionCounter() string {
	if len(resp.Format1) > 0 {
		return resp.Format1.ApplicationTransactionCounter()
	}
	return resp.Format2.ApplicationTransactionCounter
}

func (resp GenerateACResponse) ApplicationCryptogram() string {
	if len(resp.Format2.ApplicationCryptogram) > 0 {
		return resp.Format2.ApplicationCryptogram
	}
	return resp.Format1.ApplicationCryptogram()
}

func (resp GenerateACResponse) IssuerApplicationData() []byte {
	if len(resp.Format2.IssuerApplicationData) > 0 {
		return resp.Format2.IssuerApplicationData
	}
	return resp.Format1.IssuerApplicationData()
}

type GenerateACResponseFormat1 []byte

func (f1 GenerateACResponseFormat1) CryptogramInformationData() CryptogramInformationData {
	if len(f1) < 1 {
		return 0
	}
	return CryptogramInformationData(f1[0])
}

func (f1 GenerateACResponseFormat1) ApplicationTransactionCounter() string {
	const offset = 1
	if len(f1) < offset+2 {
		return ""
	}
	return fmt.Sprintf("%02X", f1[offset:offset+2])
}

func (f1 GenerateACResponseFormat1) ApplicationCryptogram() string {
	const offset = 3
	if len(f1) < offset+8 {
		return ""
	}
	return fmt.Sprintf("%02X", f1[offset:offset+8])
}

func (f1 GenerateACResponseFormat1) IssuerApplicationData() []byte {
	const offset = 11
	if len(f1) < offset+1 {
		return nil
	}
	return f1[offset:]
}

type CryptogramInformationData byte

func (cid CryptogramInformationData) AAC() bool {
	return cid>>6 == 0b00
}

func (cid CryptogramInformationData) TC() bool {
	return cid>>6 == 0b01
}

func (cid CryptogramInformationData) ARQC() bool {
	return cid>>6 == 0b10
}

func (cid CryptogramInformationData) RFU() bool {
	return cid>>6 == 0b11
}

func (cid CryptogramInformationData) String() string {
	return fmt.Sprintf("%02X", byte(cid))
}
