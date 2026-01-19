package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
	"github.com/Starostina-elena/investment_platform/services/organisation/middleware"
)

type mockService struct {
	createFunc func(ctx context.Context, o core.Org) (*core.Org, error)
	getFunc    func(ctx context.Context, id int) (*core.Org, error)
	updateFunc func(ctx context.Context, o core.Org) (*core.Org, error)
}

func (m *mockService) Create(ctx context.Context, o core.Org) (*core.Org, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, o)
	}
	return &core.Org{}, nil
}

func (m *mockService) Get(ctx context.Context, id int) (*core.Org, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockService) GetPublicInfoOrg(ctx context.Context, id int) (*core.Org, error) {
	return nil, nil
}

func (m *mockService) Update(ctx context.Context, o core.Org) (*core.Org, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, o)
	}
	return nil, nil
}

func (m *mockService) UploadAvatar(ctx context.Context, orgID int, userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return "", nil
}

func (m *mockService) DeleteAvatar(ctx context.Context, orgID int, userID int, avatarPath string) error {
	return nil
}

func (m *mockService) UpdateAvatarPath(ctx context.Context, orgID int, avatarPath string) error {
	return nil
}

func (m *mockService) UploadDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return "", nil
}

func (m *mockService) DeleteDoc(ctx context.Context, orgID int, userID int, docType core.OrgDocType) error {
	return nil
}

func (m *mockService) DownloadDoc(ctx context.Context, orgID int, userID int, isAdmin bool, docType core.OrgDocType) ([]byte, string, error) {
	return nil, "", nil
}

func (m *mockService) BanOrg(ctx context.Context, orgID int, banned bool) error {
	return nil
}

func (m *mockService) GetUsersOrgs(ctx context.Context, userID int) ([]core.Org, error) {
	return nil, nil
}

func (m *mockService) CheckUserOrgPermission(ctx context.Context, orgID int, userID int, permission core.OrgPermission) (bool, error) {
	return false, nil
}

func (m *mockService) AddEmployee(ctx context.Context, orgID int, userRequested int, userID int, orgAccMgmt, moneyMgmt, projMgmt bool) error {
	return nil
}

func (m *mockService) GetOrgEmployees(ctx context.Context, orgID int) ([]core.OrgEmployee, error) {
	return nil, nil
}

func (m *mockService) UpdateEmployeePermissions(ctx context.Context, orgID int, userRequested int, userID int, orgAccMgmt, moneyMgmt, projMgmt bool) error {
	return nil
}

func (m *mockService) DeleteEmployee(ctx context.Context, orgID int, userRequested int, userID int) error {
	return nil
}

func (m *mockService) TransferOwnership(ctx context.Context, orgID int, userRequested int, newOwnerID int) error {
	return nil
}

