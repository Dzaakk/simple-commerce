package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/model"
	emailmodel "Dzaakk/simple-commerce/internal/email/model"
	userdto "Dzaakk/simple-commerce/internal/user/dto"
	usermodel "Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthCustomerService struct {
	createFn       func(context.Context, *userdto.RegisterCustomerRequest) (string, error)
	findByEmailFn  func(context.Context, string) (*usermodel.Customer, error)
	findByIDFn     func(context.Context, string) (*userdto.CustomerRes, error)
	updateStatusFn func(context.Context, string, constant.UserStatus) error
}

func (m *mockAuthCustomerService) Create(ctx context.Context, req *userdto.RegisterCustomerRequest) (string, error) {
	if m.createFn == nil {
		return "", errors.New("unexpected customer Create call")
	}
	return m.createFn(ctx, req)
}

func (m *mockAuthCustomerService) FindByEmail(ctx context.Context, email string) (*usermodel.Customer, error) {
	if m.findByEmailFn == nil {
		return nil, errors.New("unexpected customer FindByEmail call")
	}
	return m.findByEmailFn(ctx, email)
}

func (m *mockAuthCustomerService) FindByID(ctx context.Context, customerID string) (*userdto.CustomerRes, error) {
	if m.findByIDFn == nil {
		return nil, errors.New("unexpected customer FindByID call")
	}
	return m.findByIDFn(ctx, customerID)
}

func (m *mockAuthCustomerService) UpdateStatus(ctx context.Context, customerID string, status constant.UserStatus) error {
	if m.updateStatusFn == nil {
		return errors.New("unexpected customer UpdateStatus call")
	}
	return m.updateStatusFn(ctx, customerID, status)
}

type mockAuthSellerService struct {
	createFn       func(context.Context, *userdto.RegisterSellerRequest) (string, error)
	findByEmailFn  func(context.Context, string) (*usermodel.Seller, error)
	findByIDFn     func(context.Context, string) (*userdto.SellerRes, error)
	updateStatusFn func(context.Context, string, constant.UserStatus) error
}

func (m *mockAuthSellerService) Create(ctx context.Context, req *userdto.RegisterSellerRequest) (string, error) {
	if m.createFn == nil {
		return "", errors.New("unexpected seller Create call")
	}
	return m.createFn(ctx, req)
}

func (m *mockAuthSellerService) FindByEmail(ctx context.Context, email string) (*usermodel.Seller, error) {
	if m.findByEmailFn == nil {
		return nil, errors.New("unexpected seller FindByEmail call")
	}
	return m.findByEmailFn(ctx, email)
}

func (m *mockAuthSellerService) FindByID(ctx context.Context, sellerID string) (*userdto.SellerRes, error) {
	if m.findByIDFn == nil {
		return nil, errors.New("unexpected seller FindByID call")
	}
	return m.findByIDFn(ctx, sellerID)
}

func (m *mockAuthSellerService) UpdateStatus(ctx context.Context, sellerID string, status constant.UserStatus) error {
	if m.updateStatusFn == nil {
		return errors.New("unexpected seller UpdateStatus call")
	}
	return m.updateStatusFn(ctx, sellerID, status)
}

type mockAuthEmailService struct {
	sendEmailVerificationFn func(context.Context, emailmodel.VerificationEmailReq) error
}

func (m *mockAuthEmailService) SendEmailVerification(ctx context.Context, req emailmodel.VerificationEmailReq) error {
	if m.sendEmailVerificationFn == nil {
		return errors.New("unexpected SendEmailVerification call")
	}
	return m.sendEmailVerificationFn(ctx, req)
}

type mockActivationEmailPublisher struct {
	publishVerificationEmailFn func(context.Context, emailmodel.VerificationEmailReq) error
}

func (m *mockActivationEmailPublisher) PublishVerificationEmail(ctx context.Context, req emailmodel.VerificationEmailReq) error {
	if m.publishVerificationEmailFn == nil {
		return errors.New("unexpected PublishVerificationEmail call")
	}
	return m.publishVerificationEmailFn(ctx, req)
}

type mockActivationCodeRepository struct {
	createFn     func(context.Context, *model.ActivationCode) (int64, error)
	findByCodeFn func(context.Context, string) (*model.ActivationCode, error)
	markAsUsedFn func(context.Context, int64) error
}

func (m *mockActivationCodeRepository) Create(ctx context.Context, data *model.ActivationCode) (int64, error) {
	if m.createFn == nil {
		return 0, errors.New("unexpected activation Create call")
	}
	return m.createFn(ctx, data)
}

