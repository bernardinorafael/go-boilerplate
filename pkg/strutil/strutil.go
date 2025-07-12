package strutil

import (
	"fmt"
	"regexp"
	"strings"
)

func Slugify(val string) string {
	if len(val) == 0 {
		return ""
	}

	val = strings.ToLower(val)

	var sb strings.Builder
	for _, char := range val {
		if char == ' ' {
			sb.WriteRune(char)
			continue
		}
		switch char {
		case 'á', 'à', 'ã', 'â', 'ä':
			sb.WriteRune('a')
		case 'é', 'è', 'ê', 'ë':
			sb.WriteRune('e')
		case 'í', 'ì', 'î', 'ï':
			sb.WriteRune('i')
		case 'ó', 'ò', 'õ', 'ô', 'ö':
			sb.WriteRune('o')
		case 'ú', 'ù', 'û', 'ü':
			sb.WriteRune('u')
		case 'ç':
			sb.WriteRune('c')
		default:
			sb.WriteRune(char)
		}
	}
	val = sb.String()

	r := regexp.MustCompile(`[^a-z0-9]+`)
	val = r.ReplaceAllString(val, "-")

	r, _ = regexp.Compile("-+")
	val = r.ReplaceAllString(val, "-")

	val = strings.Trim(val, "-")

	return val
}

func SlugifyWithPrefix(prefix, val string) string {
	if len(val) == 0 {
		return ""
	}

	return fmt.Sprintf("%s-%s", prefix, Slugify(val))
}

func IsSlug(val string) bool {
	return regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`).MatchString(val)
}
