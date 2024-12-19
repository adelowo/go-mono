package mono

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ayinke-llc/hermes"
)

type DirectDebitService service

// ENUM(mandate)
type DirectPayMethod string

type CreatedManadateOptions struct {
	InitiateMandateOptions
	AccountNumber string `json:"account_number,omitempty"`
	BankCode      string `json:"bank_code,omitempty"`
	Signature     string `json:"signature,omitempty"`
	AccountID     string `json:"account,omitempty"`
	Customer      string `json:"customer,omitempty"`
}

type InitiateMandateOptions struct {
	Amount      int64                          `json:"amount,omitempty"`
	Type        DirectDebitPaymentScheduleType `json:"type,omitempty"`
	Method      DirectPayMethod                `json:"method,omitempty"`
	MandateType MandateType                    `json:"mandate_type,omitempty"`
	DebitType   DebitType                      `json:"debit_type,omitempty"`
	Description string                         `json:"description,omitempty"`
	Reference   string                         `json:"reference,omitempty"`
	Customer    struct {
		ID string `json:"id,omitempty"`
	} `json:"customer,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// ENUM(recurring-debit,one-time)
type DirectDebitPaymentScheduleType string

// ENUM(fixed,variable)
type DebitType string

// ENUM(emandate,gsm,signed)
type MandateType string

type CreatedManadateDetails struct {
	MonoURL              string      `json:"mono_url,omitempty"`
	Type                 string      `json:"type,omitempty"`
	Method               string      `json:"method,omitempty"`
	MandateType          MandateType `json:"mandate_type,omitempty"`
	Amount               int64       `json:"amount,omitempty"`
	Description          string      `json:"description,omitempty"`
	Reference            string      `json:"reference,omitempty"`
	Customer             string      `json:"customer,omitempty"`
	AccountName          string      `json:"account_name,omitempty"`
	AccountNumber        string      `json:"account_number,omitempty"`
	Bank                 string      `json:"bank,omitempty"`
	RedirectURL          string      `json:"redirect_url,omitempty"`
	CreatedAt            time.Time   `json:"created_at,omitempty"`
	UpdatedAt            time.Time   `json:"updated_at,omitempty"`
	StartDate            string      `json:"start_date,omitempty"`
	EndDate              string      `json:"end_date,omitempty"`
	ID                   string      `json:"id,omitempty"`
	TransferDestinations []struct {
		BankName      string `json:"bank_name,omitempty"`
		AccountNumber string `json:"account_number,omitempty"`
		Icon          string `json:"icon,omitempty"`
	} `json:"transfer_destinations,omitempty"`
	OTPDestinations struct {
		Session string `json:"session,omitempty"`
		Methods []struct {
			Type  DirectDebitOTPMandateMethodType `json:"type,omitempty"`
			Value string                          `json:"value,omitempty"`
		} `json:"methods,omitempty"`
	} `json:"otp_destinations,omitempty"`
	Institution struct {
		BankCode string `json:"bank_code,omitempty"`
		NipCode  string `json:"nip_code,omitempty"`
		Name     string `json:"name,omitempty"`
	} `json:"institution,omitempty"`
}

type CreatedManadateDetailsResponse struct {
	BaseMonoResponse
	Data CreatedManadateDetails `json:"data,omitempty"`
}

func (d *DirectDebitService) Create(ctx context.Context,
	opts CreatedManadateOptions) (CreatedManadateDetails, error) {

	var resp CreatedManadateDetails

	body, err := ToReader(opts)
	if err != nil {
		return resp, err
	}

	req, err := d.client.newRequest(
		http.MethodPost, "/v3/payments/mandates", body)
	if err != nil {
		return resp, err
	}

	var res CreatedManadateDetailsResponse

	_, err = d.client.Do(ctx, req, &res)
	return res.Data, err
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
		BankCode string `json:"bank_code,omitempty"`
		NipCode  string `json:"nip_code,omitempty"`
		Name     string `json:"name,omitempty"`
	} `json:"institution,omitempty"`
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

	if hermes.IsStringEmpty(mandateID) {
		return resp, errors.New("please provide a valid mandate id")
	}

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

	if hermes.IsStringEmpty(mandateID) {
		return errors.New("please provide a valid mandate id")
	}

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

	if hermes.IsStringEmpty(mandateID) {
		return errors.New("please provide a valid mandate id")
	}

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

type DebitAccount struct {
	Status          string `json:"status"`
	Amount          int    `json:"amount"`
	Customer        string `json:"customer"`
	Mandate         string `json:"mandate"`
	ReferenceNumber string `json:"reference_number"`
	AccountDebited  struct {
		BankCode      string `json:"bank_code"`
		AccountName   string `json:"account_name"`
		AccountNumber string `json:"account_number"`
		BankName      string `json:"bank_name"`
	} `json:"account_debited"`
	Beneficiary struct {
		BankCode      string `json:"bank_code,omitempty"`
		AccountName   string `json:"account_name,omitempty"`
		AccountNumber string `json:"account_number,omitempty"`
		BankName      string `json:"bank_name,omitempty"`
	} `json:"beneficiary,omitempty"`
	Date time.Time `json:"date"`
}

type DebitAccountResponse struct {
	Data DebitAccount `json:"data,omitempty"`
	BaseMonoResponse
}

type DebitAccountOptions struct {
	Amount      int64  `json:"amount,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Narration   string `json:"narration,omitempty"`
	Beneficiary struct {
		Nuban   string `json:"nuban,omitempty"`
		NipCode string `json:"nip_code,omitempty"`
	} `json:"beneficiary,omitempty"`
}

