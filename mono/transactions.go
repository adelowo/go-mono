package mono

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	// We do not use this internally but can be helpful for people who want to parse the date
	// correctly
	DateFormat = "2006-01-02"
)

type TransactionService service

// ENUM(debit,credit)
type TransactionType string

type TransactionsFilter struct {
	Paginate   bool            `url:"paginate,omitempty"`
	Start, End string          `url:"start,omitempty"`
	Narration  string          `url:"narration,omitempty"`
	Type       TransactionType `url:"type,omitempty"`
	Limit      int             `url:"limit,omitempty"`
}

type Transaction struct {
	ID        string    `json:"id"`
	Narration string    `json:"narration"`
	Amount    int       `json:"amount"`
	Type      string    `json:"type"`
	Balance   int       `json:"balance"`
	Date      time.Time `json:"date"`
	Category  string    `json:"category"`
}

type TransactionMetadata struct {
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	Previous any    `json:"previous"`
	Next     string `json:"next"`
}

type TransactionsResponse struct {
	BaseMonoResponse
	Data []Transaction       `json:"data"`
	Meta TransactionMetadata `json:"meta"`
}

func (t *TransactionService) All(ctx context.Context,
	accountID string,
	opts TransactionsFilter) ([]Transaction, TransactionMetadata, error) {

	var resp []Transaction
	var metadata TransactionMetadata

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return resp, metadata, nil
	}

	v, err := query.Values(opts)
	if err != nil {
		return resp, metadata, err
	}

	req, err := t.client.newRequest(http.MethodGet,
		fmt.Sprintf("/accounts/%s/transactions?%s", accountID, v.Encode()),
		body)
	if err != nil {
		return resp, metadata, nil
	}

	var response TransactionsResponse
	_, err = t.client.Do(ctx, req, response)
	return response.Data, response.Meta, err
}
