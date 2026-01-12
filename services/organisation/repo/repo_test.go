package repo

import (
	"testing"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
)

func TestDocTypeValidation(t *testing.T) {
	tests := []struct {
		name    string
		docType core.OrgDocType
		orgType core.OrgType
		valid   bool
	}{
		{name: "phys passport photo for phys org", docType: core.DocPhysPassportPhoto, orgType: core.OrgTypePhys, valid: true},
		{name: "phys passport propiska for phys org", docType: core.DocPhysPassportPropiska, orgType: core.OrgTypePhys, valid: true},
		{name: "phys uchet for phys org", docType: core.DocPhysUchet, orgType: core.OrgTypePhys, valid: true},
		{name: "phys doc with jur org", docType: core.DocPhysPassportPhoto, orgType: core.OrgTypeJur, valid: false},
		{name: "phys doc with ip org", docType: core.DocPhysPassportPhoto, orgType: core.OrgTypeIP, valid: false},
		{name: "jur reg svid for jur org", docType: core.DocJurRegSvid, orgType: core.OrgTypeJur, valid: true},
		{name: "jur uchet for jur org", docType: core.DocJurUchet, orgType: core.OrgTypeJur, valid: true},
		{name: "jur doc with phys org", docType: core.DocJurRegSvid, orgType: core.OrgTypePhys, valid: false},
		{name: "jur doc with ip org", docType: core.DocJurRegSvid, orgType: core.OrgTypeIP, valid: false},
		{name: "ip uchet for ip org", docType: core.DocIPUchet, orgType: core.OrgTypeIP, valid: true},
		{name: "ip passport photo for ip org", docType: core.DocIPPassportPhoto, orgType: core.OrgTypeIP, valid: true},
		{name: "ip doc with phys org", docType: core.DocIPUchet, orgType: core.OrgTypePhys, valid: false},
		{name: "ip doc with jur org", docType: core.DocIPUchet, orgType: core.OrgTypeJur, valid: false},
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
