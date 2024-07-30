package utils

import "regexp"

var (
	userIDRegex = regexp.MustCompile(`/user/\d+`)
	skuIDRegex  = regexp.MustCompile(`/cart/\d+`)
)

func CleanURL(path string) string {
	cleanedPath := userIDRegex.ReplaceAllString(path, "/user/{user_id}")
	cleanedPath = skuIDRegex.ReplaceAllString(cleanedPath, "/cart/{sku_id}")
	return cleanedPath
}
