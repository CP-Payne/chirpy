# POST /api/refresh

This endpoint is used to refresh the JWT access token using a valid refresh token. It is part of a secure authentication system that allows users to obtain a new access token without re-entering their credentials.

## Request

`POST /api/refresh`

### Headers

- `Authorization`: Bearer `<refresh-token>` - The refresh token received during the initial authentication.

### Example Header

```request
Authorization: Bearer <refresh-token>
```

## Response

### Success Response

- **Code**: `200 OK`
- **Content**: JSON object containing the newly issued JWT access token.

### Example Success Response

```json
{   
	"Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." 
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
- **Unauthorized (Revoked Token)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "please log in again"}`
- **Internal Server Error (Parsing User ID)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Couldn't parse user ID from token"}`
- **Internal Server Error (Token Creation)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Something went wrong"}`

## Examples

### Refreshing an Access Token

To refresh an access token using a refresh token, you would make the following request:



```bash
curl -X POST 'http://localhost:8080/api/refresh' \
-H 'Authorization: Bearer <refresh-token>'
```

Replace `<refresh-token>` with the refresh token.

## Notes

- Refresh tokens are typically long-lived and used to request new access tokens once the access token has expired or needs to be refreshed.
- The endpoint validates the provided refresh token, checks if it is not revoked, and issues a new access token.
- Ensure the refresh token is validated against its intended use (identified by the issuer) and verify that it corresponds to an existing, non-revoked token in the database.