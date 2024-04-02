# POST /api/users

This endpoint is responsible for creating a new user account in the Chirpy application. It requires the client to send user registration details, including an email and password.

## Request

`POST /api/users`

### Headers

```request
Content-Type: application/json
```
- Indicates that the body of the request is JSON.

### Body Parameters

- `email`: A string containing the user's email. This must be a unique email that is not already associated with an existing user account.
- `password`: A string containing the user's password. This will be hashed before storage for security.

### Example Request Body

```json
{
	"email": "newuser@example.com",   
	"password": "securepassword123" 
}
```

## Response

### Success Response

- **Code**: `201 CREATED`
- **Content**: JSON object containing the newly created user's details, excluding the password for security reasons.

### Example Success Response

```json
{   
	"ID": 1,   
	"Email": "newuser@example.com",   
	"IsChirpyRed": false
}
```

### Error Responses

- **Internal Server Error (Decoding Parameters)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Couldn't decode parameters"}`
- **Internal Server Error (General)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Something went wrong"}`
- **Conflict (Email Already Exists)**:
    - **Code**: `409 CONFLICT`
    - **Content**: `{"error": "Email already exist"}`

## Examples

### Create a New User

To create a new user account, you would make the following request:

```bash
curl -X POST 'http://localhost:8080/api/users' \
-H 'Content-Type: application/json' \
-d '{"email": "newuser@example.com", "password": "securepassword123"}'
```


## Notes

- The API performs a check to confirm the uniqueness of the email and will return a `409 CONFLICT` error if the email is already in use.
- Passwords are hashed before being stored in the database to ensure the security of user credentials.
- The response includes basic user information. Sensitive information like the user's password is not included in the response for security reasons.