package slug

import "github.com/gosimple/slug"

func GenerateSlug(s string) string {
	result := slug.MakeLang(s, "en")
	if len(result) > 100 {
		return result[:100]
	}
	return result
}
