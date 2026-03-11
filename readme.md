## Requirements

- Go 1.26
- Mysql 

## How to run

1. migrate

```bash
go run main.go migrate up
```

2. run

```bash
go run main.go
```

3. set up environment variable

```bash
cp .env.example .env
```

## How to test

exposed wallet id can seen from detail wallet and list wallet, so you can use it for next request.
after you make action to wallet, you can hit api confirmation transaction to confirm it, if you don't confirm it, the transaction will not be processed.
you can use tool like jmeter or postman to test it.

1. create owner

```bash
curl -X POST http://localhost:8080/api/owner \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
  }'
```

2. create wallet

```bash
curl -X POST http://localhost:8080/api/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "owner_id": <ownerId>,
    "currency": "USD"
  }'
```

3. top up wallet

```bash
curl -X POST http://localhost:8080/api/wallet/<walletId>/top-up \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "100000"
  }'
```

4. payment wallet

```bash
curl -X POST http://localhost:8080/api/wallet/<walletId>/payment \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "100000"
  }'
```

5. transfer wallet

```bash
curl -X POST http://localhost:8080/api/wallet/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "from_wallet_id": <walletId>,
    "to_wallet_id": <walletId>,
    "amount": "100000"
  }'
```

6. suspend wallet

```bash
curl -X POST http://localhost:8080/api/wallet/<walletId>/suspend
```

7. detail wallet

```bash
curl -X GET http://localhost:8080/api/wallet/<walletId>
```

8. list wallet

```bash
curl -X GET http://localhost:8080/api/owner/<ownerId>/wallets
```

9. confirmation transaction

```bash
curl -X POST http://localhost:8080/api/transaction/confirmation \
  -H "Content-Type: application/json" \
  -d '{
    "trx_id": <transactionId>,
  }'
```