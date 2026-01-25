package mono

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ayinke-llc/hermes"
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
	Customer struct {
		ID string `json:"id,omitempty"`
	} `json:"customer,omitempty"`
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

	req, err := a.client.newRequest(http.MethodGet,
		fmt.Sprintf("/v2/accounts/%s", accountID),
		nil)
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

	if hermes.IsStringEmpty(accountID) {
		return errors.New("please provide a valid accountID")
	}

	req, err := a.client.newRequest(http.MethodPost,
		fmt.Sprintf("/v2/accounts/%s/unlink", accountID), nil)
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

	if hermes.IsStringEmpty(opts.AccountID) {
		return 0, errors.New("please provide a valid accountID")
	}

	req, err := a.client.newRequest(http.MethodGet,
		fmt.Sprintf("/v2/accounts/%s/balance", opts.AccountID), nil)
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

	_, err = a.client.Do(ctx, req, &response)
	return response.Data.Balance, err
}

// DataSyncOptions contains options for the DataSync method
type DataSyncOptions struct {
	AccountID                string
	AllowIncompleteStatement bool
}

// ENUM(SYNC_SUCCESSFUL,REAUTHORISATION_REQUIRED,INCOMPLETE_STATEMENT_ERROR)
type DataSyncCode string

// DataSyncResponse represents the response from the data sync endpoint
type DataSyncResponse struct {
	Status     string       `json:"status"`
	HasNewData bool         `json:"hasNewData"`
	Code       DataSyncCode `json:"code"`
}

func (a *AccountService) DataSync(ctx context.Context,
	opts DataSyncOptions) (DataSyncResponse, error) {

	var resp DataSyncResponse

	if hermes.IsStringEmpty(opts.AccountID) {
		return resp, errors.New("please provide a valid accountID")
	}

	path := fmt.Sprintf("/accounts/%s/sync", opts.AccountID)
	if opts.AllowIncompleteStatement {
		path += "?allow_incomplete_statement=true"
	}

	req, err := a.client.newRequest(http.MethodPost, path, nil)
	if err != nil {
		return resp, err
	}

	_, err = a.client.Do(ctx, req, &resp)
	return resp, err
}
