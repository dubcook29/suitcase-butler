package whois

import "github.com/suitcase/butler/wmpci"

var WMPApplicationBasic = wmpci.WMPBasic{
	Id:      "butler_wmp_whois",
	Name:    "butler butilt wmp whois",
	Version: "0.0.1(Alpha)",
	Copyright: []wmpci.Author{
		{AuthorName: "dubcook29"},
	},
}

var WMPApplicationRequestDefault = wmpci.WMPRequest{
	"domain": {"example.com"},
	"asn":    {"AS15133"},
	"ip":     {"23.192.228.84"},
}

var WMPApplicationResponseDefault = wmpci.WMPResponse{
	"domain_whois":     {""},
	"domain_whois_raw": {""},
	"asn_whois":        {""},
	"asn_whois_raw":    {""},
	"ip_whois":         {""},
	"ip_whois_raw":     {""},
}

var WMPApplicationConfigDefault = map[string]wmpci.WMPCustom{
	"resend": {
		Name:        "resend",
		Value:       3,
		Description: "Set the number of retries for CDN Check query requests",
	},
}
