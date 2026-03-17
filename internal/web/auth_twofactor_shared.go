package web

import (
	"errors"
	"fmt"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/appleauth"
)

var _ appleauth.SessionState = (*AuthSession)(nil)

func (s *AuthSession) TwoFactorMethod() string {
	if s == nil {
		return ""
	}
	return s.twoFactorMethod
}

func (s *AuthSession) TwoFactorPhoneID() int {
	if s == nil {
		return 0
	}
	return s.twoFactorPhoneID
}

func (s *AuthSession) TwoFactorPhoneMode() string {
	if s == nil {
		return ""
	}
	return s.twoFactorPhoneMode
}

func (s *AuthSession) TwoFactorDestination() string {
	if s == nil {
		return ""
	}
	return s.twoFactorDestination
}

func (s *AuthSession) TwoFactorCodeRequested() bool {
	return s != nil && s.twoFactorCodeRequested
}

func (s *AuthSession) SetPreparedTwoFactorState(method string, phoneID int, phoneMode, destination string, requested bool) {
	if s == nil {
		return
	}
	s.twoFactorMethod = method
	s.twoFactorPhoneID = phoneID
	s.twoFactorPhoneMode = phoneMode
	s.twoFactorDestination = destination
	s.twoFactorCodeRequested = requested
}

func (s *AuthSession) SetTwoFactorCodeRequested(requested bool) {
	if s == nil {
		return
	}
	s.twoFactorCodeRequested = requested
}

func webSharedAuthOptions(opts *authOptionsResponse) *appleauth.AuthOptions {
	if opts == nil {
		return &appleauth.AuthOptions{}
	}
	shared := &appleauth.AuthOptions{
		NoTrustedDevices: opts.NoTrustedDevices,
	}
	if len(opts.TrustedPhoneNumbers) == 0 {
		return shared
	}
	shared.TrustedPhoneNumbers = make([]appleauth.TrustedPhoneNumber, 0, len(opts.TrustedPhoneNumbers))
	for _, phone := range opts.TrustedPhoneNumbers {
		shared.TrustedPhoneNumbers = append(shared.TrustedPhoneNumbers, appleauth.TrustedPhoneNumber{
			ID:                 phone.ID,
			PushMode:           phone.PushMode,
			NumberWithDialCode: phone.NumberWithDialCode,
		})
	}
	return shared
}

func wrapWebTwoFactorFlowError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, appleauth.ErrNoTrustedPhoneNumbers) {
		return fmt.Errorf("2fa failed: %w", err)
	}
	var unsupported *appleauth.UnsupportedTwoFactorMethodError
	if errors.As(err, &unsupported) {
		return fmt.Errorf("2fa failed: %w", err)
	}
	return err
}
