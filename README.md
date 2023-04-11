This is a simple bank RESTful API built using Golang.

## Installation

1. Clone the repository: git clone https://github.com/flukis/go-simplebank.git
2. Navigate into the directory: cd go-simplebank
3. Install dependencies: go mod download

## Usage

1. Start the server: make dev
2. Use a RESTful API client (such as Postman) to interact with the API endpoints.

## API Endpoints

The API provides the following endpoints:

- POST /api/auth/signup: create a new user account
- POST /api/auth/login: authenticate and login a user
- GET /api/accounts/:id: retrieve an account by its ID
- POST /api/accounts: create a new account
- GET /api/accounts: retrieve a list of all accounts
- POST /api/transfers: create a new transfer between two accounts
- GET /api/transfers/:id: retrieve a transfer by its ID
- GET /api/transfers: retrieve a list of all transfers

## Contributing

Contributions are welcome! If you have any suggestions or find any bugs, please open an issue or submit a pull request.
