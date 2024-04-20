package manipulation

func ParseBool(value string) *bool {
	if value == "true" {
		result := true
		return &result
	} else if value == "false" {
		result := false
		return &result
	}

	return nil
}
