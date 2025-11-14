package wmpdata

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/suitcase/butler/wmpci/data/model"
)

type IPWhois struct {
	model.WMPDataModelBasicStructure
	CommonWhoisRegisterDataModel

	IpAddress string `json:"ip_address,omitempty" bson:"ip_address"`
	ASNumber  string `json:"as_number,omitempty" bson:"as_number"`
	Cidr      string `json:"cidr,omitempty" bson:"cidr"`
	IpRange   string `json:"ip_range,omitempty" bson:"ip_range"`
}

func (i IPWhois) Key() string {
	return "ip_whois"
}

func (i IPWhois) Exchange(data map[string][]interface{}) map[string][]interface{} {
	if val, ok := data[i.Key()]; ok {
		data[i.Key()] = append(val, i)
	} else {
		data[i.Key()] = []interface{}{i}
	}
	return data
}

func (i IPWhois) DataModel() model.WMPDataModelBasicStructure {
	return i.WMPDataModelBasicStructure
}

func (i IPWhois) JSON() (json.RawMessage, error) {
	return json.Marshal(&i)
}

func (i IPWhois) New(assetId string) IPWhois {
	i.WMPDataModelBasicStructure.CreatedAt = time.Now()
	i.WMPDataModelBasicStructure.UpdatedAt = time.Now()
	i.WMPDataModelBasicStructure.AssetId = assetId
	i.WMPDataModelBasicStructure.Identifier = uuid.NewString()
	i.WMPDataModelBasicStructure.DataModelKey = i.Key()
	return i
}

type ASNWhois struct {
	model.WMPDataModelBasicStructure
	CommonWhoisRegisterDataModel

	ASNumber string `json:"as_number" bson:"as_number"`
}

func (a ASNWhois) Key() string {
	return "as_whois"
}

func (a ASNWhois) Exchange(data map[string][]interface{}) map[string][]interface{} {
	if val, ok := data[a.Key()]; ok {
		data[a.Key()] = append(val, a)
	} else {
		data[a.Key()] = []interface{}{a}
	}
	return data
}

func (a ASNWhois) DataModel() model.WMPDataModelBasicStructure {
	return a.WMPDataModelBasicStructure
}

func (a ASNWhois) JSON() (json.RawMessage, error) {
	return json.Marshal(&a)
}

func (a ASNWhois) New(assetId string) ASNWhois {
	a.WMPDataModelBasicStructure.CreatedAt = time.Now()
	a.WMPDataModelBasicStructure.UpdatedAt = time.Now()
	a.WMPDataModelBasicStructure.AssetId = assetId
	a.WMPDataModelBasicStructure.Identifier = uuid.NewString()
	a.WMPDataModelBasicStructure.DataModelKey = a.Key()
	return a
}

type CommonWhoisRegisterDataModel struct {
	RegistrarOrgName    string `json:"registrar_org_name,omitempty" bson:"registrar_org_name"`
	RegistrarOrgEmail   string `json:"registrar_org_email,omitempty" bson:"registrar_org_email"`
	RegistrarOrgPhone   string `json:"registrar_org_phone,omitempty" bson:"registrar_org_phone"`
	RegistrarOrgAddress string `json:"registrar_org_address,omitempty" bson:"registrar_org_address"`
	RegistrarAbuseName  string `json:"registrar_abuse_name,omitempty" bson:"registrar_abuse_name"`
	RegistrarAbuseEmail string `json:"registrar_abuse_email,omitempty" bson:"registrar_abuse_email"`
	RegistrarAbusePhone string `json:"registrar_abuse_phone,omitempty" bson:"registrar_abuse_phone"`

	RegistrantOrgName    string `json:"registrant_org_name,omitempty" bson:"registrant_org_name"`
	RegistrantOrgEmail   string `json:"registrant_org_email,omitempty" bson:"registrant_org_email"`
	RegistrantOrgPhone   string `json:"registrant_org_phone,omitempty" bson:"registrant_org_phone"`
	RegistrantOrgAddress string `json:"registrant_org_address,omitempty" bson:"registrant_org_address"`
	RegistrantAbuseName  string `json:"registrant_abuse_name,omitempty" bson:"registrant_abuse_name"`
	RegistrantAbuseEmail string `json:"registrant_abuse_email,omitempty" bson:"registrant_abuse_email"`
	RegistrantAbusePhone string `json:"registrant_abuse_phone,omitempty" bson:"registrant_abuse_phone"`

	TechnologyName        string `json:"technology_name,omitempty" bson:"technology_name"`
	TechnologyEmail       string `json:"technology_email,omitempty" bson:"technology_email"`
	TechnologyPhone       string `json:"technology_phone,omitempty" bson:"technology_phone"`
	TechnologyTitle       string `json:"technology_title,omitempty" bson:"technology_title"`
	TechnologyReferralURL string `json:"technology_referral_url,omitempty" bson:"technology_referral_url"`

	AdministratorName        string `json:"administrator_name,omitempty" bson:"administrator_name"`
	AdministratorEmail       string `json:"administrator_email,omitempty" bson:"administrator_email"`
	AdministratorPhone       string `json:"administrator_phone,omitempty" bson:"administrator_phone"`
	AdministratorTitle       string `json:"administrator_title,omitempty" bson:"administrator_title"`
	AdministratorReferralURL string `json:"administrator_referral_url,omitempty" bson:"administrator_referral_url"`
}

