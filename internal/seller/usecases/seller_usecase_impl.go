package usecases

import (
	model "Dzaakk/simple-commerce/internal/seller/models"
	repo "Dzaakk/simple-commerce/internal/seller/repositories"
	template "Dzaakk/simple-commerce/package/templates"
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
	hashedPassword, err := template.HashPassword(data.Password)
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

	sellerId, err := s.repo.Create(ctx, seller)
	if err != nil {
		return 0, err
	}

	return sellerId, nil
}

func (s *SellerUseCaseImpl) Deactivate(ctx context.Context, sellerId int64) (int64, error) {

	rowsAffected, err := s.repo.Deactive(ctx, sellerId)
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil

}

func (s *SellerUseCaseImpl) FindById(ctx context.Context, sellerId int64) (*model.ResData, error) {
	sellerData, err := s.repo.FindById(ctx, sellerId)
	if err != nil {
		return nil, err
	}

	res := &model.ResData{
		Id:       fmt.Sprintf("%d", sellerData.Id),
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
		Id:       fmt.Sprintf("%d", sellerData.Id),
		Username: sellerData.Username,
		Email:    sellerData.Email,
		Balance:  fmt.Sprintf("%.2f", sellerData.Balance),
	}

	return res, nil
}

func (s *SellerUseCaseImpl) Update(ctx context.Context, data model.ReqUpdate) (int64, error) {
	sellerId, _ := strconv.ParseInt(data.Id, 0, 64)
	existingData, err := s.repo.FindById(ctx, sellerId)
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

func (s *SellerUseCaseImpl) ChangePassword(ctx context.Context, sellerId int64, newPassword string) (int64, error) {
	hashedPassword, err := template.HashPassword(newPassword)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := s.repo.UpdatePassword(ctx, sellerId, string(hashedPassword))
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
