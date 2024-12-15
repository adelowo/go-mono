package mono

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ayinke-llc/hermes"
)

type CustomerService service

type CustomerDetails struct {
	Type               string `json:"type,omitempty"`
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	FirstName          string `json:"first_name,omitempty"`
	LastName           string `json:"last_name,omitempty"`
	Email              string `json:"email,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Address            string `json:"address,omitempty"`
	IdentificationNo   string `json:"identification_no,omitempty"`
	IdentificationType string `json:"identification_type,omitempty"`
	BVN                string `json:"bvn,omitempty"`
}

type CustomerDetailsResponse struct {
	BaseMonoResponse
	Data CustomerDetails `json:"data"`
}

func (c *CustomerService) Details(ctx context.Context,
	accountID string) (CustomerDetails, error) {

	var resp CustomerDetails

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return resp, err
	}

	req, err := c.client.newRequest(http.MethodGet,
		fmt.Sprintf("/v2/customers/%s", accountID),
		body)
	if err != nil {
		return resp, err
	}

	var customer CustomerDetailsResponse

	_, err = c.client.Do(ctx, req, &customer)
	if err != nil {
		return resp, err
	}

	return customer.Data, nil
}

// ENUM(individual,business)
type CustomerType string

type CreateCustomerOptions struct {
	Email     string       `json:"email,omitempty"`
	Type      CustomerType `json:"type,omitempty"`
	FirstName string       `json:"first_name,omitempty"`
	LastName  string       `json:"last_name,omitempty"`
	Address   string       `json:"address,omitempty"`
	Phone     string       `json:"phone,omitempty"`
	Identity  struct {
		Type   string `json:"type,omitempty"`
		Number string `json:"number,omitempty"`
	} `json:"identity,omitempty"`
}

func (c *CustomerService) Create(ctx context.Context,
	opts CreateCustomerOptions) (CustomerDetails, error) {

	var resp CustomerDetails

	body, err := ToReader(opts)
	if err != nil {
		return resp, err
	}

	req, err := c.client.newRequest(http.MethodPatch, "/v2/customers", body)
	if err != nil {
		return resp, err
	}

	var customer CustomerDetailsResponse

	_, err = c.client.Do(ctx, req, &customer)
	return customer.Data, err
}

type UpdateCustomerOptions struct {
	Phone    string `json:"phone,omitempty"`
	Address  string `json:"address,omitempty"`
	Identity struct {
		Type   string `json:"type,omitempty"`
		Number string `json:"number,omitempty"`
	} `json:"identity,omitempty"`
}

func (c *CustomerService) Update(ctx context.Context,
	customerID string, opts *UpdateCustomerOptions) error {

	body, err := ToReader(opts)
	if err != nil {
		return err
	}

	req, err := c.client.newRequest(http.MethodPatch,
		fmt.Sprintf("/v2/customers/%s", customerID),
		body)
	if err != nil {
		return err
	}

	var customer CustomerDetailsResponse

	_, err = c.client.Do(ctx, req, &customer)
	return err
}

func (a *CustomerService) Delete(ctx context.Context,
	customerID string) error {

	if hermes.IsStringEmpty(customerID) {
		return errors.New("please provide a valid customerID")
	}

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return err
	}

	req, err := a.client.newRequest(http.MethodDelete,
		fmt.Sprintf("/customers/%s", customerID), body)
	if err != nil {
		return err
	}

	_, err = a.client.Do(ctx, req, nil)
	return err
}
