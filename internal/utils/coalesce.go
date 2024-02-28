package utils

func CoalesceString(a string, others ...string) string {
	if a != "" {
		return a
	}
	for _, v := range others {
		if v != "" {
			return v
		}
	}
	return ""
}

func CoalesceBytes(a []byte, others ...[]byte) []byte {
	if len(a) != 0 {
		return a
	}
	for _, v := range others {
		if len(v) != 0 {
			return v
		}
	}
	return make([]byte, 0)
}
