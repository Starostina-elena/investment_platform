package repo

import (
	"context"
	"database/sql"
	"errors"
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
