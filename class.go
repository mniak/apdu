package apdu

type Class byte

func (c Class) Invalid() bool {
	return c == 0xFF
}

// func (c Class) FirstInterindustry() bool {
// 	return c&0b1110_0000 == 0b0000_0000
// }

// func (c Class) SecondInterindustry() bool {
// 	return c&0b1000_0000 != 0
// }

func (c Class) Proprietary() bool {
	return c&0b1000_0000 != 0
}

func (c Class) LastInChain() bool {
	return !c.Proprietary() && c&0b0001_0000 == 0
}

// func (c Class) SecureMessagingIndication() SecureMessagingIndication {
// 	if c.Proprietary() {
// 		return NoSMOrNoIndication
// 	}
// 	return NoSMOrNoIndication
// }

type SecureMessagingIndication string

const (
	NoSMOrNoIndication                  = "No SM or no indication"
	ProprietarySMFormat                 = "Proprietary SM format"
	SMAccordingToSection6_NotProcessed  = "SM according to section 6, command header not processed according to 6.2.3.1"
	SMAccordingToSection6_Authenticated = "SM according to section 6, command header authenticated according to 6.2.3.1"
)
