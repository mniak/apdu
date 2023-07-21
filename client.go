package apdu

type Client struct {
	RawClient
	LowLevelCommands LowLevelCommands
	HighLevelCommands
}

func NewClient(driver Driver) Client {
	raw := NewRawClient(driver)
	low := _LowLevelClient{
		RawClient: raw,
	}
	high := _HighLevelClient{
		Low: low,
	}
	return Client{
		RawClient:         raw,
		LowLevelCommands:  low,
		HighLevelCommands: high,
	}
}
