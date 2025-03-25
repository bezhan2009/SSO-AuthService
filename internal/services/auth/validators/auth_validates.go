package validators

import (
	"SSO/pkg/errs"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
)

const (
	emptyValue = 0
)

func ValidateUserRegisterRequest(userRequest *ssov1.RegisterRequest) (err error) {
	if userRequest.GetUsername() == "" {
		return errs.ErrUsernameIsRequired
	}

	if userRequest.GetFirstName() == "" {
		return errs.ErrFirstNameIsRequired
	}

	if userRequest.GetLastName() == "" {
		return errs.ErrLastNameIsRequired
	}

	if userRequest.GetEmail() == "" {
		return errs.ErrEmailIsRequired
	}

	if userRequest.GetPassword() == "" {
		return errs.ErrPasswordIsRequired
	}

	return nil
}

func ValidateUserLoginRequest(userRequest *ssov1.LoginRequest) (err error) {
	if userRequest.GetUsername() == "" {
		return errs.ErrUsernameIsRequired
	}

	if userRequest.GetPassword() == "" {
		return errs.ErrPasswordIsRequired
	}

	if userRequest.GetAppLogin() == emptyValue {
		return errs.ErrAppLoginIsRequired
	}

	return nil
}
