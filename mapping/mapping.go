package mapping

func StrToPtr(x string) *string {
	if x == "" {
		return nil
	}
	return &x
}

func IntToPtr(x int64) *int64 {
	return &x
}

func BoolToPtr(x bool) *bool {
	return &x
}
