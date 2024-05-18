# diploma

## Requirements

- docker 26.0.0
- docker-compose 1.29.2

## About
The application is a personal financial manager, divided into microservices
for more efficient operation. It allows users to manage their finances: view
transaction information, add transactions, set goals, add accounts, and
manage expense and income categories.

![piglet.drawio.svg](piglet.drawio.svg)

All microservices are implemented using gRPC to ensure efficient
communication between them. Various tools and technologies are also used,
such as Docker for containerization, protobuf for defining messages and
services, and various libraries such as validator, pgx, envconfig and others.

## Docs
You can view the documentation in Swagger format [here](https://github.com/missbulochka/protos/tree/main/openapiv2).

## Run
You can run dockerized application via docker compose:
```bash
docker compose up --build
```

## Using
There are three entities that a user can manage: *accounts*, *transactions*, and *categories*.

### Accounts
Accounts are entities to which a person associates all his transactions.
Money is transferred either from them or to them.
They can be of two types: cash account and goal.

By *account*, we mean card or cash. A default account is automatically
created in the application.

The *goal* refers to the amount of money that the user wants to accumulate.
The user enters the amount and the date by which it needs to be accumulated,
and the application calculates the monthly payment.

#### Create

To create an account the user needs to make a request in terminal:
```bash
curl -X 'POST' \
  'http://piglet-gateway:8083/piglet/bills/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --data '{}'
```

The curl body data depends on the account type:
- account:
  ```bash
  --data '{
    "billType": true,
    "billName": "account",
  }'
  ```
- goal:
    ```bash
  --data '{
    "billType": false,
    "billName": "goal",
    "goalSum": 500000,
    "date": "2024-05-18"
  }'
  ```

The command will return the account body, filling it with the values written
to the database.

#### Update

To update a transaction, the rules are similar to creation, but the user needs
to know the account `id`. Replace `{id}` with the account ID to be
deleted.

To update an account the user needs to make a request in terminal:
```bash
curl -X 'PUT' \
  'http://piglet-gateway:8083/piglet/bills/{id}' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --data '{}'
```

#### Get all accounts
To get all the accounts:
```bash
curl -X 'GET' \
  'http://piglet-gateway:8083/piglet/bills/accounts' \
  -H 'accept: application/json'
```

To get all the goals:
```bash
curl -X 'GET' \
  'http://piglet-gateway:8083/piglet/bills/goals' \
  -H 'accept: application/json'
```

#### Get an account
To get an account:
```bash
curl -X 'GET' \
  'http://piglet-gateway:8083/piglet/bills/{id}' \
  -H 'accept: application/json'
```

Replace `{id}` with the account ID to be deleted.

#### Delete the account
To delete an account:
```bash
curl -X 'DELETE' \
  'http://piglet-gateway:8083/piglet/bills/{id}' \
  -H 'accept: application/json'
```

Replace `{id}` with the account ID to be deleted.

### Transactions

Transactions are records of money movements, such as income, expenses, debts,
and transfers.

#### Create
To create a transaction, the user needs to make a request in the terminal:

```bash
curl -X 'POST' \
  'http://piglet-gateway:8083/piglet/transactions/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --data '{}'
```

The curl body data depends on the transaction type:
- income:
  ```bash
  --data '{
    "transType": 1,
    "sum": 500,
    "date": "2024-05-18",
    "idCategory": "00000000-0000-0000-0000-000000000001",
    "idBillTo": "00000000-0000-0000-0000-000000000001",
    "repeat": true
  }'
  ```
  Parameters `idCategory` and `idBillTo` are required for this transaction type.
  `Repeat` is an optional (if undefined the transaction will be non-repeating).
- expense:
    ```bash
  --data '{
    "transType": 2,
    "sum": 500,
    "idCategory": "00000000-0000-0000-0000-000000000000",
    "idBillFrom": "00000000-0000-0000-0000-000000000001",
    "person": "Oleg"
  }'
  ```
  Parameters `idCategory` and `idBillFrom` are required for this transaction type.
  `Repeat` is an optional (if undefined the transaction will be non-repeating).
- debt:
  ```bash
  --data '{
    "transType": 3,
    "sum": 500,
    "debtType": true,
    "idBillTo": "",
    "idBillFrom": "00000000-0000-0000-0000-000000000001"
  }'
  ```
  Parameters `debtType`, `idBillTo` and `idBillFrom` required for each debt type.
  If you are the creditor `debtType` sets true; otherwise, sets false.
- transfer:
  ```bash
  --data '{
    "transType": 4,
    "sum": 500,
    "idBillTo": "e3662e81-5c56-4159-bbd2-0f989759c305",
    "idBillFrom": "00000000-0000-0000-0000-000000000001"
  }'  
  ```
  
Parameters `transType`, `sum` required for **all** transactions.
Parameters `transDate`, `comment`, `person` are optional (if undefined,
the date becomes equal to the transaction creation date, and the comment and
person remains empty).

#### Update

To update a transaction, the rules are similar to creation, but the user needs
to know the transaction `id`. Replace `{id}` with the transaction ID.

```bash
curl -X 'PUT' \
  'http://piglet-gateway:8083/piglet/transactions/{id}' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --data '{}'
```

#### Delete
To delete a transaction, the user needs to make a request in the terminal:
```bash
curl -X 'DELETE' \
  'http://piglet-gateway:8083/piglet/transactions/{id}' \
  -H 'accept: application/json'
```
Replace {id} with the transaction ID to be deleted.

### Get
To get details of a specific transaction, the user needs to make a request in
the terminal:

```bash
curl -X 'GET' \
  'http://piglet-gateway:8083/piglet/transactions/{id}' \
  -H 'accept: application/json'
```
Replace {id} with the transaction ID.

To get all transactions, the user needs to make a request in the terminal:

```bash
curl -X 'GET' \
  'http://piglet-gateway:8083/piglet/transactions/all' \
  -H 'accept: application/json'
```

### Categories

Categories are a transaction property used to classify transactions into
different types, such as expenses or incomes and mandatory or optional.
The system automatically creates default categories of expenses and income.
This will include all transactions for which the user has not specified
a category.

#### Add
To add a category, the user needs to make a request in the terminal:

```bash
curl -X 'POST' \
  'http://piglet-gateway:8083/piglet/categories/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --data '{
    "type": true,
    "name": "Groceries",
    "mandatory": true
  }'
```

Parameters:
- `type` true for expense, false for income
- `mandatory` true for mandatory transaction

#### Update
To update a category, the user needs to know the category `id` and make
a request in the terminal:

```bash
curl -X 'PUT' \
  'http://piglet-gateway:8083/piglet/categories/{id}' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  --data '{
    "type": false,
    "name": "Salary",
    "mandatory": true
  }'
```
Replace {id} with the category ID, and the category type, name, and mandatory
flag as needed.

#### Delete
To delete a category, the user needs to make a request in the terminal:

```bash
curl -X 'DELETE' \
  'http://piglet-gateway:8083/piglet/categories/{id}' \
  -H 'accept: application/json'
```
Replace {id} with the category ID to be deleted.

#### Get
To get details of a specific category, the user needs to make a request in
the terminal:

```bash
curl -X 'GET' \
    'http://piglet-gateway:8083/piglet/categories/{id}' \
    -H 'accept: application/json'
```
Replace {id} with the category ID.

To get all categories, the user needs to make a request in the terminal:

```bash
curl -X 'GET' \
    'http://piglet-gateway:8083/piglet/categories/all' \
    -H 'accept: application/json'
```
