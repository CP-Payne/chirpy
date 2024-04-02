# PUT /api/users

This endpoint allows authenticated users to update their account information, including their email and password.

## Request

`PUT /api/users`

### Headers

- `Content-Type`: `application/json` - Indicates that the body of the request is JSON.
- `Authorization`: Bearer `<jwt-token>` - The JWT token for authenticating the user.

### Body Parameters

- `email`: A string containing the user's new email. This must be a unique email that is not already associated with another user account.
- `password`: A string containing the user's new password. This will be hashed before storage for security.

### Example Request Body

```json
{   
	"email": "updatedemail@example.com",   
	"password": "newsecurepassword123" 
}
```

## Response

### Success Response

- **Code**: `200 OK`
- **Content**: JSON object containing the updated user's details, excluding the password for security reasons.

### Example Success Response

```json
{  
	"ID": 1,   
	"Email": "updatedemail@example.com"
}
```

### Error Responses

- **Unauthorized (No JWT Found)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Couldn't find JWT"}`
- **Unauthorized (Invalid JWT)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Couldn't validate JWT"}`
- **Conflict (Email Already Exists)**:
    - **Code**: `409 CONFLICT`
    - **Content**: `{"error": "Email already exist"}`
- **Internal Server Error (Decoding Parameters)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Couldn't decode parameters"}`
- **Internal Server Error (General)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Couldn't update user information"}`

## Examples

### Update User Information

To update a user's email and password, you would make the following request:

```bash
curl -X PUT 'http://localhost:8080/api/users' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer <jwt-token>' \
-d '{"email": "updatedemail@example.com", "password": "newsecurepassword123"}'
```

Replace `<jwt-token>` with the JWT token of the user.

## Notes

- A valid JWT token must be included in the `Authorization` header to authenticate the request.
- The user can update their email and/or password. If the email is already in use by another account, the API will return a `409 CONFLICT` error.
- Passwords are hashed before being stored in the database to ensure the security of user credentials.