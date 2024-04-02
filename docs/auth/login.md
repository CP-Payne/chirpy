# POST /api/login

This endpoint authenticates a user based on their email and password, issuing JWT access and refresh tokens upon successful authentication.

## Request

`POST /api/login`

### Headers

- `Content-Type`: `application/json` - Indicates that the body of the request is JSON.

### Body Parameters

- `email`: A string containing the user's email address.
- `password`: A string containing the user's password.

### Example Request Body

```json
{   
	"email": "user@example.com",   
	"password": "userpassword" 
}
```

## Response

### Success Response

- **Code**: `200 OK`
- **Content**: JSON object containing the authenticated user's details along with the issued JWT access and refresh tokens.

### Example Success Response

```json
{   
	"ID": 1,   
	"Email": "user@example.com",   
	"IsChirpyRed": false,   
	"Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",   
	"RefreshToken": "def50200..." 
}
```

### Error Responses

- **Internal Server Error (Decoding Parameters)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Couldn't decode parameters"}`
- **Internal Server Error (General)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Something went wrong"}`
- **Unauthorized (Incorrect Credentials)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Incorrect credentials"}`

## Examples

### Authenticate a User

To authenticate a user and receive JWT tokens, you would make the following request:

```bash
curl -X POST 'http://localhost:8080/api/login' \
-H 'Content-Type: application/json' \
-d '{"email": "user@example.com", "password": "userpassword"}'
```

## Notes

- This endpoint validates the user's credentials against the stored records. If the credentials are correct, it issues an access token and a refresh token.
- The access token is used for authenticating subsequent requests by the user. The refresh token can be used to obtain a new access token when the current one expires.
- The password is checked against its hashed version stored in the database.
- Tokens are to be added to the `Authorization` header in subsequent requests in the format `Bearer <token>`.