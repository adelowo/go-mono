# Go-Mono

go-mono is a Go client library for accessing the Mono v2 API.

> [!WARNING]
> This is not feature complete nor does it cover all available APIs.
> I am building to use for <https://tewotewo.com> . So I will prioritize
> features needed over there. But feel free to open an issue or PR if you need
> support for a specific endpoint

## Usage

```go

 client, err := mono.New(mono.WithAPISecret("test_sk_ddd"))
 if err != nil {
  log.Fatal(err)
 }

 // or other methods as you use
 err = client.Account.Unlink(context.Background(), "id")
 if err != nil {
  log.Fatal(err)
 }

```

### Using your own HTTP Client

Sometimes, you need to reuse an http client instead of letting the library
create and manage one for you. Sometimes, this might also be because you
want to trace the request amongst others.

```go

httpClient := http.Client{
 Transport: otelhttp.NewTransport(http.DefaultTransport),
}

client, err := mono.New(
  mono.WithAPISecret("test_sk_ddd"),
  mono.WithHTTPClient(httpClient),)


```

## Features

- [x] Transactions
  - [x] List transactions
- [x] Authorisation
  - [x] Exchange token
  - [x] Reauthorise account
- [x] Inflow/Outflow
  - [x] Credits
  - [x] Debits
- [x] Accounts
  - [x] Details
  - [x] Unlink account
  - [ ] Identity
  - [x] Balance
  - [ ] Income
  - [ ] Income records
  - [ ] All accounts linked to business
  - [ ] Creditworthiness
- [ ] Statement
  - [ ] Fetch statement
  - [ ] Poll statement
- [ ] Investments
  - [ ] Assets
  - [ ] Earnings
- [ ] Data enrichment
  - [ ] Transaction Categorisation
  - [ ] Statement insights
