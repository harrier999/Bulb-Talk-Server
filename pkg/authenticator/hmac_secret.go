package authenticator

import "os"

func GetHMACSecret() string {
	return os.Getenv("JWT_SECRET")
}
