package service

import (
	"Dzaakk/simple-commerce/internal/user/domain"
	"Dzaakk/simple-commerce/internal/user/dto"
	"context"
	"errors"
	"strconv"
)

type SellerServiceImpl struct {
	Repo SellerRepository
}

func NewSellerService(repo SellerRepository) SellerService {
	return &SellerServiceImpl{Repo: repo}
}

func (s *SellerServiceImpl) Create(ctx context.Context, req *dto.SellerCreateReq) (string, error) {
	data := req.ToCreateData()

	id, err := s.Repo.Create(ctx, data)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *SellerServiceImpl) Update(ctx context.Context, req *dto.SellerUpdateReq) error {
	sellerID, err := strconv.ParseInt(req.SellerID, 0, 64)
	if err != nil {
		return err
	}

	if sellerID <= 0 {
		return errors.New("invalid parameter seller id")
	}

	data := req.ToUpdateData(sellerID)

	rowsAffected, err := s.Repo.Update(ctx, data)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (s *SellerServiceImpl) FindByEmail(ctx context.Context, email string) (*domain.Seller, error) {
	data, err := s.Repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *SellerServiceImpl) FindByID(ctx context.Context, sellerID string) (*dto.SellerRes, error) {
	data, err := s.Repo.FindByID(ctx, sellerID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	seller := dto.ToSellerRes(data)

	return &seller, nil
}
