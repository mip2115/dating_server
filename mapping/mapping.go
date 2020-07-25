package mapping

func StrToPtr(x string) *string {
	if x == "" {
		return nil
	}
	return &x
}

func StrToV(x *string) string {
	if x == nil {
		return ""
	}
	return *x
}

func IntToPtr(x int64) *int64 {
	return &x
}

func BoolToPtr(x bool) *bool {
	return &x
}

func BoolToV(x *bool) bool {
	if x == nil {
		return false
	}
	return *x
}
