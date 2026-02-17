package destiny

type MembershipResponse struct {
	Response struct {
		DestinyMemberships  []DestinyMembership `json:"destinyMemberships"`
		PrimaryMembershipId string              `json:"primaryMembershipId"`
	} `json:"Response"`
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

type DestinyMembership struct {
	MembershipType          int    `json:"membershipType"`
	MembershipId            string `json:"membershipId"`
	DisplayName             string `json:"displayName"`
	BungieGlobalDisplayName string `json:"bungieGlobalDisplayName"`
}
