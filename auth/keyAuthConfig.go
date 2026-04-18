package auth

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
)

var publicURLs = []*regexp.Regexp{
	regexp.MustCompile("^/health$"),
}

var KeyAuthConfig keyauth.Config = keyauth.Config{
	Next: func(c fiber.Ctx) bool {
		path := strings.ToLower(c.Path())
		for _, pattern := range publicURLs {
			if pattern.MatchString(path) {
				return true
			}
		}

		return false
	},
	SuccessHandler: nil,
	ErrorHandler:   nil,
	Validator:      ApiKeyValidator,
	Realm:          "Restricted",
	Extractor:      extractors.FromAuthHeader("Bearer"),
}
