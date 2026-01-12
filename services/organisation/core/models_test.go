package core

import (
	"testing"
)

func TestIsValidDocType(t *testing.T) {
	tests := []struct {
		name    string
		docType OrgDocType
		valid   bool
	}{
		{name: "phys passport photo", docType: DocPhysPassportPhoto, valid: true},
		{name: "phys passport propiska", docType: DocPhysPassportPropiska, valid: true},
		{name: "phys uchet", docType: DocPhysUchet, valid: true},
		{name: "jur reg svid", docType: DocJurRegSvid, valid: true},
		{name: "jur uchet", docType: DocJurUchet, valid: true},
		{name: "jur appointment protocol", docType: DocJurAppointmentProtocol, valid: true},
		{name: "jur usn", docType: DocJurUSN, valid: true},
		{name: "jur ustav", docType: DocJurUstav, valid: true},
		{name: "ip uchet", docType: DocIPUchet, valid: true},
		{name: "ip passport photo", docType: DocIPPassportPhoto, valid: true},
		{name: "ip passport propiska", docType: DocIPPassportPropiska, valid: true},
		{name: "ip usn", docType: DocIPUSN, valid: true},
		{name: "ip ogrnip", docType: DocIPOGRNIP, valid: true},
		{name: "invalid doc type", docType: OrgDocType("invalid"), valid: false},
		{name: "empty doc type", docType: OrgDocType(""), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidDocType(tt.docType)
			if result != tt.valid {
				t.Errorf("IsValidDocType() = %v, want %v", result, tt.valid)
			}
		})
	}
}

func TestIsValidForOrgType(t *testing.T) {
	tests := []struct {
		name    string
		docType OrgDocType
		orgType OrgType
		valid   bool
	}{
		{name: "phys passport photo for phys org", docType: DocPhysPassportPhoto, orgType: OrgTypePhys, valid: true},
		{name: "phys passport propiska for phys org", docType: DocPhysPassportPropiska, orgType: OrgTypePhys, valid: true},
		{name: "phys uchet for phys org", docType: DocPhysUchet, orgType: OrgTypePhys, valid: true},
		{name: "phys doc for jur org", docType: DocPhysPassportPhoto, orgType: OrgTypeJur, valid: false},
		{name: "phys doc for ip org", docType: DocPhysPassportPhoto, orgType: OrgTypeIP, valid: false},

		{name: "jur reg svid for jur org", docType: DocJurRegSvid, orgType: OrgTypeJur, valid: true},
		{name: "jur uchet for jur org", docType: DocJurUchet, orgType: OrgTypeJur, valid: true},
		{name: "jur appointment protocol for jur org", docType: DocJurAppointmentProtocol, orgType: OrgTypeJur, valid: true},
		{name: "jur usn for jur org", docType: DocJurUSN, orgType: OrgTypeJur, valid: true},
		{name: "jur ustav for jur org", docType: DocJurUstav, orgType: OrgTypeJur, valid: true},
		{name: "jur doc for phys org", docType: DocJurRegSvid, orgType: OrgTypePhys, valid: false},
		{name: "jur doc for ip org", docType: DocJurRegSvid, orgType: OrgTypeIP, valid: false},

		{name: "ip uchet for ip org", docType: DocIPUchet, orgType: OrgTypeIP, valid: true},
		{name: "ip passport photo for ip org", docType: DocIPPassportPhoto, orgType: OrgTypeIP, valid: true},
		{name: "ip passport propiska for ip org", docType: DocIPPassportPropiska, orgType: OrgTypeIP, valid: true},
		{name: "ip usn for ip org", docType: DocIPUSN, orgType: OrgTypeIP, valid: true},
		{name: "ip ogrnip for ip org", docType: DocIPOGRNIP, orgType: OrgTypeIP, valid: true},
		{name: "ip doc for phys org", docType: DocIPUchet, orgType: OrgTypePhys, valid: false},
		{name: "ip doc for jur org", docType: DocIPUchet, orgType: OrgTypeJur, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.docType.IsValidForOrgType(tt.orgType)
			if result != tt.valid {
				t.Errorf("IsValidForOrgType(%s) = %v, want %v", tt.orgType, result, tt.valid)
			}
		})
	}
}

func TestOrgTypeValidation(t *testing.T) {
	tests := []struct {
		name    string
		orgType OrgType
		valid   bool
	}{
		{name: "valid phys", orgType: OrgTypePhys, valid: true},
		{name: "valid jur", orgType: OrgTypeJur, valid: true},
		{name: "valid ip", orgType: OrgTypeIP, valid: true},
		{name: "invalid org type", orgType: OrgType("invalid"), valid: false},
		{name: "empty org type", orgType: OrgType(""), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.orgType == OrgTypePhys || tt.orgType == OrgTypeJur || tt.orgType == OrgTypeIP
			if valid != tt.valid {
				t.Errorf("OrgType validation failed for %s", tt.orgType)
			}
		})
	}
}
