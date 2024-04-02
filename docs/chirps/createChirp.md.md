
# POST /api/chirps

This endpoint allows an authenticated user to post a new chirp.

## Request

`POST /api/chirps`

### Headers
```request
Authorization: Bearer <jwt-token>
```
 

Replace `<jwt-token>` with the JWT token received upon authentication.

### Body Parameters

- `body`: A string containing the chirp's content. Must be 140 characters or less.

### Example Body

```json
{
  "body": "Excited to be using Chirpy!"
}
````

## Response

### Success Response

- **Code**: `201 CREATED`
- **Content**:

```json
{   
	"id": 1,   
	"body": "Excited to be using Chirpy!",   
	"author_id": 42 
}
```



### Error Responses

- **Unauthorized (No Token Found)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Couldn't find JWT"}`
- **Unauthorized (Invalid Token)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Couldn't validate JWT"}`
- **Internal Server Error (User ID Parse Error)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Couldn't parse user ID from token"}`
- **Internal Server Error (Body Decode Error)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Couldn't decode parameters"}`
- **Bad Request (Chirp Too Long)**:
    - **Code**: `400 BAD REQUEST`
    - **Content**: `{"error": "Chirp is too long"}`

## Examples

### Post a New Chirp

```bash
curl -X POST 'http://localhost:8080/api/chirps' \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer <jwt-token>' \
-d '{"body": "Excited to be using Chirpy!"}'
```

Replace `<jwt-token>` with the JWT token.

## Notes

- The chirp `body` cannot exceed 140 characters. If it does, the API will respond with a `400 BAD REQUEST`.
- The JWT token must be valid and contain a `subject` that is the user ID and an `issuer` that matches `chirpy-access`. Otherwise, a `401 UNAUTHORIZED` response will be returned.
- The `body` of the chirp will be cleaned of any predefined 'bad words' before saving.