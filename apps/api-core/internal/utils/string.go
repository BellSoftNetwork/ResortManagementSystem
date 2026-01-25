package utils

// StringPtrToString converts *string to string, returning empty string if nil
func StringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
