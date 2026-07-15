package auth

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
)

var publicURLs = []*regexp.Regexp{
	regexp.MustCompile("^/.well-known"),
	regexp.MustCompile("^/health$"),
	regexp.MustCompile("^/docs"),
	// INFO: This is needed to allow new user creation via invite link
	regexp.MustCompile("^/api/v1/users$"),
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
	Validator:      APIKeyValidator,
	Realm:          "Restricted",
	Extractor:      extractors.FromAuthHeader("Bearer"),
}
