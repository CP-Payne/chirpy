
# Chirpy API

The Chirpy API, a platform similar to Twitter where users can share messages, or "chirps". The API allows you to interact programmatically with our main features such as posting messages and getting lists of chirps. This project is a learning endeavour and is not intended for production use.


## Overview
Chirpy API is built in Go, as part of the Learn Web Servers course in `boot.dev`.

## API Endpoints

Below is a list of available endpoints, grouped by functionality. Detailed documentation and examples for each endpoint is available in the `/doc` folder, which you can access by clicking on the endpoint name.

### Chirps
- [`POST /api/chirps`](./docs/chirps/createChirp.md) - Post a new chirp.
- [`GET /api/chirps`](./docs/chirps/getChirps.md) - Retrieve a list of all chirps.
- [`GET /api/chirps/{chirp_id}`](./docs/chirps/getChirpByID.md) - Get a single chirp by its ID.
- [`DELETE /api/chirps/{chirp_id}`](./docs/chirps/deleteChirp.md) - Delete a chirp by its ID.

### Metrics

- [`GET /admin/metrics`](./docs/metrics/getMetrics.md) - Obtain the number of visits to the home page.
- [`POST /api/reset`](./docs/metrics/resetMetrics.md) - Reset the number of visits to home page

### Users

- [`POST /api/users`](./docs/users/createUser.md) - Register a new user.
- [`PUT /api/users`](./docs/users/updateUser.md) - Update an existing user's profile.

### Auth

- [`POST /api/login`](./docs/auth/login.md) - Authenticate a user and receive a token.
- [`POST /api/refresh`](./docs/auth/refreshToken.md) - Refresh the authentication token.
- [`POST /api/revoke`](./docs/auth/revokeRefreshToken.md) - Revoke a refresh token.

### Webhooks

- [`POST /api/polka/webhooks`](./docs/webhooks/upgradeUser.md) - Upgrade the user role through a webhook.

### Other
- [`GET /app`](./docs/other/servingFiles.md) - A simple endpoint to display the homepage content.

## Current Data Storage Implementation

The application currently uses a file-based storage system for managing data. All application data, including users, chirps, and authentication tokens, are stored in a `database.json` file located at the root of the project directory.  An example of the file structure can be found here: [Database file example](./docs/databaseExample/databaseExample.md)

## Installation

To get the Chirpy API up and running on your local machine, follow these steps:

### Prerequisites

- Ensure you have Go installed on your system. The Chirpy API requires Go version 1.22 or later.
- Git should be installed to clone the repository.

### Cloning the Repository

First, clone the Chirpy API repository to your local machine using Git. Open your terminal, and run the following command:

```bash
git clone https://github.com/CP-Payne/chirpy
cd chirpy
```
### Installing Dependencies
The project uses `go.mod` for managing dependencies. To install all the necessary dependencies, run:

```bash
go mod download
```

This command will download and install all the required packages and dependencies for the Chirpy API.

### Environment Setup

Before running the Chirpy API, you need to set up the required environment variables. Create a `.env` file in the root directory of the project with the following content:

```plaintext
JWT_SECRET=YOUR_JWT_SECRET_HERE
````
Replace `YOUR_JWT_SECRET_HERE` with your own secret key for JWT authentication. This is crucial for generating and validating JWT tokens securely.

**Note:** The `WEBHOOK_API` variable is used for a specific function in the API that integrates with an external webhook service provided by boot.dev. This feature requires a subscription with boot.dev to obtain an API key. Without this subscription, you can still run and use the majority of the API; however, the webhook-related functionality will not be available.

#### Running the Server Without the Webhook API

If you don't have access to the `WEBHOOK_API` key, you can still run the server. The API will function normally for all endpoints that don't require this integration. Just keep in mind that any features relying on the webhook service will not be operational. (Currently only one endpoint)

#### Running the Server With the Webhook API

If you have a subscription with boot.dev and have obtained your API key, you can add it to the `.env` file like so:

```plaintext
WEBHOOK_API=YOUR_API_KEY_HERE
```
Replace `YOUR_API_KEY_HERE` with your actual API key from boot.dev. Once set, the full functionality of the API, including webhook-related features, will be available.

Save the `.env` file and follow the instructions in the `Running the Server` section to start your API
### Running the Server

To start the server, you simply need to run the `main.go` file. From the root directory of the project, execute:

```bash
go run main.go
```


After running the above command, the API server should be up and listening for requests on the port `8080`, which you can then access via `localhost`.

### Testing Endpoints

You can test the endpoints using tools like `curl` or Postman. For example, to get all chirps, you might use:

```bash
curl http://localhost:port/api/chirps
```

Replace `port` with the actual port number your server is listening on.

## Future Enhancements

- Create a frontend for the API
- Implement a database (Postgres, MySQL, etc) instead of using a file to store the information
- Add feature to comment on chirps
- Add profile management