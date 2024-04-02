# GET /api/chirps/{chirp_id}

This endpoint retrieves a specific chirp by its unique identifier.

## Request

`GET /api/chirps/{chirp_id}`

### Path Parameters

- `chirp_id`: The unique identifier of the chirp you want to retrieve.

## Response

### Success Response

- **Code**: `200 OK`
- **Content**: A JSON object representing the chirp with the specified ID.

### Example Response

```json
{   
	"id": 1,
	"body": "Here is a chirp message",   
	"author_id": 42,   
	"created_at": "2024-04-02T14:00:00Z" 
}
```

### Error Responses

- **Invalid Chirp ID**:
    - **Code**: `400 BAD REQUEST`
    - **Content**: `{"error": "Invalid chirp id"}`
- **Chirp Not Found**:
    - **Code**: `404 NOT FOUND`
    - **Content**: `{"error": "Not Found"}`
- **Internal Server Error**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Something went wrong"}`

## Examples

### Get a Chirp by ID

To retrieve a chirp with ID `123`, you would make the following request:
```bash
curl -X GET 'http://localhost:8080/api/chirps/123' -H 'Accept: application/json'
```


## Notes

- The `chirp_id` must be an integer. If a non-integer is provided, the API will respond with a `400 BAD REQUEST`.
- If there is no chirp matching the provided ID, the API will respond with a `404 NOT FOUND`.