func TestCreateOrgHandler_Success(t *testing.T) {
	ms := &mockService{
		createFunc: func(ctx context.Context, o core.Org) (*core.Org, error) {
			o.ID = 42
			return &o, nil
		},
	}

	h := NewHandler(ms, *slog.Default())

	req := CreateOrgRequest{
		Name:    "Test Org",
		Email:   "test@example.com",
		OrgType: core.OrgTypePhys,
		PhysFace: &core.PhysFace{
			BIC:                  "044525225",
			CheckingAccount:      "40702810000000000000",
			CorrespondentAccount: "30101810400000000225",
			FIO:                  "Иванов И.И.",
			INN:                  "123456789012",
			PassportSeries:       1234,
			PassportNumber:       567890,
			PassportGivenBy:      "МВД РФ",
			RegistrationAddress:  "г. Москва",
			PostAddress:          "г. Москва",
		},
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/orgs", bytes.NewReader(body))
	userClaims := &middleware.UserClaims{UserID: 1}
	httpReq = httpReq.WithContext(
		middleware.SetClaimsInContext(httpReq.Context(), userClaims),
	)

	w := httptest.NewRecorder()
	CreateOrgHandler(h).ServeHTTP(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp core.Org
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ID != 42 {
		t.Errorf("expected org id 42, got %d", resp.ID)
	}
}

func TestCreateOrgHandler_MissingFields(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms, *slog.Default())

	tests := []struct {
		name    string
		request CreateOrgRequest
	}{
		{
			name: "missing name",
			request: CreateOrgRequest{
				Name:    "",
				Email:   "test@example.com",
				OrgType: core.OrgTypePhys,
			},
		},
		{
			name: "missing email",
			request: CreateOrgRequest{
				Name:    "Test",
				Email:   "",
				OrgType: core.OrgTypePhys,
			},
		},
		{
			name: "missing phys face",
			request: CreateOrgRequest{
				Name:    "Test",
				Email:   "test@example.com",
				OrgType: core.OrgTypePhys,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.request)
			httpReq := httptest.NewRequest("POST", "/orgs", bytes.NewReader(body))
			userClaims := &middleware.UserClaims{UserID: 1}
			httpReq = httpReq.WithContext(
				middleware.SetClaimsInContext(httpReq.Context(), userClaims),
			)

			w := httptest.NewRecorder()
			CreateOrgHandler(h).ServeHTTP(w, httpReq)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestCreateOrgHandler_Unauthorized(t *testing.T) {
	ms := &mockService{}
	h := NewHandler(ms, *slog.Default())

	req := CreateOrgRequest{
		Name:    "Test Org",
		Email:   "test@example.com",
		OrgType: core.OrgTypePhys,
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/orgs", bytes.NewReader(body))

	w := httptest.NewRecorder()
	CreateOrgHandler(h).ServeHTTP(w, httpReq)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestValidatePhysFace(t *testing.T) {
	tests := []struct {
		name    string
		face    *core.PhysFace
		wantErr bool
	}{
		{
			name: "valid phys face",
			face: &core.PhysFace{
				BIC:                  "044525225",
				CheckingAccount:      "40702810000000000000",
				CorrespondentAccount: "30101810400000000225",
				FIO:                  "Иванов И.И.",
				INN:                  "123456789012",
				PassportSeries:       1234,
				PassportNumber:       567890,
				PassportGivenBy:      "МВД РФ",
				RegistrationAddress:  "г. Москва",
				PostAddress:          "г. Москва",
			},
			wantErr: false,
		},
		{
			name: "missing BIC",
			face: &core.PhysFace{
				CheckingAccount:      "40702810000000000000",
				CorrespondentAccount: "30101810400000000225",
				FIO:                  "Иванов И.И.",
				INN:                  "123456789012",
				PassportSeries:       1234,
				PassportNumber:       567890,
				PassportGivenBy:      "МВД РФ",
				RegistrationAddress:  "г. Москва",
				PostAddress:          "г. Москва",
			},
			wantErr: true,
		},
		{
			name: "missing passport_series",
			face: &core.PhysFace{
				BIC:                  "044525225",
				CheckingAccount:      "40702810000000000000",
				CorrespondentAccount: "30101810400000000225",
				FIO:                  "Иванов И.И.",
				INN:                  "123456789012",
				PassportNumber:       567890,
				PassportGivenBy:      "МВД РФ",
				RegistrationAddress:  "г. Москва",
				PostAddress:          "г. Москва",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePhysFace(tt.face)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePhysFace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJurFace(t *testing.T) {
	tests := []struct {
		name    string
		face    *core.JurFace
		wantErr bool
	}{
		{
			name: "valid jur face",
			face: &core.JurFace{
				ActsOnBase:            "Устав",
				Position:              "Директор",
				BIC:                   "044525225",
				CheckingAccount:       "40702810000000000000",
				CorrespondentAccount:  "30101810400000000225",
				FullOrganisationName:  "ООО Компания",
				ShortOrganisationName: "Компания",
				INN:                   "7701234567",
				OGRN:                  "1227746234567",
				KPP:                   "770101001",
				JurAddress:            "г. Москва",
				FactAddress:           "г. Москва",
				PostAddress:           "г. Москва",
			},
			wantErr: false,
		},
		{
			name: "missing ActsOnBase",
			face: &core.JurFace{
				Position:              "Директор",
				BIC:                   "044525225",
				CheckingAccount:       "40702810000000000000",
				CorrespondentAccount:  "30101810400000000225",
				FullOrganisationName:  "ООО Компания",
				ShortOrganisationName: "Компания",
				INN:                   "7701234567",
				OGRN:                  "1227746234567",
				KPP:                   "770101001",
				JurAddress:            "г. Москва",
				FactAddress:           "г. Москва",
				PostAddress:           "г. Москва",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJurFace(tt.face)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateJurFace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateIPFace(t *testing.T) {
	tests := []struct {
		name    string
		face    *core.IPFace
		wantErr bool
	}{
		{
			name: "valid ip face",
			face: &core.IPFace{
				BIC:           "044525225",
				RasSchot:      "40702810000000000000",
				KorSchot:      "30101810400000000225",
				FIO:           "Петров П.П.",
				IpSvidSerial:  123456,
				IpSvidNumber:  789012,
				IpSvidGivenBy: "МВД РФ",
				INN:           "123456789012",
				OGRN:          "304001234567890",
				JurAddress:    "г. Москва",
				FactAddress:   "г. Москва",
				PostAddress:   "г. Москва",
			},
			wantErr: false,
		},
		{
			name: "missing BIC",
			face: &core.IPFace{
				RasSchot:      "40702810000000000000",
				KorSchot:      "30101810400000000225",
				FIO:           "Петров П.П.",
				IpSvidSerial:  123456,
				IpSvidNumber:  789012,
				IpSvidGivenBy: "МВД РФ",
				INN:           "123456789012",
				OGRN:          "304001234567890",
				JurAddress:    "г. Москва",
				FactAddress:   "г. Москва",
				PostAddress:   "г. Москва",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIPFace(tt.face)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateIPFace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
