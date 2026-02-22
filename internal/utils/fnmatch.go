package utils

func FnMatch(pattern string, str string) bool {

	if len(pattern) == 0 {
		return len(str) == 0
	}

	if pattern[0] == '*' {
		for i := 0; i <= len(str); i++ {
			if FnMatch(pattern[1:], str[i:]) {
				return true
			}
		}
		return false
	}

	if len(str) == 0 {
		return len(pattern) == 0
	}

	if pattern[0] == '?' {
		return FnMatch(pattern[1:], str[1:])
	}

	if pattern[0] != str[0] {
		return false
	}

	return FnMatch(pattern[1:], str[1:])
}
