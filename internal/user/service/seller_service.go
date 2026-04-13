package service

import (
	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"net/http"
	"strconv"
)

type SellerServiceImpl struct {
	Repo SellerRepository
}

func NewSellerService(repo SellerRepository) SellerService {
	return &SellerServiceImpl{Repo: repo}
}

func (s *SellerServiceImpl) Create(ctx context.Context, req *dto.RegisterSellerRequest) (string, error) {
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
		return response.NewAppError(http.StatusBadRequest, "invalid parameter seller id")
	}

	if sellerID <= 0 {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter seller id")
	}

	data := req.ToUpdateData(sellerID)

	rowsAffected, err := s.Repo.Update(ctx, data)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return response.NewAppError(http.StatusNotFound, "seller not found")
	}

	return nil
}

func (s *SellerServiceImpl) FindByEmail(ctx context.Context, email string) (*model.Seller, error) {
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

func (s *SellerServiceImpl) FindByShopName(ctx context.Context, name string) ([]dto.SellerRes, error) {
	data, err := s.Repo.FindByShopName(ctx, name)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []dto.SellerRes{}, nil
	}

	result := make([]dto.SellerRes, 0, len(data))
	for _, seller := range data {
		if seller == nil {
			continue
		}
		res := dto.ToSellerRes(seller)
		result = append(result, res)
	}

	return result, nil
}

func (s *SellerServiceImpl) UpdateStatus(ctx context.Context, sellerID string, status constant.UserStatus) error {
	return s.Repo.UpdateStatus(ctx, sellerID, status)
}

func (s *SellerServiceImpl) UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, sellerID string, status constant.UserStatus) error {
	return s.Repo.UpdateStatusWithTx(ctx, tx, sellerID, status)
}
