package mono

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type DirectDebitService service

// ENUM(mandate)
type DirectPayMethod string

type InitiateMandateOptions struct {
	Amount      int64                          `json:"amount,omitempty"`
	Type        DirectDebitPaymentScheduleType `json:"type,omitempty"`
	Method      DirectPayMethod                `json:"method,omitempty"`
	MandateType MandateType                    `json:"mandate_type,omitempty"`
	DebitType   DebitType                      `json:"debit_type,omitempty"`
	Description string                         `json:"description,omitempty"`
	Reference   string                         `json:"reference,omitempty"`
	Customer    string                         `json:"customer,omitempty"`
	StartDate   time.Time                      `json:"start_date,omitempty"`
	EndDate     time.Time                      `json:"end_date,omitempty"`
	RedirectURL string                         `json:"redirect_url,omitempty"`
}

// ENUM(recurring-debit,one-time)
type DirectDebitPaymentScheduleType string

// ENUM(fixed,variable)
type DebitType string

// ENUM(emandata,gsm,signed)
type MandateType string

type CreatedManadateDetails struct {
	MonoURL     string      `json:"mono_url"`
	Type        string      `json:"type"`
	Method      string      `json:"method"`
	MandateType MandateType `json:"mandate_type"`
	Amount      int64       `json:"amount"`
	Description string      `json:"description"`
	Reference   string      `json:"reference"`
	Customer    string      `json:"customer"`
	RedirectURL string      `json:"redirect_url"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	StartDate   string      `json:"start_date"`
	EndDate     string      `json:"end_date"`
}

type CreatedManadateDetailsResponse struct {
	BaseMonoResponse
	Data CreatedManadateDetails `json:"data,omitempty"`
}

func (d *DirectDebitService) InitiateMandate(ctx context.Context,
	opts InitiateMandateOptions) (CreatedManadateDetails, error) {

	var resp CreatedManadateDetails

	body, err := ToReader(opts)
	if err != nil {
		return resp, err
	}

	req, err := d.client.newRequest(
		http.MethodPost, "/v2/payments/initiate", body)
	if err != nil {
		return resp, err
	}

	var res CreatedManadateDetailsResponse

	_, err = d.client.Do(ctx, req, &res)
	return res.Data, err
}

type FetchMandateDetails struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	MandateType   string `json:"mandate_type"`
	DebitType     string `json:"debit_type"`
	Approved      bool   `json:"approved"`
	Amount        int    `json:"amount"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	Institution   struct {
		BankCode string `json:"bank_code"`
		NipCode  string `json:"nip_code"`
		Name     string `json:"name"`
	} `json:"institution"`
	Customer  string    `json:"customer"`
	Narration string    `json:"narration"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Date      time.Time `json:"date"`
}

type FetchMandateDetailsResponse struct {
	BaseMonoResponse
	Data FetchMandateDetails `json:"data,omitempty"`
}

func (d *DirectDebitService) Details(ctx context.Context,
	mandateID string) (FetchMandateDetails, error) {

	var resp FetchMandateDetails

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return resp, err
	}

	req, err := d.client.newRequest(http.MethodGet,
		fmt.Sprintf("/v3/payments/mandates/%s", mandateID), body)
	if err != nil {
		return resp, err
	}

	var mandate FetchMandateDetailsResponse

	_, err = d.client.Do(ctx, req, &mandate)
	return mandate.Data, err
}

func (d *DirectDebitService) Reinstate(ctx context.Context, mandateID string) error {

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return err
	}

	req, err := d.client.newRequest(http.MethodPatch,
		fmt.Sprintf("/v3/payments/mandates/%s/reinstate", mandateID), body)
	if err != nil {
		return err
	}

	_, err = d.client.Do(ctx, req, nil)
	return err
}

func (d *DirectDebitService) Pause(ctx context.Context, mandateID string) error {

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return err
	}

	req, err := d.client.newRequest(http.MethodPatch,
		fmt.Sprintf("/v3/payments/mandates/%s/pause", mandateID), body)
	if err != nil {
		return err
	}

	_, err = d.client.Do(ctx, req, nil)
	return err
}
