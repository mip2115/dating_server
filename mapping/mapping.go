package mapping

import "time"

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

func IntToPtr(x int) *int {
	return &x
}

func IntToV(x *int) int {
	if x == nil {
		return 0
	}
	return *x
}

func Int64ToPtr(x int64) *int64 {
	return &x
}

func Int64ToV(x *int64) int64 {
	if x == nil {
		return 0
	}
	return *x
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

func TimeToPtr(x time.Time) *time.Time {
	return &x
}

func TimeToV(x *time.Time) time.Time {
	if x == nil {
		return time.Time{}
	}
	return *x
}