func (m *mockActivationCodeRepository) FindByCode(ctx context.Context, code string) (*model.ActivationCode, error) {
	if m.findByCodeFn == nil {
		return nil, errors.New("unexpected activation FindByCode call")
	}
	return m.findByCodeFn(ctx, code)
}

func (m *mockActivationCodeRepository) MarkAsUsed(ctx context.Context, id int64) error {
	if m.markAsUsedFn == nil {
		return errors.New("unexpected activation MarkAsUsed call")
	}
	return m.markAsUsedFn(ctx, id)
}

type mockRefreshTokenRepository struct {
	createFn          func(context.Context, *model.RefreshToken) (int64, error)
	findByTokenHashFn func(context.Context, string) (*model.RefreshToken, error)
	revokeFn          func(context.Context, string) error
	revokeAllByUserFn func(context.Context, string, constant.UserType) error
}

func (m *mockRefreshTokenRepository) Create(ctx context.Context, data *model.RefreshToken) (int64, error) {
	if m.createFn == nil {
		return 0, errors.New("unexpected refresh Create call")
	}
	return m.createFn(ctx, data)
}

func (m *mockRefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	if m.findByTokenHashFn == nil {
		return nil, errors.New("unexpected refresh FindByTokenHash call")
	}
	return m.findByTokenHashFn(ctx, tokenHash)
}

func (m *mockRefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	if m.revokeFn == nil {
		return errors.New("unexpected refresh Revoke call")
	}
	return m.revokeFn(ctx, tokenHash)
}

func (m *mockRefreshTokenRepository) RevokeAllByUser(ctx context.Context, userID string, userType constant.UserType) error {
	if m.revokeAllByUserFn == nil {
		return errors.New("unexpected refresh RevokeAllByUser call")
	}
	return m.revokeAllByUserFn(ctx, userID, userType)
}

type mockTransactor struct {
	called bool
	fn     func(context.Context, func(context.Context) error) error
}

func (m *mockTransactor) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	m.called = true
	if m.fn != nil {
		return m.fn(ctx, fn)
	}
	return fn(ctx)
}

func TestAuthServiceRegisterCustomer(t *testing.T) {
	ctx := context.Background()
	t.Setenv("BASE_URL", "https://commerce.test")

	req := &dto.RegisterCustomerRequest{
		Email:    "customer@example.com",
		Password: "plain-password",
		FullName: "Customer Name",
		Phone:    "08123456789",
	}

	var activationCode string
	customerSvc := &mockAuthCustomerService{
		findByEmailFn: func(_ context.Context, email string) (*usermodel.Customer, error) {
			if email != req.Email {
				t.Fatalf("email = %q, want %q", email, req.Email)
			}
			return nil, nil
		},
		createFn: func(_ context.Context, got *userdto.RegisterCustomerRequest) (string, error) {
			if got.Email != req.Email || got.FullName != req.FullName || got.Phone != req.Phone {
				t.Fatalf("create request = %#v", got)
			}
			if got.Password == req.Password {
				t.Fatal("password must be hashed before creating customer")
			}
			if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(req.Password)); err != nil {
				t.Fatalf("password hash does not match plain password: %v", err)
			}
			return "customer-1", nil
		},
	}
	activationRepo := &mockActivationCodeRepository{
		createFn: func(_ context.Context, got *model.ActivationCode) (int64, error) {
			if got.Email != req.Email {
				t.Fatalf("activation email = %q, want %q", got.Email, req.Email)
			}
			if got.UserType != string(constant.Customer) {
				t.Fatalf("user type = %q, want %q", got.UserType, constant.Customer)
			}
			if len(got.Code) != activationCodeBytes*2 {
				t.Fatalf("code length = %d, want %d", len(got.Code), activationCodeBytes*2)
			}
			if time.Until(got.ExpiresAt) < 14*time.Minute || time.Until(got.ExpiresAt) > 16*time.Minute {
				t.Fatalf("expires_at = %v, want about 15 minutes from now", got.ExpiresAt)
			}
			activationCode = got.Code
			return 1, nil
		},
	}
	publisher := &mockActivationEmailPublisher{
		publishVerificationEmailFn: func(_ context.Context, got emailmodel.VerificationEmailReq) error {
			if got.Email != req.Email || got.Username != req.FullName {
				t.Fatalf("email request = %#v", got)
			}
			wantLink := "https://commerce.test/api/v1/auth/verify-email?code=" + activationCode
			if got.ActivationLink != wantLink {
				t.Fatalf("activation link = %q, want %q", got.ActivationLink, wantLink)
			}
			return nil
		},
	}

	err := NewAuthService(nil, customerSvc, nil, nil, publisher, activationRepo, nil).RegisterCustomer(ctx, req)
	if err != nil {
		t.Fatalf("RegisterCustomer returned error: %v", err)
	}
}

