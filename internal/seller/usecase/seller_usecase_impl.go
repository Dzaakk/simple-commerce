package usecase

import (
	"Dzaakk/simple-commerce/internal/seller/model"
	repo "Dzaakk/simple-commerce/internal/seller/repository"
	"Dzaakk/simple-commerce/package/template"
	"Dzaakk/simple-commerce/package/util"
	"context"
	"fmt"
	"strconv"
	"time"
)

type SellerUseCaseImpl struct {
	repo repo.SellerRepository
}

func NewSellerUseCase(repo repo.SellerRepository) SellerUseCase {
	return &SellerUseCaseImpl{repo}
}

func (s *SellerUseCaseImpl) Create(ctx context.Context, data model.ReqCreate) (int64, error) {
	hashedPassword, err := util.HashPassword(data.Password)
	if err != nil {
		return 0, err
	}

	seller := model.TSeller{
		Username: data.Username,
		Email:    data.Email,
		Password: string(hashedPassword),
		Balance:  0,
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}

	sellerID, err := s.repo.Create(ctx, seller)
	if err != nil {
		return 0, err
	}

	return sellerID, nil
}

func (s *SellerUseCaseImpl) Deactivate(ctx context.Context, sellerID int64) (int64, error) {

	rowsAffected, err := s.repo.Deactive(ctx, sellerID)
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil

}
func (s *SellerUseCaseImpl) FindAll(ctx context.Context) ([]*model.ResData, error) {
	sellers, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	listSeller := ConvertSellersToResData(sellers)
	return listSeller, nil
}

func (s *SellerUseCaseImpl) FindBySellerID(ctx context.Context, sellerID int64) (*model.ResData, error) {
	sellerData, err := s.repo.FindBySellerID(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	res := &model.ResData{
		Id:       fmt.Sprintf("%d", sellerData.ID),
		Username: sellerData.Username,
		Email:    sellerData.Email,
		Balance:  fmt.Sprintf("%.2f", sellerData.Balance),
	}

	return res, nil
}

func (s *SellerUseCaseImpl) FindByUsername(ctx context.Context, username string) (*model.ResData, error) {
	sellerData, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	res := &model.ResData{
		Id:       fmt.Sprintf("%d", sellerData.ID),
		Username: sellerData.Username,
		Email:    sellerData.Email,
		Balance:  fmt.Sprintf("%.2f", sellerData.Balance),
	}

	return res, nil
}

func (s *SellerUseCaseImpl) Update(ctx context.Context, data model.ReqUpdate) (int64, error) {
	sellerID, _ := strconv.ParseInt(data.Id, 0, 64)
	existingData, err := s.repo.FindBySellerID(ctx, sellerID)
	if err != nil {
		return 0, err
	}
	updatedData := generateDataUpdate(*existingData, data)

	rowsAffected, err := s.repo.Update(ctx, updatedData)
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func generateDataUpdate(existingData model.TSeller, newData model.ReqUpdate) model.TSeller {
	updatedData := existingData
	var email, username string

	if newData.Email != existingData.Email {
		email = newData.Email
	} else {
		email = existingData.Email
	}

	if newData.Username != existingData.Username {
		username = newData.Username
	} else {
		username = existingData.Username
	}

	updatedData.Email = email
	updatedData.Username = username

	return updatedData
}

func (s *SellerUseCaseImpl) ChangePassword(ctx context.Context, sellerID int64, newPassword string) (int64, error) {
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := s.repo.UpdatePassword(ctx, sellerID, string(hashedPassword))
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (s *SellerUseCaseImpl) FindByEmail(ctx context.Context, email string) (*model.TSeller, error) {
	sellerData, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return sellerData, nil
}
