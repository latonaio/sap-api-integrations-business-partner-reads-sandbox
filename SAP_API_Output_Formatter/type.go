package sap_api_output_formatter

type BusinessPartner struct {
	ConnectionKey   string `json:"connection_key"`
	Result          bool   `json:"result"`
	RedisKey        string `json:"redis_key"`
	Filepath        string `json:"filepath"`
	APISchema       string `json:"api_schema"`
	BusinessPartner string `json:"business_partner_code"`
	Deleted         bool   `json:"deleted"`
}

type Role struct {
	BusinessPartner     string `json:"BusinessPartner"`
	BusinessPartnerRole string `json:"BusinessPartnerRole"`
	ValidFrom           string `json:"ValidFrom"`
	ValidTo             string `json:"ValidTo"`
}

type Address struct {
	BusinessPartner   string `json:"BusinessPartner"`
	AddressID         string `json:"AddressID"`
	ValidityStartDate string `json:"ValidityStartDate"`
	ValidityEndDate   string `json:"ValidityEndDate"`
	Country           string `json:"Country"`
	Region            string `json:"Region"`
	StreetName        string `json:"StreetName"`
	CityName          string `json:"CityName"`
	PostalCode        string `json:"PostalCode"`
	Language          string `json:"Language"`
}
