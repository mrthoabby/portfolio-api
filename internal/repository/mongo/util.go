package mongo

import (
	"net/url"
	"strings"
)

func sanitizeMongoURLForLogging(dbURL string) string {
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		if strings.Contains(dbURL, "@") {
			parts := strings.SplitN(dbURL, "@", 2)
			if len(parts) == 2 {
				return "****:****@" + parts[1]
			}
		}
		return dbURL
	}

	if parsedURL.User != nil {
		parsedURL.User = url.UserPassword("****", "****")
	}

	return parsedURL.String()
}
