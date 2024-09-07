package mono

import (
	"context"
	"net/http"
	"time"
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
	Status    string         `json:"status"`
	Message   string         `json:"message"`
	Timestamp time.Time      `json:"timestamp"`
	Data      AccountDetails `json:"data"`
}

type AccountService service

func (a *AccountService) Details(ctx context.Context,
	accountID string) (AccountDetails, error) {

	var resp AccountDetails

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return resp, err
	}

	req, err := a.client.newRequest(http.MethodGet, "/providers", body)
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