func TestAuthServiceRegisterCustomerReturnsEmailAlreadyExists(t *testing.T) {
	createCalled := false
	customerSvc := &mockAuthCustomerService{
		findByEmailFn: func(context.Context, string) (*usermodel.Customer, error) {
			return &usermodel.Customer{ID: "customer-1"}, nil
		},
		createFn: func(context.Context, *userdto.RegisterCustomerRequest) (string, error) {
			createCalled = true
			return "", nil
		},
	}

	err := NewAuthService(nil, customerSvc, nil, nil, nil, nil, nil).RegisterCustomer(context.Background(), &dto.RegisterCustomerRequest{})
	if !errors.Is(err, response.ErrEmailAlreadyExist) {
		t.Fatalf("error = %v, want %v", err, response.ErrEmailAlreadyExist)
	}
	if createCalled {
		t.Fatal("customer Create must not be called when email already exists")
	}
}

func TestAuthServiceRegisterSeller(t *testing.T) {
	ctx := context.Background()
	t.Setenv("BASE_URL", "https://commerce.test")

	req := &dto.RegisterSellerRequest{
		Email:    "seller@example.com",
		Password: "plain-password",
		FullName: "Seller Name",
		Phone:    "08123456789",
		ShopName: "Seller Shop",
	}

	var activationCode string
	sellerSvc := &mockAuthSellerService{
		findByEmailFn: func(_ context.Context, email string) (*usermodel.Seller, error) {
			if email != req.Email {
				t.Fatalf("email = %q, want %q", email, req.Email)
			}
			return nil, nil
		},
		createFn: func(_ context.Context, got *userdto.RegisterSellerRequest) (string, error) {
			if got.Email != req.Email || got.FullName != req.FullName || got.Phone != req.Phone || got.ShopName != req.ShopName {
				t.Fatalf("create request = %#v", got)
			}
			if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(req.Password)); err != nil {
				t.Fatalf("password hash does not match plain password: %v", err)
			}
			return "seller-1", nil
		},
	}
	activationRepo := &mockActivationCodeRepository{
		createFn: func(_ context.Context, got *model.ActivationCode) (int64, error) {
			if got.Email != req.Email || got.UserType != string(constant.Seller) {
				t.Fatalf("activation data = %#v", got)
			}
			activationCode = got.Code
			return 1, nil
		},
	}
	publisher := &mockActivationEmailPublisher{
		publishVerificationEmailFn: func(_ context.Context, got emailmodel.VerificationEmailReq) error {
			if !strings.HasSuffix(got.ActivationLink, "/api/v1/auth/verify-email?code="+activationCode) {
				t.Fatalf("activation link = %q", got.ActivationLink)
			}
			return nil
		},
	}

	err := NewAuthService(nil, nil, sellerSvc, nil, publisher, activationRepo, nil).RegisterSeller(ctx, req)
	if err != nil {
		t.Fatalf("RegisterSeller returned error: %v", err)
	}
}

