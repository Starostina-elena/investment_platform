package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Starostina-elena/investment_platform/services/organisation/core"
	"github.com/Starostina-elena/investment_platform/services/organisation/repo"
)

type Service interface {
	Create(ctx context.Context, org core.Org) (*core.Org, error)
	Get(ctx context.Context, id int) (*core.Org, error)
}

type service struct {
	repo repo.RepoInterface
	log  slog.Logger
}

func NewService(r repo.RepoInterface, log slog.Logger) Service {
	return &service{repo: r, log: log}
}

func (s *service) Create(ctx context.Context, org core.Org) (*core.Org, error) {
	org.Balance = 0.0
	org.CreatedAt = time.Now()
	org.IsBanned = false
	id, err := s.repo.Create(ctx, &org)
	if err != nil {
		s.log.Error("failed to create organisation", "error", err)
		return nil, err
	}
	org.ID = id
	return &org, nil
}

func (s *service) Get(ctx context.Context, id int) (*core.Org, error) {
	org, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Error("failed to get organisation", "error", err)
		return nil, err
	}

	org.RegistrationCompleted = true

	switch org.OrgType {
	case core.OrgTypePhys:
		if org.PhysFace == nil {
			org.RegistrationCompleted = false
		} else if org.PhysFace.PassportPageWithPhotoPath == "" ||
			org.PhysFace.PassportPageWithPropiskaPath == "" ||
			org.PhysFace.SvidOPostanovkeNaUchetPhysLitsaPath == "" {
			org.RegistrationCompleted = false
		}
	case core.OrgTypeJur:
		if org.JurFace == nil {
			org.RegistrationCompleted = false
		} else if org.JurFace.SvidORegistratsiiJurLitsaPath == "" ||
			org.JurFace.SvidOPostanovkeNaNalogUchetPath == "" ||
			org.JurFace.ProtocolONasznacheniiLitsaPath == "" ||
			org.JurFace.USNPath == "" ||
			org.JurFace.UstavPath == "" {
			org.RegistrationCompleted = false
		}
	case core.OrgTypeIP:
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

	return org, nil
}