type WhoisRegisterData struct {
	Name         string `json:"name,omitempty" bson:"name"`
	Organization string `json:"organization,omitempty" bson:"organization"`
	Street       string `json:"street,omitempty" bson:"street"`
	City         string `json:"city,omitempty" bson:"city"`
	Province     string `json:"province,omitempty" bson:"province"`
	PostalCode   string `json:"postal_code,omitempty" bson:"postal_code"`
	Country      string `json:"country,omitempty" bson:"country"`
	Phone        string `json:"phone,omitempty" bson:"phone"`
	PhoneExt     string `json:"phone_ext,omitempty" bson:"phone_ext"`
	Fax          string `json:"fax,omitempty" bson:"fax"`
	FaxExt       string `json:"fax_ext,omitempty" bson:"fax_ext"`
	Email        string `json:"email,omitempty" bson:"email"`
	ReferralURL  string `json:"referral_url,omitempty" bson:"referral_url"`
}

type DomainWhois struct {
	model.WMPDataModelBasicStructure

	ID             string   `json:"id,omitempty"`
	Domain         string   `json:"domain,omitempty"`
	Punycode       string   `json:"punycode,omitempty"`
	Name           string   `json:"name,omitempty"`
	Extension      string   `json:"extension,omitempty"`
	WhoisServer    string   `json:"whois_server,omitempty"`
	Status         []string `json:"status,omitempty"`
	NameServers    []string `json:"name_servers,omitempty"`
	DNSSec         bool     `json:"dnssec,omitempty"`
	CreatedDate    string   `json:"created_date,omitempty"`
	UpdatedDate    string   `json:"updated_date,omitempty"`
	ExpirationDate string   `json:"expiration_date,omitempty"`

	WhoisBasicInfo
}

func (d DomainWhois) Key() string {
	return "domain_whois"
}

func (d DomainWhois) Exchange(data map[string][]interface{}) map[string][]interface{} {
	if val, ok := data[d.Key()]; ok {
		data[d.Key()] = append(val, d)
	} else {
		data[d.Key()] = []interface{}{d}
	}
	return data
}

func (d DomainWhois) DataModel() model.WMPDataModelBasicStructure {
	return d.WMPDataModelBasicStructure
}

func (d DomainWhois) JSON() (json.RawMessage, error) {
	return json.Marshal(&d)
}

func (d DomainWhois) New(assetId string) DomainWhois {
	d.WMPDataModelBasicStructure.CreatedAt = time.Now()
	d.WMPDataModelBasicStructure.UpdatedAt = time.Now()
	d.WMPDataModelBasicStructure.AssetId = assetId
	d.WMPDataModelBasicStructure.Identifier = uuid.NewString()
	d.WMPDataModelBasicStructure.DataModelKey = d.Key()
	return d
}

type WhoisBasicInfo struct {
	Registrar      Contact `json:"registrar,omitempty"`
	Registrant     Contact `json:"registrant,omitempty"`
	Administrative Contact `json:"administrative,omitempty"`
	Technical      Contact `json:"technical,omitempty"`
	Billing        Contact `json:"billing,omitempty"`
}

type Contact struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Organization string `json:"organization,omitempty"`
	Street       string `json:"street,omitempty"`
	City         string `json:"city,omitempty"`
	Province     string `json:"province,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	Country      string `json:"country,omitempty"`
	Phone        string `json:"phone,omitempty"`
	PhoneExt     string `json:"phone_ext,omitempty"`
	Fax          string `json:"fax,omitempty"`
	FaxExt       string `json:"fax_ext,omitempty"`
	Email        string `json:"email,omitempty"`
	ReferralURL  string `json:"referral_url,omitempty"`
}
