# POST /api/revoke

This endpoint enables users to revoke their refresh token, preventing the refresh token from being used to generate new access tokens in the future.

## Request

`POST /api/revoke`

### Headers

- `Authorization`: Bearer `<refresh-token>` - The refresh token that the user wishes to revoke.

### Example Header


```request 
Authorization: Bearer <refresh-token>
```

## Response

### Success Response

- **Code**: `200 OK`
- **Content**: JSON object indicating that the token has been successfully revoked.

### Example Success Response

```json
{   
	"Success": "Token revoked"
}
```

### Error Responses

- **Unauthorized (No JWT Found)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Couldn't find JWT"}`
- **Unauthorized (Invalid JWT)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Couldn't validate JWT"}`
- **Unauthorized (Not a Refresh Token)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Provide a refresh token"}`
- **Conflict (Token Already Revoked)**:
    - **Code**: `409 CONFLICT`
    - **Content**: `{"error": "Token already revoked"}`
- **Internal Server Error (Revoking Token)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Something went wrong"}`

## Examples

### Revoking a Refresh Token

To revoke a refresh token, thereby preventing it from being used again, you would make the following request:

```bash
curl -X POST 'http://localhost:8080/api/revoke' \
-H 'Authorization: Bearer <refresh-token>'
```

Replace `<refresh-token>` with the refresh token that needs to be revoked.

## Notes

- Revoking a refresh token is a critical step in ensuring that a user's session can be securely ended. Once revoked, the refresh token cannot be used to generate new access tokens.
- The endpoint requires a valid refresh token to be provided in the `Authorization` header. It checks whether the token is valid, has not been revoked already, and belongs to the refresh token issuer.