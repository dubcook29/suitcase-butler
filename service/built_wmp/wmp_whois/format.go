package wmp_whois

import (
	"encoding/json"

	whoisparser "github.com/likexian/whois-parser"
	"github.com/suitcase/butler/wmpci/data/wmpdata"
)

func (wmp WMPWhois) domainWhoisFormatFunc(w string) (wmpdata.DomainWhois, error) {
	if data, err := whoisparser.Parse(w); err != nil {
		return wmpdata.DomainWhois{}, err // fix a bug
	} else {
		return whoisparserToWMPDataModelDomainWhois(data, wmp.assetId), nil
	}
}

func whoisparserToWMPDataModelDomainWhois(whoisInfo whoisparser.WhoisInfo, assetId string) wmpdata.DomainWhois {
	var result wmpdata.DomainWhois
	result = result.New(assetId)

	jumpContact(whoisInfo.Registrar, &result.Registrar)
	jumpContact(whoisInfo.Registrant, &result.Registrant)
	jumpContact(whoisInfo.Administrative, &result.Administrative)
	jumpContact(whoisInfo.Technical, &result.Technical)
	jumpContact(whoisInfo.Billing, &result.Billing)
	jumpDomain(whoisInfo.Domain, &result)

	return result
}

func jumpContact(s *whoisparser.Contact, d *wmpdata.Contact) {
	data, _ := json.Marshal(s)
	json.Unmarshal(data, d)
}

func jumpDomain(s *whoisparser.Domain, d *wmpdata.DomainWhois) {
	data, _ := json.Marshal(s)
	json.Unmarshal(data, d)
}

func ASWhoisFormatFunc(w string) {

}

func IPWhoisFormatFunc(w string) {

}
