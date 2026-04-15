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

type CustomerServiceImpl struct {
	Repo   CustomerRepository
	Logger *logging.Logger
}

func NewCustomerService(repo CustomerRepository) CustomerService {
	return &CustomerServiceImpl{
		Repo:   repo,
		Logger: logging.NewLogger("user", "customer_service"),
	}
}

func (c *CustomerServiceImpl) Create(ctx context.Context, req *dto.RegisterCustomerRequest) (string, error) {

	data := req.ToCreateData()

	id, err := c.Repo.Create(ctx, data)
	if err != nil {
		c.Logger.Error(ctx, "customer_create_failed", map[string]interface{}{
			"operation": "create_customer",
		})
		return "", err
	}

	c.Logger.Info(ctx, "customer_created", map[string]interface{}{
		"customer_id": id,
	})
	return id, nil
}

func (c *CustomerServiceImpl) Update(ctx context.Context, req *dto.UpdateReq) error {

	customerID, err := strconv.ParseInt(req.CustomerID, 0, 64)
	if err != nil {
		c.Logger.Warn(ctx, "customer_update_invalid_id", map[string]interface{}{
			"operation": "update_customer",
		})
		return response.NewAppError(http.StatusBadRequest, "invalid parameter customer id")
	}

	if customerID <= 0 {
		c.Logger.Warn(ctx, "customer_update_invalid_id", map[string]interface{}{
			"operation": "update_customer",
		})
		return response.NewAppError(http.StatusBadRequest, "invalid parameter customer id")
	}

	data := req.ToUpdateData(customerID)

	rowsAffected, err := c.Repo.Update(ctx, data)
	if err != nil {
		c.Logger.Error(ctx, "customer_update_failed", map[string]interface{}{
			"customer_id": customerID,
			"operation":   "update_customer",
		})
		return err
	}
	if rowsAffected == 0 {
		c.Logger.Warn(ctx, "customer_update_not_found", map[string]interface{}{
			"customer_id": customerID,
			"operation":   "update_customer",
		})
		return response.NewAppError(http.StatusNotFound, "customer not found")
	}

	c.Logger.Info(ctx, "customer_updated", map[string]interface{}{
		"customer_id": customerID,
	})
	return nil
}

func (c *CustomerServiceImpl) FindByEmail(ctx context.Context, email string) (*model.Customer, error) {

	data, err := c.Repo.FindByEmail(ctx, email)
	if err != nil {
		c.Logger.Error(ctx, "customer_find_by_email_failed", map[string]interface{}{
			"operation": "find_customer_by_email",
		})
		return nil, err
	}
	if data == nil {
		c.Logger.Info(ctx, "customer_not_found", map[string]interface{}{
			"lookup": "email",
		})
		return nil, nil
	}

	c.Logger.Info(ctx, "customer_found", map[string]interface{}{
		"lookup":      "email",
		"customer_id": data.ID,
	})
	return data, nil
}

func (c *CustomerServiceImpl) FindByID(ctx context.Context, customerID string) (*dto.CustomerRes, error) {

	data, err := c.Repo.FindByID(ctx, customerID)
	if err != nil {
		c.Logger.Error(ctx, "customer_find_by_id_failed", map[string]interface{}{
			"customer_id": customerID,
			"operation":   "find_customer_by_id",
		})
		return nil, err
	}
	if data == nil {
		c.Logger.Info(ctx, "customer_not_found", map[string]interface{}{
			"lookup":      "id",
			"customer_id": customerID,
		})
		return nil, nil
	}

	customer := dto.ToCustomerRes(data)

	c.Logger.Info(ctx, "customer_found", map[string]interface{}{
		"lookup":      "id",
		"customer_id": customerID,
	})
	return &customer, nil
}

func (c *CustomerServiceImpl) UpdateStatus(ctx context.Context, customerID string, status constant.UserStatus) error {
	if err := c.Repo.UpdateStatus(ctx, customerID, status); err != nil {
		c.Logger.Error(ctx, "customer_status_update_failed", map[string]interface{}{
			"customer_id": customerID,
			"status":      status,
		})
		return err
	}

	c.Logger.Info(ctx, "customer_status_updated", map[string]interface{}{
		"customer_id": customerID,
		"status":      status,
	})
	return nil
}

func (c *CustomerServiceImpl) UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, customerID string, status constant.UserStatus) error {
	if err := c.Repo.UpdateStatusWithTx(ctx, tx, customerID, status); err != nil {
		c.Logger.Error(ctx, "customer_status_update_with_tx_failed", map[string]interface{}{
			"customer_id": customerID,
			"status":      status,
		})
		return err
	}

	c.Logger.Info(ctx, "customer_status_updated_with_tx", map[string]interface{}{
		"customer_id": customerID,
		"status":      status,
	})
	return nil
}
