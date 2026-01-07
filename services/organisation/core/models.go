package core

import "time"

type OrgType string

const (
	OrgTypeJur  OrgType = "jur"
	OrgTypePhys OrgType = "phys"
	OrgTypeIP   OrgType = "ip"
)

type OrgBase struct {
	ID                    int       `json:"id"`
	Name                  string    `json:"name"`
	OwnerId               int       `json:"owner_id" db:"owner"`
	AvatarPath            *string   `json:"-" db:"avatar_path"`
	Email                 string    `json:"email"`
	Balance               float64   `json:"balance"`
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
	INN                                 string `json:"inn" db:"inn"`
	PassportSeries                      int    `json:"passport_series" db:"pasport_series"`
	PassportNumber                      int    `json:"passport_number" db:"pasport_number"`
	PassportGivenBy                     string `json:"passport_givenby" db:"pasport_givenby"`
	RegistrationAddress                 string `json:"registration_address" db:"registration_address"`
	PostAddress                         string `json:"post_address" db:"post_address"`
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
	INN                             string `json:"inn" db:"inn"`
	OGRN                            string `json:"ogrn" db:"ogrn"`
	KPP                             string `json:"kpp" db:"kpp"`
	JurAddress                      string `json:"jur_address" db:"jur_address"`
	FactAddress                     string `json:"fact_address" db:"fact_address"`
	PostAddress                     string `json:"post_address" db:"post_address"`
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
	IpSvidSerial                    int64  `json:"ip_svid_serial" db:"ip_svid_serial"`
	IpSvidNumber                    int64  `json:"ip_svid_number" db:"ip_svid_number"`
	IpSvidGivenBy                   string `json:"ip_svid_givenby" db:"ip_svid_givenby"`
	INN                             string `json:"inn" db:"inn"`
	OGRN                            string `json:"ogrn" db:"ogrn"`
	JurAddress                      string `json:"jur_address" db:"jur_address"`
	FactAddress                     string `json:"fact_address" db:"fact_address"`
	PostAddress                     string `json:"post_address" db:"post_address"`
	SvidOPostanovkeNaNalogUchetPath string `json:"-" db:"svid_o_postanovke_na_nalog_uchet_path"`
	IpPassportPhotoPagePath         string `json:"-" db:"ip_pasport_photo_page_path"`
	IpPassportPropiskaPath          string `json:"-" db:"ip_pasport_propiska_path"`
	USNPath                         string `json:"-" db:"usn_path"`
	OGRNIPPath                      string `json:"-" db:"ogrnip_path"`
}

type Org struct {
	OrgBase
	PhysFace *PhysFace `json:"phys_face,omitempty"`
	JurFace  *JurFace  `json:"jur_face,omitempty"`
	IPFace   *IPFace   `json:"ip_face,omitempty"`
}
