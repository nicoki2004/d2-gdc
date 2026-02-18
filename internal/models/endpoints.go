package models

import "fmt"

const (
	BASEURL               = "https://www.bungie.net/Platform"
	AUTH_URL_PREFIX       = "https://www.bungie.net/en/OAuth/Authorize?"
	AUTH_TOKEN_URL_PREFIX = "https://www.bungie.net/Platform/App/OAuth/Token/"
	URL_MEMBERSHIP_USER   = "https://www.bungie.net/Platform/User/GetMembershipsForCurrentUser/"
)

func GetProfileURL(membershipType int, membershipId string) string {
	return fmt.Sprintf("%s/Destiny2/%d/Profile/%s/", BASEURL, membershipType, membershipId)
}

func GetCharacterURL(membershipType int, membershipId string, characterId string) string {
	return fmt.Sprintf("%s/Destiny2/%d/Profile/%s/Character/%s/", BASEURL, membershipType, membershipId, characterId)
}
