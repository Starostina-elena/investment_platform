package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
)

type Repo struct {
	db  *sqlx.DB
	log slog.Logger
}

type RepoInterface interface {
	Create(ctx context.Context, o *core.Org) (int, error)
	Get(ctx context.Context, id int) (*core.Org, error)
	Update(ctx context.Context, o *core.Org) (*core.Org, error)
	UpdateAvatarPath(ctx context.Context, orgID int, avatarPath *string) error
	UpdateDocPath(ctx context.Context, orgID int, docType core.OrgDocType, path string) error
	GetDocPath(ctx context.Context, orgID int, docType core.OrgDocType) (string, error)
	GetUsersOrgs(ctx context.Context, userID int) ([]core.Org, error)
	BanOrg(ctx context.Context, orgID int, banned bool) error
}

func NewRepo(db *sqlx.DB, log slog.Logger) RepoInterface {
	return &Repo{db: db, log: log}
}

func (r *Repo) Create(ctx context.Context, o *core.Org) (int, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var OrgId, detailedOrgId int

	// write to database additinal info (physical/juridical/ip face)
	switch {
	case o.OrgType == core.OrgTypePhys && o.PhysFace != nil:
		err = tx.QueryRowContext(ctx, `
		INSERT INTO physical_face_project_account 
		(BIC, checking_account, correspondent_account, FIO, INN, pasport_series,
		pasport_number, pasport_givenby, registration_address, post_address,
		pasport_page_with_photo_path, pasport_page_with_propiska_path,
		svid_o_postanovke_na_uchet_phys_litsa_path) VALUES
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id
		`, o.PhysFace.BIC, o.PhysFace.CheckingAccount, o.PhysFace.CorrespondentAccount,
			o.PhysFace.FIO, o.PhysFace.INN, o.PhysFace.PassportSeries, o.PhysFace.PassportNumber,
			o.PhysFace.PassportGivenBy, o.PhysFace.RegistrationAddress, o.PhysFace.PostAddress,
			"", "", "").Scan(&detailedOrgId)
	case o.OrgType == core.OrgTypeJur && o.JurFace != nil:
		err = tx.QueryRowContext(ctx, `
		INSERT INTO juridical_face_project_accout
		(acts_on_base, position, BIC, checking_account, correspondent_account,
		full_organisation_name, short_organisation_name, INN, OGRN, KPP,
		jur_address, fact_address, post_address,
		svid_o_registratsii_jur_litsa_path, svid_o_postanovke_na_nalog_uchet_path,
		protocol_o_nasznachenii_litsa_path, USN_path, ustav_path) VALUES
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18) RETURNING id
		`, o.JurFace.ActsOnBase, o.JurFace.Position, o.JurFace.BIC,
			o.JurFace.CheckingAccount, o.JurFace.CorrespondentAccount,
			o.JurFace.FullOrganisationName, o.JurFace.ShortOrganisationName, o.JurFace.INN,
			o.JurFace.OGRN, o.JurFace.KPP, o.JurFace.JurAddress, o.JurFace.FactAddress,
			o.JurFace.PostAddress, "", "", "", "", "").Scan(&detailedOrgId)
	case o.OrgType == core.OrgTypeIP && o.IPFace != nil:
		err = tx.QueryRowContext(ctx, `
		INSERT INTO ip_project_account
		(BIC, ras_schot, kor_schot, FIO, ip_svid_serial, ip_svid_number, ip_svid_givenby,
		INN, OGRN, jur_address, fact_address, post_address,
		svid_o_postanovke_na_nalog_uchet_path, ip_pasport_photo_page_path,
		ip_pasport_propiska_path, USN_path, OGRNIP_path) VALUES
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17) RETURNING id
		`, o.IPFace.BIC, o.IPFace.RasSchot, o.IPFace.KorSchot, o.IPFace.FIO,
			o.IPFace.IpSvidSerial, o.IPFace.IpSvidNumber, o.IPFace.IpSvidGivenBy,
			o.IPFace.INN, o.IPFace.OGRN, o.IPFace.JurAddress, o.IPFace.FactAddress,
			o.IPFace.PostAddress, "", "", "", "", "").Scan(&detailedOrgId)
	default:
		return 0, errors.New("invalid organisation type or missing face details")
	}

	if err != nil {
		r.log.Error("failed to create detailed organisation info", "error", err)
		return 0, err
	}

	// save main organisation info
	err = tx.QueryRowContext(ctx, `
		INSERT INTO organizations (name, owner, email, type, org_type_id)
		VALUES ($1,$2,$3,$4,$5) RETURNING id
		`, o.Name, o.OwnerId, o.Email, o.OrgType, detailedOrgId).Scan(&OrgId)
	if err != nil {
		r.log.Error("failed to create organisation", "error", err)
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return OrgId, nil
}

func (r *Repo) Get(ctx context.Context, id int) (*core.Org, error) {
	o := &core.Org{}
	if err := r.db.GetContext(ctx, &o.OrgBase, `
		SELECT id, name, owner, avatar_path, email, balance, type, org_type_id, created_at, is_banned
		FROM organizations WHERE id = $1
	`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
		r.log.Error("failed to get organisation", "id", id, "error", err)
		return nil, err
	}

	switch o.OrgType {
	case core.OrgTypePhys:
		phys := &core.PhysFace{}
		if err := r.db.GetContext(ctx, phys, `
			SELECT id, BIC, checking_account, correspondent_account, FIO, INN, pasport_series,
			pasport_number, pasport_givenby, registration_address, post_address,
			pasport_page_with_photo_path, pasport_page_with_propiska_path,
			svid_o_postanovke_na_uchet_phys_litsa_path
			FROM physical_face_project_account WHERE id = $1
		`, o.OrgTypeId); err != nil && !errors.Is(err, sql.ErrNoRows) {
			r.log.Error("failed to get phys face", "id", o.OrgTypeId, "error", err)
			return nil, err
		} else if err == nil {
			o.PhysFace = phys
		}

	case core.OrgTypeJur:
		jur := &core.JurFace{}
		if err := r.db.GetContext(ctx, jur, `
			SELECT id, acts_on_base, position, BIC, checking_account, correspondent_account,
			full_organisation_name, short_organisation_name, INN, OGRN, KPP,
			jur_address, fact_address, post_address, svid_o_registratsii_jur_litsa_path,
			svid_o_postanovke_na_nalog_uchet_path, protocol_o_nasznachenii_litsa_path,
			USN_path, ustav_path
			FROM juridical_face_project_accout WHERE id = $1
		`, o.OrgTypeId); err != nil && !errors.Is(err, sql.ErrNoRows) {
			r.log.Error("failed to get jur face", "id", o.OrgTypeId, "error", err)
			return nil, err
		} else if err == nil {
			o.JurFace = jur
		}

	case core.OrgTypeIP:
		ip := &core.IPFace{}
		if err := r.db.GetContext(ctx, ip, `
			SELECT id, BIC, ras_schot, kor_schot, FIO, ip_svid_serial, ip_svid_number, ip_svid_givenby,
			INN, OGRN, jur_address, fact_address, post_address,
			svid_o_postanovke_na_nalog_uchet_path, ip_pasport_photo_page_path,
			ip_pasport_propiska_path, USN_path, OGRNIP_path
			FROM ip_project_account WHERE id = $1
		`, o.OrgTypeId); err != nil && !errors.Is(err, sql.ErrNoRows) {
			r.log.Error("failed to get ip face", "id", o.OrgTypeId, "error", err)
			return nil, err
		} else if err == nil {
			o.IPFace = ip
		}
	}

	return o, nil
}

func (r *Repo) Update(ctx context.Context, o *core.Org) (*core.Org, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, `
		UPDATE organizations 
		SET name = $1, email = $2
		WHERE id = $3
	`, o.Name, o.Email, o.ID)
	if err != nil {
		r.log.Error("failed to update organisation", "error", err)
		return nil, err
	}

	switch {
	case o.OrgType == core.OrgTypePhys && o.PhysFace != nil:
		_, err = tx.ExecContext(ctx, `
			UPDATE physical_face_project_account
			SET BIC = $1, checking_account = $2, correspondent_account = $3, FIO = $4, 
				INN = $5, pasport_series = $6, pasport_number = $7, pasport_givenby = $8,
				registration_address = $9, post_address = $10
			WHERE id = $11
		`, o.PhysFace.BIC, o.PhysFace.CheckingAccount, o.PhysFace.CorrespondentAccount,
			o.PhysFace.FIO, o.PhysFace.INN, o.PhysFace.PassportSeries, o.PhysFace.PassportNumber,
			o.PhysFace.PassportGivenBy, o.PhysFace.RegistrationAddress, o.PhysFace.PostAddress,
			o.OrgTypeId)
	case o.OrgType == core.OrgTypeJur && o.JurFace != nil:
		_, err = tx.ExecContext(ctx, `
			UPDATE juridical_face_project_accout
			SET acts_on_base = $1, position = $2, BIC = $3, checking_account = $4,
				correspondent_account = $5, full_organisation_name = $6, short_organisation_name = $7,
				INN = $8, OGRN = $9, KPP = $10, jur_address = $11, fact_address = $12,
				post_address = $13
			WHERE id = $14
		`, o.JurFace.ActsOnBase, o.JurFace.Position, o.JurFace.BIC,
			o.JurFace.CheckingAccount, o.JurFace.CorrespondentAccount,
			o.JurFace.FullOrganisationName, o.JurFace.ShortOrganisationName, o.JurFace.INN,
			o.JurFace.OGRN, o.JurFace.KPP, o.JurFace.JurAddress, o.JurFace.FactAddress,
			o.JurFace.PostAddress, o.OrgTypeId)
	case o.OrgType == core.OrgTypeIP && o.IPFace != nil:
		_, err = tx.ExecContext(ctx, `
			UPDATE ip_project_account
			SET BIC = $1, ras_schot = $2, kor_schot = $3, FIO = $4, ip_svid_serial = $5,
				ip_svid_number = $6, ip_svid_givenby = $7, INN = $8, OGRN = $9,
				jur_address = $10, fact_address = $11, post_address = $12
			WHERE id = $13
		`, o.IPFace.BIC, o.IPFace.RasSchot, o.IPFace.KorSchot, o.IPFace.FIO,
			o.IPFace.IpSvidSerial, o.IPFace.IpSvidNumber, o.IPFace.IpSvidGivenBy,
			o.IPFace.INN, o.IPFace.OGRN, o.IPFace.JurAddress, o.IPFace.FactAddress,
			o.IPFace.PostAddress, o.OrgTypeId)
	default:
		return nil, errors.New("invalid organisation type or missing face details")
	}

	if err != nil {
		r.log.Error("failed to update detailed organisation info", "error", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return o, nil
}

func (r *Repo) UpdateAvatarPath(ctx context.Context, orgID int, avatarPath *string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE organizations SET avatar_path = $1 WHERE id = $2`, avatarPath, orgID)
	return err
}

func (r *Repo) UpdateDocPath(ctx context.Context, orgID int, docType core.OrgDocType, path string) error {
	query, err := buildDocUpdateQuery(docType)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, query, path, orgID)
	return err
}

func (r *Repo) GetDocPath(ctx context.Context, orgID int, docType core.OrgDocType) (string, error) {
	query, err := buildDocSelectQuery(docType)
	if err != nil {
		return "", err
	}
	var path sql.NullString
	if err := r.db.QueryRowContext(ctx, query, orgID).Scan(&path); err != nil {
		return "", err
	}
	if !path.Valid || path.String == "" {
		return "", nil
	}
	return path.String, nil
}

func buildDocUpdateQuery(docType core.OrgDocType) (string, error) {
	switch docType {
	case core.DocPhysPassportPhoto:
		return `UPDATE physical_face_project_account SET pasport_page_with_photo_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocPhysPassportPropiska:
		return `UPDATE physical_face_project_account SET pasport_page_with_propiska_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocPhysUchet:
		return `UPDATE physical_face_project_account SET svid_o_postanovke_na_uchet_phys_litsa_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocJurRegSvid:
		return `UPDATE juridical_face_project_accout SET svid_o_registratsii_jur_litsa_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocJurUchet:
		return `UPDATE juridical_face_project_accout SET svid_o_postanovke_na_nalog_uchet_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocJurAppointmentProtocol:
		return `UPDATE juridical_face_project_accout SET protocol_o_nasznachenii_litsa_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocJurUSN:
		return `UPDATE juridical_face_project_accout SET usn_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocJurUstav:
		return `UPDATE juridical_face_project_accout SET ustav_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocIPUchet:
		return `UPDATE ip_project_account SET svid_o_postanovke_na_nalog_uchet_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocIPPassportPhoto:
		return `UPDATE ip_project_account SET ip_pasport_photo_page_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocIPPassportPropiska:
		return `UPDATE ip_project_account SET ip_pasport_propiska_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocIPUSN:
		return `UPDATE ip_project_account SET usn_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	case core.DocIPOGRNIP:
		return `UPDATE ip_project_account SET ogrnip_path = $1 WHERE id = (SELECT org_type_id FROM organizations WHERE id = $2)`, nil
	default:
		return "", fmt.Errorf("unknown doc type: %s", docType)
	}
}

func buildDocSelectQuery(docType core.OrgDocType) (string, error) {
	switch docType {
	case core.DocPhysPassportPhoto:
		return `SELECT pasport_page_with_photo_path FROM physical_face_project_account WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocPhysPassportPropiska:
		return `SELECT pasport_page_with_propiska_path FROM physical_face_project_account WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocPhysUchet:
		return `SELECT svid_o_postanovke_na_uchet_phys_litsa_path FROM physical_face_project_account WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocJurRegSvid:
		return `SELECT svid_o_registratsii_jur_litsa_path FROM juridical_face_project_accout WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocJurUchet:
		return `SELECT svid_o_postanovke_na_nalog_uchet_path FROM juridical_face_project_accout WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocJurAppointmentProtocol:
		return `SELECT protocol_o_nasznachenii_litsa_path FROM juridical_face_project_accout WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocJurUSN:
		return `SELECT usn_path FROM juridical_face_project_accout WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocJurUstav:
		return `SELECT ustav_path FROM juridical_face_project_accout WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocIPUchet:
		return `SELECT svid_o_postanovke_na_nalog_uchet_path FROM ip_project_account WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocIPPassportPhoto:
		return `SELECT ip_pasport_photo_page_path FROM ip_project_account WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocIPPassportPropiska:
		return `SELECT ip_pasport_propiska_path FROM ip_project_account WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocIPUSN:
		return `SELECT usn_path FROM ip_project_account WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	case core.DocIPOGRNIP:
		return `SELECT ogrnip_path FROM ip_project_account WHERE id = (SELECT org_type_id FROM organizations WHERE id = $1)`, nil
	default:
		return "", fmt.Errorf("unknown doc type: %s", docType)
	}
}

func (r *Repo) GetUsersOrgs(ctx context.Context, userID int) ([]core.Org, error) {
	var orgIDs []core.Org
	err := r.db.SelectContext(ctx, &orgIDs, `
		SELECT id
		FROM organizations WHERE owner = $1
	`, userID)
	if err != nil {
		r.log.Error("failed to get user's organisations", "user_id", userID, "error", err)
		return nil, err
	}

	var orgs []core.Org
	for _, orgID := range orgIDs {
		org, err := r.Get(ctx, orgID.ID)
		if err != nil {
			r.log.Error("failed to get organisation details", "org_id", orgID.ID, "error", err)
			return nil, err
		}
		orgs = append(orgs, *org)
	}

	return orgs, nil
}

func (r *Repo) BanOrg(ctx context.Context, orgID int, banned bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE organizations SET is_banned = $1 WHERE id = $2`, banned, orgID)
	return err
}
