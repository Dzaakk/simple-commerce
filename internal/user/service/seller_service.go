package service

import (
	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/logging"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"net/http"
	"strconv"
)

type SellerServiceImpl struct {
	Repo   SellerRepository
	Logger *logging.Logger
}

func NewSellerService(repo SellerRepository) SellerService {
	return &SellerServiceImpl{
		Repo:   repo,
		Logger: logging.NewLogger("user", "seller_service"),
	}
}

func (s *SellerServiceImpl) Create(ctx context.Context, req *dto.RegisterSellerRequest) (string, error) {
	data := req.ToCreateData()

	id, err := s.Repo.Create(ctx, data)
	if err != nil {
		s.Logger.Error(ctx, "seller_create_failed", map[string]interface{}{
			"operation": "create_seller",
		})
		return "", err
	}

	s.Logger.Info(ctx, "seller_created", map[string]interface{}{
		"seller_id": id,
	})
	return id, nil
}

func (s *SellerServiceImpl) Update(ctx context.Context, req *dto.SellerUpdateReq) error {
	sellerID, err := strconv.ParseInt(req.SellerID, 0, 64)
	if err != nil {
		s.Logger.Warn(ctx, "seller_update_invalid_id", map[string]interface{}{
			"operation": "update_seller",
		})
		return response.NewAppError(http.StatusBadRequest, "invalid parameter seller id")
	}

	if sellerID <= 0 {
		s.Logger.Warn(ctx, "seller_update_invalid_id", map[string]interface{}{
			"operation": "update_seller",
		})
		return response.NewAppError(http.StatusBadRequest, "invalid parameter seller id")
	}

	data := req.ToUpdateData(sellerID)

	rowsAffected, err := s.Repo.Update(ctx, data)
	if err != nil {
		s.Logger.Error(ctx, "seller_update_failed", map[string]interface{}{
			"seller_id":  sellerID,
			"operation":  "update_seller",
		})
		return err
	}
	if rowsAffected == 0 {
		s.Logger.Warn(ctx, "seller_update_not_found", map[string]interface{}{
			"seller_id":  sellerID,
			"operation":  "update_seller",
		})
		return response.NewAppError(http.StatusNotFound, "seller not found")
	}

	s.Logger.Info(ctx, "seller_updated", map[string]interface{}{
		"seller_id": sellerID,
	})
	return nil
}

func (s *SellerServiceImpl) FindByEmail(ctx context.Context, email string) (*model.Seller, error) {
	data, err := s.Repo.FindByEmail(ctx, email)
	if err != nil {
		s.Logger.Error(ctx, "seller_find_by_email_failed", map[string]interface{}{
			"operation": "find_seller_by_email",
		})
		return nil, err
	}
	if data == nil {
		s.Logger.Info(ctx, "seller_not_found", map[string]interface{}{
			"lookup": "email",
		})
		return nil, nil
	}

	s.Logger.Info(ctx, "seller_found", map[string]interface{}{
		"lookup":    "email",
		"seller_id": data.ID,
	})
	return data, nil
}

func (s *SellerServiceImpl) FindByID(ctx context.Context, sellerID string) (*dto.SellerRes, error) {
	data, err := s.Repo.FindByID(ctx, sellerID)
	if err != nil {
		s.Logger.Error(ctx, "seller_find_by_id_failed", map[string]interface{}{
			"seller_id":  sellerID,
			"operation":  "find_seller_by_id",
		})
		return nil, err
	}
	if data == nil {
		s.Logger.Info(ctx, "seller_not_found", map[string]interface{}{
			"lookup":    "id",
			"seller_id": sellerID,
		})
		return nil, nil
	}

	seller := dto.ToSellerRes(data)

	s.Logger.Info(ctx, "seller_found", map[string]interface{}{
		"lookup":    "id",
		"seller_id": sellerID,
	})
	return &seller, nil
}

func (s *SellerServiceImpl) FindByShopName(ctx context.Context, name string) ([]dto.SellerRes, error) {
	data, err := s.Repo.FindByShopName(ctx, name)
	if err != nil {
		s.Logger.Error(ctx, "seller_find_by_shop_name_failed", map[string]interface{}{
			"operation": "find_seller_by_shop_name",
		})
		return nil, err
	}
	if len(data) == 0 {
		s.Logger.Info(ctx, "seller_not_found", map[string]interface{}{
			"lookup": "shop_name",
		})
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

	s.Logger.Info(ctx, "seller_found", map[string]interface{}{
		"lookup": "shop_name",
		"count":  len(result),
	})
	return result, nil
}

func (s *SellerServiceImpl) UpdateStatus(ctx context.Context, sellerID string, status constant.UserStatus) error {
	if err := s.Repo.UpdateStatus(ctx, sellerID, status); err != nil {
		s.Logger.Error(ctx, "seller_status_update_failed", map[string]interface{}{
			"seller_id": sellerID,
			"status":    status,
		})
		return err
	}

	s.Logger.Info(ctx, "seller_status_updated", map[string]interface{}{
		"seller_id": sellerID,
		"status":    status,
	})
	return nil
}

func (s *SellerServiceImpl) UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, sellerID string, status constant.UserStatus) error {
	if err := s.Repo.UpdateStatusWithTx(ctx, tx, sellerID, status); err != nil {
		s.Logger.Error(ctx, "seller_status_update_with_tx_failed", map[string]interface{}{
			"seller_id": sellerID,
			"status":    status,
		})
		return err
	}

	s.Logger.Info(ctx, "seller_status_updated_with_tx", map[string]interface{}{
		"seller_id": sellerID,
		"status":    status,
	})
	return nil
}
