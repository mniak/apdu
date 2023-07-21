package apdu

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
	PDOL                             []byte            `tlv:"9f38"`
	IssuerDiscretionaryData          map[string][]byte `tlv:"bf0c"`
}

type RecordTemplate struct {
	EMVProprietaryTemplates []EMVProprietaryTemplate `tlv:"70"`

	// RawData tlv.TLV `tlv:"raw"`
}

type EMVProprietaryTemplate struct {
	ApplicationTemplates []ApplicationTemplate `tlv:"61"`

	Track2EquivalentData       *string `tlv:"57,hex"`
	CardholderName             *string `tlv:"5f20"`
	PAN                        string  `tlv:"5a,hex"`
	PANSequenceNumber          string  `tlv:"5f34,hex"`
	ExpirationDate             string  `tlv:"5f24,hex"`
	UsageControl               string  `tlv:"9f07,hex"`
	IssuerCountryCode          string  `tlv:"5f28,hex"`
	EffectiveDate              string  `tlv:"5f25,hex"`
	IssuerActionCodeDenial     string  `tlv:"9f0e,hex"`
	IssuerActionCodeOnline     string  `tlv:"9f0f,hex"`
	IssuerActionCodeDefault    string  `tlv:"9f0d,hex"`
	CAPublicKeyIndex1          string  `tlv:"8f,hex"`
	IssuerPublicKeyExponent    string  `tlv:"9f32,hex"`
	IssuerPublicKeyCertificate string  `tlv:"90,hex"`
	CurrencyCode               string  `tlv:"9f42,hex"`
	CurrencyExponent           string  `tlv:"9f44,hex"`
	CDOL1                      string  `tlv:"8c,hex"`
	CDOL2                      string  `tlv:"8d,hex"`
	VersionNumber1             string  `tlv:"9f08,hex"`
	ICCPublicKeyCertificate    string  `tlv:"9f46,hex"`
	ICCPublicKeyExponent       string  `tlv:"9f47,hex"`
	DDOL                       string  `tlv:"9f49,hex"`
	CVMList                    string  `tlv:"8e,hex"`

	UnknownTag9F69 []byte `tlv:"9f69"`
	// RawData        tlv.TLV `tlv:"raw"`
}

type ApplicationTemplate struct {
	ID                []byte  `tlv:"4F"`
	Label             *string `tlv:"50"`
	PriorityIndicator *int    `tlv:"87"`
	PreferredName     *string `tlv:"9F12"`

	// RawData tlv.TLV `tlv:"raw"`
}
