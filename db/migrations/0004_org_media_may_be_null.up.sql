ALTER TABLE physical_face_project_account
    ALTER COLUMN pasport_page_with_photo_path DROP NOT NULL;
ALTER TABLE physical_face_project_account
    ALTER COLUMN pasport_page_with_propiska_path DROP NOT NULL;
ALTER TABLE physical_face_project_account
    ALTER COLUMN svid_o_postanovke_na_uchet_phys_litsa_path DROP NOT NULL;

ALTER TABLE juridical_face_project_accout
    ALTER COLUMN svid_o_registratsii_jur_litsa_path DROP NOT NULL;
ALTER TABLE juridical_face_project_accout
    ALTER COLUMN svid_o_postanovke_na_nalog_uchet_path DROP NOT NULL;
ALTER TABLE juridical_face_project_accout
    ALTER COLUMN protocol_o_nasznachenii_litsa_path DROP NOT NULL;
ALTER TABLE juridical_face_project_accout
    ALTER COLUMN USN_path DROP NOT NULL;
ALTER TABLE juridical_face_project_accout
    ALTER COLUMN ustav_path DROP NOT NULL;

ALTER TABLE ip_project_account
    ALTER COLUMN svid_o_postanovke_na_nalog_uchet_path DROP NOT NULL;
ALTER TABLE ip_project_account
    ALTER COLUMN ip_pasport_photo_page_path DROP NOT NULL;
ALTER TABLE ip_project_account
    ALTER COLUMN ip_pasport_propiska_path DROP NOT NULL;
ALTER TABLE ip_project_account
    ALTER COLUMN USN_path DROP NOT NULL;
ALTER TABLE ip_project_account
    ALTER COLUMN OGRNIP_path DROP NOT NULL;
