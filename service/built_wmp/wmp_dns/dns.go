package wmp_dns

import (
	"github.com/miekg/dns"
	"github.com/projectdiscovery/retryabledns"
)

func NSLookup(hostname string, retries int, resolvers []string) (*retryabledns.DNSData, error) {

	dnsClient, err := retryabledns.New(resolvers, retries)
	if err != nil {
		return nil, err
	}

	// Query Types: dns.TypeA, dns.TypeNS, dns.TypeCNAME, dns.TypeSOA, dns.TypePTR, dns.TypeMX, dns.TypeANY
	// dns.TypeTXT, dns.TypeAAAA, dns.TypeSRV (from github.com/miekg/dns)
	requestTypes := []uint16{
		dns.TypeA,
		dns.TypeMX,
		dns.TypeNS,
		dns.TypeTXT,
		dns.TypePTR,
		dns.TypeAAAA,
		dns.TypeCNAME,
	}
	dnsResponses, err := dnsClient.QueryMultiple(hostname, requestTypes)
	if err != nil {
		return nil, err
	}

	return dnsResponses, nil
}