func TestAuthServiceLoginCustomer(t *testing.T) {
	ctx := context.Background()
	t.Setenv("JWT_SECRET", "test-secret")

	passwordHash, err := hashPassword("plain-password")
	if err != nil {
		t.Fatalf("hashPassword returned error: %v", err)
	}
	var storedRefresh *model.RefreshToken

	customerSvc := &mockAuthCustomerService{
		findByEmailFn: func(_ context.Context, email string) (*usermodel.Customer, error) {
			if email != "customer@example.com" {
				t.Fatalf("email = %q, want customer@example.com", email)
			}
			return &usermodel.Customer{
				ID:           "customer-1",
				Email:        "customer@example.com",
				PasswordHash: passwordHash,
				Status:       string(constant.StatusActive),
			}, nil
		},
	}
	refreshRepo := &mockRefreshTokenRepository{
		createFn: func(_ context.Context, got *model.RefreshToken) (int64, error) {
			if got.UserID != "customer-1" || got.UserType != constant.Customer {
				t.Fatalf("refresh token = %#v", got)
			}
			if len(got.TokenHash) != 64 {
				t.Fatalf("token hash length = %d, want 64", len(got.TokenHash))
			}
			if time.Until(got.ExpiresAt) < (refreshTokenDuration - time.Minute) {
				t.Fatalf("expires_at = %v, want about %v from now", got.ExpiresAt, refreshTokenDuration)
			}
			storedRefresh = got
			return 1, nil
		},
	}

	got, err := NewAuthService(nil, customerSvc, nil, nil, nil, nil, refreshRepo).Login(ctx, &dto.LoginRequest{
		Email:    "customer@example.com",
		Password: "plain-password",
		UserType: string(constant.Customer),
	})
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if got.AccessToken == "" || got.RefreshToken == "" {
		t.Fatalf("tokens must not be empty: %#v", got)
	}
	if got.ExpiresIn != int(accessTokenDuration.Seconds()) {
		t.Fatalf("expires_in = %d, want %d", got.ExpiresIn, int(accessTokenDuration.Seconds()))
	}
	if storedRefresh == nil {
		t.Fatal("refresh token was not stored")
	}
	if hashRefreshToken(got.RefreshToken) != storedRefresh.TokenHash {
		t.Fatalf("stored hash = %q, want hash of returned refresh token", storedRefresh.TokenHash)
	}
	assertAccessTokenClaims(t, got.AccessToken, "customer-1", string(constant.Customer), "customer@example.com")
}

func TestAuthServiceLoginRejectsPendingCustomer(t *testing.T) {
	passwordHash, err := hashPassword("plain-password")
	if err != nil {
		t.Fatalf("hashPassword returned error: %v", err)
	}

	customerSvc := &mockAuthCustomerService{
		findByEmailFn: func(context.Context, string) (*usermodel.Customer, error) {
			return &usermodel.Customer{
				ID:           "customer-1",
				Email:        "customer@example.com",
				PasswordHash: passwordHash,
				Status:       string(constant.StatusPending),
			}, nil
		},
	}

	_, err = NewAuthService(nil, customerSvc, nil, nil, nil, nil, nil).Login(context.Background(), &dto.LoginRequest{
		Email:    "customer@example.com",
		Password: "plain-password",
		UserType: string(constant.Customer),
	})
	if !errors.Is(err, response.ErrEmailNotVerified) {
		t.Fatalf("error = %v, want %v", err, response.ErrEmailNotVerified)
	}
}

func TestAuthServiceLoginRejectsWrongPassword(t *testing.T) {
	passwordHash, err := hashPassword("plain-password")
	if err != nil {
		t.Fatalf("hashPassword returned error: %v", err)
	}

	customerSvc := &mockAuthCustomerService{
		findByEmailFn: func(context.Context, string) (*usermodel.Customer, error) {
			return &usermodel.Customer{
				ID:           "customer-1",
				Email:        "customer@example.com",
				PasswordHash: passwordHash,
				Status:       string(constant.StatusActive),
			}, nil
		},
	}

	_, err = NewAuthService(nil, customerSvc, nil, nil, nil, nil, nil).Login(context.Background(), &dto.LoginRequest{
		Email:    "customer@example.com",
		Password: "wrong-password",
		UserType: string(constant.Customer),
	})
	if !errors.Is(err, response.ErrInvalidCredentials) {
		t.Fatalf("error = %v, want %v", err, response.ErrInvalidCredentials)
	}
}

func TestAuthServiceRefreshTokenCustomer(t *testing.T) {
	ctx := context.Background()
	t.Setenv("JWT_SECRET", "test-secret")

	rawRefresh := "refresh-token"
	refreshRepo := &mockRefreshTokenRepository{
		findByTokenHashFn: func(_ context.Context, tokenHash string) (*model.RefreshToken, error) {
			if tokenHash != hashRefreshToken(rawRefresh) {
				t.Fatalf("token hash = %q, want hash of raw refresh token", tokenHash)
			}
			return &model.RefreshToken{
				UserID:   "customer-1",
				UserType: constant.Customer,
			}, nil
		},
	}
	customerSvc := &mockAuthCustomerService{
		findByIDFn: func(_ context.Context, customerID string) (*userdto.CustomerRes, error) {
			if customerID != "customer-1" {
				t.Fatalf("customer id = %q, want customer-1", customerID)
			}
			return &userdto.CustomerRes{ID: "customer-1", Email: "customer@example.com"}, nil
		},
	}

	got, err := NewAuthService(nil, customerSvc, nil, nil, nil, nil, refreshRepo).RefreshToken(ctx, rawRefresh)
	if err != nil {
		t.Fatalf("RefreshToken returned error: %v", err)
	}
	if got.AccessToken == "" {
		t.Fatal("access token must not be empty")
	}
	if got.ExpiresIn != int(accessTokenDuration.Seconds()) {
		t.Fatalf("expires_in = %d, want %d", got.ExpiresIn, int(accessTokenDuration.Seconds()))
	}
	assertAccessTokenClaims(t, got.AccessToken, "customer-1", string(constant.Customer), "customer@example.com")
}

