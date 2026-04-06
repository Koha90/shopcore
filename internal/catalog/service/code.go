package service

import (
	"strings"
	"unicode"
)

// SuggestCategoryCode builds a predictable category code from category name.
//
// The function uses a small built-in Cyrillic transliteration table and
// converts separators into dash-delimited ASCII slug form.
func SuggestCategoryCode(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(name))

	prevDash := false

	for _, r := range strings.ToLower(name) {
		switch {
		case isASCIIAlphaNum(r):
			b.WriteRune(r)
			prevDash = false

		case isCyrillic(r):
			s := translitRune(r)
			if s == "" {
				if !prevDash && b.Len() > 0 {
					b.WriteByte('-')
					prevDash = true
				}
				continue
			}

			for _, tr := range s {
				if isASCIIAlphaNum(tr) {
					b.WriteRune(tr)
					prevDash = false
					continue
				}

				if !prevDash && b.Len() > 0 {
					b.WriteByte('-')
					prevDash = true
				}
			}

		case isSeparator(r):
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}

	code := strings.Trim(b.String(), "-")
	code = collapseDashes(code)

	return code
}

func isASCIIAlphaNum(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}

func isCyrillic(r rune) bool {
	return unicode.In(r, unicode.Cyrillic)
}

func isSeparator(r rune) bool {
	switch {
	case unicode.IsSpace(r):
		return true
	case unicode.IsPunct(r):
		return true
	case unicode.IsSymbol(r):
		return true
	default:
		return false
	}
}

func collapseDashes(s string) string {
	if s == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(s))

	prevDash := false
	for _, r := range s {
		if r == '-' {
			if prevDash {
				continue
			}
			prevDash = true
			b.WriteRune(r)
			continue
		}

		prevDash = false
		b.WriteRune(r)
	}

	return strings.Trim(b.String(), "-")
}

func translitRune(r rune) string {
	switch r {
	case 'а':
		return "a"
	case 'б':
		return "b"
	case 'в':
		return "v"
	case 'г':
		return "g"
	case 'д':
		return "d"
	case 'е':
		return "e"
	case 'ё':
		return "e"
	case 'ж':
		return "zh"
	case 'з':
		return "z"
	case 'и':
		return "i"
	case 'й':
		return "y"
	case 'к':
		return "k"
	case 'л':
		return "l"
	case 'м':
		return "m"
	case 'н':
		return "n"
	case 'о':
		return "o"
	case 'п':
		return "p"
	case 'р':
		return "r"
	case 'с':
		return "s"
	case 'т':
		return "t"
	case 'у':
		return "u"
	case 'ф':
		return "f"
	case 'х':
		return "h"
	case 'ц':
		return "ts"
	case 'ч':
		return "ch"
	case 'ш':
		return "sh"
	case 'щ':
		return "sch"
	case 'ъ':
		return ""
	case 'ы':
		return "y"
	case 'ь':
		return ""
	case 'э':
		return "e"
	case 'ю':
		return "yu"
	case 'я':
		return "ya"
	default:
		return ""
	}
}
