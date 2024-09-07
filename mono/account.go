package mono

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/adelowo/go-mono/mono/util"
)

type AccountDetails struct {
	Account struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Currency      string `json:"currency"`
		Type          string `json:"type"`
		AccountNumber string `json:"account_number"`
		Balance       int    `json:"balance"`
		Bvn           string `json:"bvn"`
		Institution   struct {
			Name     string `json:"name"`
			BankCode string `json:"bank_code"`
			Type     string `json:"type"`
		} `json:"institution"`
	} `json:"account"`
	Meta struct {
		DataStatus string `json:"data_status"`
		AuthMethod string `json:"auth_method"`
	} `json:"meta"`
}

type AccountDetailsResponse struct {
	BaseMonoResponse
	Data AccountDetails `json:"data"`
}

type AccountService service

func (a *AccountService) Details(ctx context.Context,
	accountID string) (AccountDetails, error) {

	var resp AccountDetails

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return resp, err
	}

	req, err := a.client.newRequest(http.MethodGet,
		fmt.Sprintf("/accounts/%s", accountID),
		body)
	if err != nil {
		return resp, err
	}

	var account AccountDetailsResponse
	_, err = a.client.Do(ctx, req, &account)
	if err != nil {
		return resp, err
	}

	return account.Data, nil
}

func (a *AccountService) Unlink(ctx context.Context,
	accountID string) error {

	if util.IsStringEmpty(accountID) {
		return errors.New("please provide a valid accountID")
	}

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return err
	}

	req, err := a.client.newRequest(http.MethodPost,
		fmt.Sprintf("/accounts/%s/unlink", accountID), body)
	if err != nil {
		return err
	}

	_, err = a.client.Do(ctx, req, nil)
	return err
}

type FetchBalanceOptions struct {
	AccountID string
	Realtime  bool
}

func (a *AccountService) Balance(ctx context.Context,
	opts FetchBalanceOptions) (int64, error) {

	if util.IsStringEmpty(opts.AccountID) {
		return 0, errors.New("please provide a valid accountID")
	}

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return 0, err
	}

	req, err := a.client.newRequest(http.MethodGet,
		fmt.Sprintf("/accounts/%s/balance", opts.AccountID), body)
	if err != nil {
		return 0, err
	}

	if opts.Realtime {
		req.Header.Add("X-REALTIME", "true")
	}

	var response struct {
		BaseMonoResponse
		Data struct {
			Balance int64 `json:"balance,omitempty"`
		} `json:"data,omitempty"`
	}

	_, err = a.client.Do(ctx, req, response)
	return response.Data.Balance, err
}
