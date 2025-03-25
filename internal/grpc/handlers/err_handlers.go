package handlers

import (
	"SSO/pkg/errs"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Обработка ошибок, которые приводят к статусу 400 (Bad Request)
func handleBadRequestErrors(err error) bool {
	return errors.Is(err, errs.ErrUsernameUniquenessFailed) ||
		errors.Is(err, errs.ErrUsernameIsRequired) ||
		errors.Is(err, errs.ErrPasswordIsRequired) ||
		errors.Is(err, errs.ErrFirstNameIsRequired) ||
		errors.Is(err, errs.ErrLastNameIsRequired) ||
		errors.Is(err, errs.ErrAppLoginIsRequired) ||
		errors.Is(err, errs.ErrEmailIsRequired) ||
		errors.Is(err, errs.ErrIncorrectUsernameOrPassword) ||
		errors.Is(err, errs.ErrCategoryNameUniquenessFailed) ||
		errors.Is(err, errs.ErrOrderStatusNameUniquenessFailed) ||
		errors.Is(err, errs.ErrInvalidCredentials) ||
		errors.Is(err, errs.ErrOrderNotFound) ||
		errors.Is(err, errs.ErrInvalidStoreReviewID) ||
		errors.Is(err, errs.ErrPathParametrized) ||
		errors.Is(err, errs.ErrInvalidProductID) ||
		errors.Is(err, errs.ErrInvalidAddressID) ||
		errors.Is(err, errs.ErrInvalidOrderID) ||
		errors.Is(err, errs.ErrInvalidProductReviewID) ||
		errors.Is(err, errs.ErrInvalidAccountID) ||
		errors.Is(err, errs.ErrInvalidFeaturedProductID) ||
		errors.Is(err, errs.ErrInvalidAddressName) ||
		errors.Is(err, errs.ErrInvalidAccountNumber) ||
		errors.Is(err, errs.ErrAddressNameUniquenessFailed) ||
		errors.Is(err, errs.ErrAccountNumberUniquenessFailed) ||
		errors.Is(err, errs.ErrInvalidMinPrice) ||
		errors.Is(err, errs.ErrInvalidMaxPrice) ||
		errors.Is(err, errs.ErrInvalidPrice) ||
		errors.Is(err, errs.ErrInvalidID) ||
		errors.Is(err, errs.ErrInvalidPaymentID) ||
		errors.Is(err, errs.ErrInvalidRating) ||
		errors.Is(err, errs.ErrInvalidContent) ||
		errors.Is(err, errs.ErrInvalidComment) ||
		errors.Is(err, errs.ErrInvalidField) ||
		errors.Is(err, errs.ErrInvalidCategory) ||
		errors.Is(err, errs.ErrEmailIsEmpty) ||
		errors.Is(err, errs.ErrPasswordIsEmpty) ||
		errors.Is(err, errs.ErrUsernameIsEmpty) ||
		errors.Is(err, errs.ErrInvalidStore) ||
		errors.Is(err, errs.ErrInvalidStoreID) ||
		errors.Is(err, errs.ErrValidationFailed) ||
		errors.Is(err, errs.ErrStoreNameUniquenessFailed) ||
		errors.Is(err, errs.ErrNotEnoughProductInStock) ||
		errors.Is(err, errs.ErrDeleteFailed) ||
		errors.Is(err, errs.ErrInvalidTitle) ||
		errors.Is(err, errs.ErrOrderAlreadyPayed) ||
		errors.Is(err, errs.ErrInvalidDescription) ||
		errors.Is(err, errs.ErrInvalidAmount) ||
		errors.Is(err, errs.ErrInvalidQuantity) ||
		errors.Is(err, errs.ErrInsufficientFunds)
}

// Обработка ошибок, которые приводят к статусу 404 (Not Found)
func handleNotFoundErrors(err error) bool {
	return errors.Is(err, errs.ErrRecordNotFound) ||
		errors.Is(err, errs.ErrCategoryNotFound) ||
		errors.Is(err, errs.ErrOrderStatusNotFound) ||
		errors.Is(err, errs.ErrOrderNotFound) ||
		errors.Is(err, errs.ErrProductReviewNotFound) ||
		errors.Is(err, errs.ErrProductNotFound) ||
		errors.Is(err, errs.ErrAddressNotFound) ||
		errors.Is(err, errs.ErrFeaturedProductNotFound) ||
		errors.Is(err, errs.ErrPaymentNotFound) ||
		errors.Is(err, errs.ErrAccountNotFound) ||
		errors.Is(err, errs.ErrStoreNotFound) ||
		errors.Is(err, errs.ErrStoreReviewNotFound)
}

// Обработка ошибок, которые приводят к статусу 401 (Unauthorized)
func handleUnauthorizedErrors(err error) bool {
	return errors.Is(err, errs.ErrInvalidToken) ||
		errors.Is(err, errs.ErrUnauthorized) ||
		errors.Is(err, errs.ErrRefreshTokenExpired)
}

// HandleError Основная функция обработки ошибок
func HandleError(err error) error {
	if handleBadRequestErrors(err) {
		return status.Error(codes.InvalidArgument, err.Error())
	} else if errors.Is(err, errs.ErrPermissionDenied) {
		return status.Error(codes.PermissionDenied, err.Error())
	} else if handleNotFoundErrors(err) {
		return status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, errs.ErrFetchingProducts) {
		return status.Error(codes.FailedPrecondition, err.Error())
	} else if errors.Is(err, errs.WarningNoProductsFound) {
		return status.Error(codes.OK, err.Error())
	} else if handleUnauthorizedErrors(err) {
		return status.Error(codes.Unauthenticated, err.Error())
	} else {
		return status.Error(codes.Internal, err.Error())
	}
}
