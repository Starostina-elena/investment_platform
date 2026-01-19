package core

import "time"

type OrgType string

const (
	OrgTypeJur  OrgType = "jur"
	OrgTypePhys OrgType = "phys"
	OrgTypeIP   OrgType = "ip"
)

type OrgPermission string

const (
	OrgAccountManagement OrgPermission = "org_account_management"
	MoneyManagement      OrgPermission = "money_management"
	ProjectManagement    OrgPermission = "project_management"
)

type OrgEmployee struct {
	OrgID      int    `json:"org_id" db:"org_id"`
	UserID     int    `json:"user_id" db:"user_id"`
	UserEmail  string `json:"user_email" db:"user_email"`
	UserName   string `json:"nickname" db:"nickname"`
	OrgAccMgmt bool   `json:"org_account_management" db:"org_account_management"`
	MoneyMgmt  bool   `json:"money_management" db:"money_management"`
	ProjMgmt   bool   `json:"project_management" db:"project_management"`
}

type OrgBase struct {
	ID                    int       `json:"id"`
	Name                  string    `json:"name"`
	OwnerId               int       `json:"owner_id" db:"owner"`
	AvatarPath            *string   `json:"-" db:"avatar_path"`
	Email                 string    `json:"email"`
	Balance               float64   `json:"balance,omitempty"`
	OrgType               OrgType   `json:"org_type" db:"type"`
	OrgTypeId             int       `json:"-" db:"org_type_id"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	IsBanned              bool      `json:"is_banned" db:"is_banned"`
	RegistrationCompleted bool      `json:"registration_completed"`
}

type PhysFace struct {
	ID                                  int    `json:"-"`
	BIC                                 string `json:"bic" db:"bic"`
	CheckingAccount                     string `json:"checking_account" db:"checking_account"`
	CorrespondentAccount                string `json:"correspondent_account" db:"correspondent_account"`
	FIO                                 string `json:"fio" db:"fio"`
	INN                                 string `json:"inn,omitempty" db:"inn"`
	PassportSeries                      int    `json:"passport_series,omitempty" db:"pasport_series"`
	PassportNumber                      int    `json:"passport_number,omitempty" db:"pasport_number"`
	PassportGivenBy                     string `json:"passport_givenby,omitempty" db:"pasport_givenby"`
	RegistrationAddress                 string `json:"registration_address,omitempty" db:"registration_address"`
	PostAddress                         string `json:"post_address,omitempty" db:"post_address"`
	PassportPageWithPhotoPath           string `json:"-" db:"pasport_page_with_photo_path"`
	PassportPageWithPropiskaPath        string `json:"-" db:"pasport_page_with_propiska_path"`
	SvidOPostanovkeNaUchetPhysLitsaPath string `json:"-" db:"svid_o_postanovke_na_uchet_phys_litsa_path"`
}

type JurFace struct {
	ID                              int    `json:"-"`
	ActsOnBase                      string `json:"acts_on_base" db:"acts_on_base"`
	Position                        string `json:"position" db:"position"`
	BIC                             string `json:"bic" db:"bic"`
	CheckingAccount                 string `json:"checking_account" db:"checking_account"`
	CorrespondentAccount            string `json:"correspondent_account" db:"correspondent_account"`
	FullOrganisationName            string `json:"full_organisation_name" db:"full_organisation_name"`
	ShortOrganisationName           string `json:"short_organisation_name" db:"short_organisation_name"`
	INN                             string `json:"inn,omitempty" db:"inn"`
	OGRN                            string `json:"ogrn" db:"ogrn"`
	KPP                             string `json:"kpp,omitempty" db:"kpp"`
	JurAddress                      string `json:"jur_address,omitempty" db:"jur_address"`
	FactAddress                     string `json:"fact_address,omitempty" db:"fact_address"`
	PostAddress                     string `json:"post_address,omitempty" db:"post_address"`
	SvidORegistratsiiJurLitsaPath   string `json:"-" db:"svid_o_registratsii_jur_litsa_path"`
	SvidOPostanovkeNaNalogUchetPath string `json:"-" db:"svid_o_postanovke_na_nalog_uchet_path"`
	ProtocolONasznacheniiLitsaPath  string `json:"-" db:"protocol_o_nasznachenii_litsa_path"`
	USNPath                         string `json:"-" db:"usn_path"`
	UstavPath                       string `json:"-" db:"ustav_path"`
}

type IPFace struct {
	ID                              int    `json:"-"`
	BIC                             string `json:"bic" db:"bic"`
	RasSchot                        string `json:"ras_schot" db:"ras_schot"`
	KorSchot                        string `json:"kor_schot" db:"kor_schot"`
	FIO                             string `json:"fio" db:"fio"`
	IpSvidSerial                    int64  `json:"ip_svid_serial,omitempty" db:"ip_svid_serial"`
	IpSvidNumber                    int64  `json:"ip_svid_number,omitempty" db:"ip_svid_number"`
	IpSvidGivenBy                   string `json:"ip_svid_givenby,omitempty" db:"ip_svid_givenby"`
	INN                             string `json:"inn,omitempty" db:"inn"`
	OGRN                            string `json:"ogrn" db:"ogrn"`
	JurAddress                      string `json:"jur_address,omitempty" db:"jur_address"`
	FactAddress                     string `json:"fact_address,omitempty" db:"fact_address"`
	PostAddress                     string `json:"post_address,omitempty" db:"post_address"`
	SvidOPostanovkeNaNalogUchetPath string `json:"-" db:"svid_o_postanovke_na_nalog_uchet_path"`
	IpPassportPhotoPagePath         string `json:"-" db:"ip_pasport_photo_page_path"`
	IpPassportPropiskaPath          string `json:"-" db:"ip_pasport_propiska_path"`
	USNPath                         string `json:"-" db:"usn_path"`
	OGRNIPPath                      string `json:"-" db:"ogrnip_path"`
}

type OrgDocType string

const (
	DocPhysPassportPhoto      OrgDocType = "phys_passport_photo"
	DocPhysPassportPropiska   OrgDocType = "phys_passport_propiska"
	DocPhysUchet              OrgDocType = "phys_svid_uchet"
	DocJurRegSvid             OrgDocType = "jur_reg_svid"
	DocJurUchet               OrgDocType = "jur_svid_uchet"
	DocJurAppointmentProtocol OrgDocType = "jur_appointment_protocol"
	DocJurUSN                 OrgDocType = "jur_usn"
	DocJurUstav               OrgDocType = "jur_ustav"
	DocIPUchet                OrgDocType = "ip_svid_uchet"
	DocIPPassportPhoto        OrgDocType = "ip_passport_photo"
	DocIPPassportPropiska     OrgDocType = "ip_passport_propiska"
	DocIPUSN                  OrgDocType = "ip_usn"
	DocIPOGRNIP               OrgDocType = "ip_ogrnip"
)

func IsValidDocType(docType OrgDocType) bool {
	validTypes := map[OrgDocType]bool{
		DocPhysPassportPhoto:      true,
		DocPhysPassportPropiska:   true,
		DocPhysUchet:              true,
		DocJurRegSvid:             true,
		DocJurUchet:               true,
		DocJurAppointmentProtocol: true,
		DocJurUSN:                 true,
		DocJurUstav:               true,
		DocIPUchet:                true,
		DocIPPassportPhoto:        true,
		DocIPPassportPropiska:     true,
		DocIPUSN:                  true,
		DocIPOGRNIP:               true,
	}
	return validTypes[docType]
}

func (d OrgDocType) IsValidForOrgType(orgType OrgType) bool {
	switch orgType {
	case OrgTypePhys:
		return d == DocPhysPassportPhoto || d == DocPhysPassportPropiska || d == DocPhysUchet
	case OrgTypeJur:
		return d == DocJurRegSvid || d == DocJurUchet || d == DocJurAppointmentProtocol || d == DocJurUSN || d == DocJurUstav
	case OrgTypeIP:
		return d == DocIPUchet || d == DocIPPassportPhoto || d == DocIPPassportPropiska || d == DocIPUSN || d == DocIPOGRNIP
	}
	return false
}

type Org struct {
	OrgBase
	PhysFace *PhysFace `json:"phys_face,omitempty"`
	JurFace  *JurFace  `json:"jur_face,omitempty"`
	IPFace   *IPFace   `json:"ip_face,omitempty"`
}
