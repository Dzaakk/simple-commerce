package usecase

import (
	"Dzaakk/simple-commerce/internal/seller/model"
	"strconv"
)

func ConvertSellersToResData(sellers []*model.TSeller) []*model.ResData {
	resDataList := make([]*model.ResData, len(sellers))
	for i, seller := range sellers {
		resData := &model.ResData{
			Id:       strconv.FormatInt(seller.Id, 10),
			Username: seller.Username,
			Email:    seller.Email,
			Balance:  strconv.FormatFloat(seller.Balance, 'f', -1, 64),
		}
		resDataList[i] = resData
	}
	return resDataList
}
