package mono

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/adelowo/go-mono/mono/util"
	"github.com/oklog/ulid/v2"
)

type AuthorisationService service

type ExchangeTokenRequest struct {
	Code string `json:"code,omitempty"`
}

type ExchangeTokenResponse struct {
	BaseMonoResponse
	Data struct {
		ID string `json:"id,omitempty"`
	} `json:"data,omitempty"`
}

func (a *AuthorisationService) ExchangeToken(ctx context.Context,
	opts ExchangeTokenRequest) (string, error) {

	if util.IsStringEmpty(opts.Code) {
		return "", errors.New("please provide the code to exchange")
	}

	var emptyResp string

	body, err := ToReader(opts)
	if err != nil {
		return emptyResp, err
	}

	req, err := a.client.newRequest(http.MethodPost,
		"/accounts/auth", body)
	if err != nil {
		return emptyResp, err
	}

	var tokenExchange ExchangeTokenResponse
	_, err = a.client.Do(ctx, req, &tokenExchange)
	if err != nil {
		return emptyResp, err
	}

	return tokenExchange.Data.ID, nil
}

// ENUM(reauth)
type ReauthorisationRequestScope string

type ReauthorisationRequest struct {
	Meta struct {
		Ref string `json:"ref,omitempty"`
	} `json:"meta,omitempty"`
	RedirectURL string                      `json:"redirect_url,omitempty"`
	AccountID   string                      `json:"account,omitempty"`
	Scope       ReauthorisationRequestScope `json:"scope,omitempty"`
}

type ReauthorisationResponse struct {
	BaseMonoResponse
	Data AccountReauthorisation `json:"data,omitempty"`
}

type AccountReauthorisation struct {
	MonoURL  string `json:"mono_url"`
	Customer string `json:"customer"`
	Meta     struct {
		Ref string `json:"ref"`
	} `json:"meta"`
	Scope       string    `json:"scope"`
	RedirectURL string    `json:"redirect_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func (a *AuthorisationService) Reauthorise(ctx context.Context,
	opts ReauthorisationRequest) (AccountReauthorisation, error) {

	var emptyResp AccountReauthorisation

	if util.IsStringEmpty(opts.AccountID) {
		return emptyResp, errors.New("please provide the account id to reauthorise")
	}

	if !opts.Scope.IsValid() {
		return emptyResp, fmt.Errorf("unsupported scope. only reauth is supported right now")
	}

	_, err := url.Parse(opts.RedirectURL)
	if err != nil {
		return emptyResp, fmt.Errorf("please provide a valid redirect url....%w", err)
	}

	if util.IsStringEmpty(opts.Meta.Ref) {
		id, err := ulid.New(ulid.Timestamp(time.Now()), ulid.DefaultEntropy())
		if err != nil {
			return emptyResp, fmt.Errorf("could not generate unique ULID reference..%w", err)
		}

		opts.Meta.Ref = id.String()
	}

	body, err := ToReader(opts)
	if err != nil {
		return emptyResp, err
	}

	req, err := a.client.newRequest(http.MethodPost,
		"/accounts/initiate", body)
	if err != nil {
		return emptyResp, err
	}

	var tokenExchange ReauthorisationResponse
	_, err = a.client.Do(ctx, req, &tokenExchange)
	if err != nil {
		return emptyResp, err
	}

	return tokenExchange.Data, nil
}
