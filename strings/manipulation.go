package strings

import (
	"regexp"
	"strings"

	"github.com/fiam/gounidecode/unidecode"
)

var fileSanitizeRegex *regexp.Regexp
var fileSanDupUnders *regexp.Regexp

func init() {
	fileSanitizeRegex, _ = regexp.Compile(`([^\w\.\_\-])`)
	fileSanDupUnders, _ = regexp.Compile(`(\_){2,}`)
}

func ToASCII(s string) string {
	var re *regexp.Regexp

	// transliteration
	s = unidecode.Unidecode(s)

	// replace whitespace and some punctuation with underscores
	re = regexp.MustCompile("[\\s@_/\\\\-]+")
	s = re.ReplaceAllString(s, "_")

	// remove invalid characters
	re = regexp.MustCompile("[^\\w.]+")
	s = re.ReplaceAllString(s, "")

	// remove duplicated underscores
	re = regexp.MustCompile("__+")
	s = re.ReplaceAllString(s, "_")

	// trim
	s = strings.Trim(s, "_ ")

	return s
}

func ToSlug(s string) string {
	var re *regexp.Regexp

	// transliteration
	s = unidecode.Unidecode(s)

	// convert to lowercase
	s = strings.ToLower(s)

	// replace whitespace and some punctuation with underscores
	re = regexp.MustCompile("[\\s/\\\\-]+")
	s = re.ReplaceAllString(s, "_")

	// remove invalid characters
	re = regexp.MustCompile("[^\\w]+")
	s = re.ReplaceAllString(s, "")

	// remove duplicated underscores
	re = regexp.MustCompile("__+")
	s = re.ReplaceAllString(s, "_")

	// trim
	s = strings.Trim(s, "_ ")

	return s
}

func SanitizeFilename(s string) (newS string) {
	s = unidecode.Unidecode(s)
	newS = fileSanitizeRegex.ReplaceAllString(s, "_")
	newS = fileSanDupUnders.ReplaceAllString(newS, "_")
	return
}
func SanitizeFolderPath(s string) (newS string) {

	for i, u := range strings.Split(s, "/") {
		u = SanitizeFilename(u)
		if i == 0 {
			newS = u
		} else {
			newS = newS + "/" + u
		}
	}

	return
}

func DigitsOnly(s string) string {
	re := regexp.MustCompile("[^0-9]+")
	s = re.ReplaceAllString(s, "")
	return s
}