func TestAuthServiceRefreshTokenRejectsUnknownToken(t *testing.T) {
	refreshRepo := &mockRefreshTokenRepository{
		findByTokenHashFn: func(context.Context, string) (*model.RefreshToken, error) {
			return nil, nil
		},
	}

	_, err := NewAuthService(nil, nil, nil, nil, nil, nil, refreshRepo).RefreshToken(context.Background(), "missing")
	if !errors.Is(err, response.ErrInvalidRefreshToken) {
		t.Fatalf("error = %v, want %v", err, response.ErrInvalidRefreshToken)
	}
}

func TestAuthServiceLogoutRevokesHashedRefreshToken(t *testing.T) {
	rawRefresh := "refresh-token"
	refreshRepo := &mockRefreshTokenRepository{
		revokeFn: func(_ context.Context, tokenHash string) error {
			if tokenHash != hashRefreshToken(rawRefresh) {
				t.Fatalf("token hash = %q, want hash of raw refresh token", tokenHash)
			}
			return nil
		},
	}

	err := NewAuthService(nil, nil, nil, nil, nil, nil, refreshRepo).Logout(context.Background(), rawRefresh)
	if err != nil {
		t.Fatalf("Logout returned error: %v", err)
	}
}

func TestAuthServiceVerifyEmailRejectsInvalidCode(t *testing.T) {
	activationRepo := &mockActivationCodeRepository{
		findByCodeFn: func(_ context.Context, code string) (*model.ActivationCode, error) {
			if code != "bad-code" {
				t.Fatalf("code = %q, want bad-code", code)
			}
			return nil, nil
		},
	}

	err := NewAuthService(nil, nil, nil, nil, nil, activationRepo, nil).VerifyEmail(context.Background(), "bad-code")
	if !errors.Is(err, response.ErrInvalidActivationCode) {
		t.Fatalf("error = %v, want %v", err, response.ErrInvalidActivationCode)
	}
}

func TestAuthServiceVerifyEmailActivatesCustomerWithinTransaction(t *testing.T) {
	txManager := &mockTransactor{}
	activationRepo := &mockActivationCodeRepository{
		findByCodeFn: func(_ context.Context, code string) (*model.ActivationCode, error) {
			if code != "valid-code" {
				t.Fatalf("code = %q, want valid-code", code)
			}
			return &model.ActivationCode{ID: 7, Email: "customer@example.com", UserType: string(constant.Customer)}, nil
		},
		markAsUsedFn: func(_ context.Context, id int64) error {
			if id != 7 {
				t.Fatalf("activation code id = %d, want 7", id)
			}
			return nil
		},
	}
	customerSvc := &mockAuthCustomerService{
		findByEmailFn: func(_ context.Context, email string) (*usermodel.Customer, error) {
			if email != "customer@example.com" {
				t.Fatalf("email = %q, want customer@example.com", email)
			}
			return &usermodel.Customer{ID: "customer-1", Email: email}, nil
		},
		updateStatusFn: func(_ context.Context, customerID string, status constant.UserStatus) error {
			if customerID != "customer-1" {
				t.Fatalf("customer id = %q, want customer-1", customerID)
			}
			if status != constant.StatusActive {
				t.Fatalf("status = %q, want %q", status, constant.StatusActive)
			}
			return nil
		},
	}

	err := NewAuthService(txManager, customerSvc, nil, nil, nil, activationRepo, nil).VerifyEmail(context.Background(), "valid-code")
	if err != nil {
		t.Fatalf("VerifyEmail returned error: %v", err)
	}
	if !txManager.called {
		t.Fatal("VerifyEmail must run activation updates inside a transaction")
	}
}

func assertAccessTokenClaims(t *testing.T, tokenString, userID, userType, email string) {
	t.Helper()

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse access token: %v", err)
	}
	if !token.Valid {
		t.Fatal("access token is not valid")
	}
	if claims.UserID != userID || claims.UserType != userType || claims.Email != email {
		t.Fatalf("claims = %#v, want user_id=%q user_type=%q email=%q", claims, userID, userType, email)
	}
}
