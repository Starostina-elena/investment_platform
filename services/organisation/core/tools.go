package core

func (org *Org) SetIsRegistrationCompleted() {
	org.RegistrationCompleted = true

	switch org.OrgType {
	case OrgTypePhys:
		if org.PhysFace == nil {
			org.RegistrationCompleted = false
		} else if org.PhysFace.PassportPageWithPhotoPath == "" ||
			org.PhysFace.PassportPageWithPropiskaPath == "" ||
			org.PhysFace.SvidOPostanovkeNaUchetPhysLitsaPath == "" {
			org.RegistrationCompleted = false
		}
	case OrgTypeJur:
		if org.JurFace == nil {
			org.RegistrationCompleted = false
		} else if org.JurFace.SvidORegistratsiiJurLitsaPath == "" ||
			org.JurFace.SvidOPostanovkeNaNalogUchetPath == "" ||
			org.JurFace.ProtocolONasznacheniiLitsaPath == "" ||
			org.JurFace.USNPath == "" ||
			org.JurFace.UstavPath == "" {
			org.RegistrationCompleted = false
		}
	case OrgTypeIP:
		if org.IPFace == nil {
			org.RegistrationCompleted = false
		} else if org.IPFace.SvidOPostanovkeNaNalogUchetPath == "" ||
			org.IPFace.IpPassportPhotoPagePath == "" ||
			org.IPFace.IpPassportPropiskaPath == "" ||
			org.IPFace.USNPath == "" ||
			org.IPFace.OGRNIPPath == "" {
			org.RegistrationCompleted = false
		}
	}
}