func (d *DirectDebitService) DebitAccount(ctx context.Context,
	mandateID string, opts DebitAccountOptions) error {

	if hermes.IsStringEmpty(mandateID) {
		return errors.New("please provide a valid mandate id")
	}

	if hermes.IsStringEmpty(opts.Reference) {
		return errors.New("please provide your reference")
	}

	if hermes.IsStringEmpty(opts.Narration) {
		return errors.New("please provide your narration")
	}

	if opts.Amount < 10000 {
		return errors.New("you can only debit a minimum of 100 naira")
	}

	body, err := ToReader(opts)
	if err != nil {
		return err
	}

	req, err := d.client.newRequest(http.MethodPatch,
		fmt.Sprintf("/v3/payments/mandates/%s/debit", mandateID), body)
	if err != nil {
		return err
	}

	_, err = d.client.Do(ctx, req, nil)
	return err
}

func (d *DirectDebitService) Balance(ctx context.Context,
	mandateID string, amount int64) (int64, error) {

	if hermes.IsStringEmpty(mandateID) {
		return 0, errors.New("please provide a valid mandate id")
	}

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return 0, err
	}

	req, err := d.client.newRequest(http.MethodPatch,
		fmt.Sprintf("/v3/payments/mandates/%s/balance-inquiry/?amount=%d", mandateID, amount), body)
	if err != nil {
		return 0, err
	}

	_, err = d.client.Do(ctx, req, nil)
	return 0, err
}

type Banks []Bank

type Bank struct {
	Name        string `json:"name,omitempty"`
	BankCode    string `json:"bank_code,omitempty"`
	NIPCode     string `json:"nip_code,omitempty"`
	DirectDebit bool   `json:"direct_debit,omitempty"`
}

type BanksResponse struct {
	Data struct {
		Banks Banks `json:"banks,omitempty"`
	} `json:"data,omitempty"`

	BaseMonoResponse
}

func (d *DirectDebitService) Banks(ctx context.Context) (Banks, error) {

	var resp BanksResponse

	req, err := d.client.newRequest(http.MethodGet, "/v3/banks/list", strings.NewReader(""))
	if err != nil {
		return resp.Data.Banks, nil
	}

	_, err = d.client.Do(ctx, req, &resp)
	return resp.Data.Banks, err
}

// ENUM(phone_number,email)
type DirectDebitOTPMandateMethodType string

type SetOTPMethodOptions struct {
	Method  DirectDebitOTPMandateMethodType `json:"method,omitempty"`
	Session string                          `json:"session,omitempty"`
}

func (d *DirectDebitService) SetOTPMethod(ctx context.Context,
	opts SetOTPMethodOptions) error {

	if hermes.IsStringEmpty(opts.Session) {
		return errors.New("please provide a valid session id")
	}

	if !opts.Method.IsValid() {
		return errors.New("invalid otp method selected")
	}

	body, err := ToReader(opts)
	if err != nil {
		return err
	}

	req, err := d.client.newRequest(http.MethodPost,
		"/v3/payments/mandates/verify/otp", body)
	if err != nil {
		return err
	}

	_, err = d.client.Do(ctx, req, nil)
	return err
}

type VerifyOTPOptions struct {
	Session string `json:"session,omitempty"`
	OTP     string `json:"otp,omitempty"`
}

func (d *DirectDebitService) VerifyOTP(ctx context.Context,
	opts VerifyOTPOptions) error {

	if hermes.IsStringEmpty(opts.Session) {
		return errors.New("please provide a valid session id")
	}

	if hermes.IsStringEmpty(opts.OTP) {
		return errors.New("please provide an OTP")
	}

	body, err := ToReader(opts)
	if err != nil {
		return err
	}

	req, err := d.client.newRequest(http.MethodPost,
		"/v3/payments/mandates/verify/otp", body)
	if err != nil {
		return err
	}

	_, err = d.client.Do(ctx, req, nil)
	return err
}
