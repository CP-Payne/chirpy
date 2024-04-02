# DELETE /api/chirps/{chirp_id}

This endpoint is responsible for deleting a specific chirp by its unique identifier, provided that the requester is authenticated and is the author of the chirp.

## Request

`DELETE /api/chirps/{chirp_id}`

### Path Parameters

- `chirp_id`: The unique identifier for the chirp to be deleted.

### Headers
```request
Authorization: Bearer <jwt-token>
```
Replace `<jwt-token>` with the JWT token received upon authentication.

## Response

### Success Response

- **Code**: `200 OK`
- **Content**: A confirmation message indicating that the chirp has been deleted successfully.

### Error Responses

- **Invalid Chirp ID**:    
    - **Code**: `400 BAD REQUEST`
    - **Content**: `{"error": "Invalid chirp id"}`
- **Unauthorized (No JWT Found)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Couldn't find JWT"}`
- **Unauthorized (Invalid JWT)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "Couldn't validate JWT"}`
- **Forbidden (Resource Ownership)**:
    - **Code**: `403 FORBIDDEN`
    - **Content**: `{"error": "You do not own this resource"}`
- **Not Found (Chirp Doesn't Exist)**:
    - **Code**: `404 NOT FOUND`
    - **Content**: `{"error": "Chirp not found"}`
- **Internal Server Error**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Something went wrong"}`

## Examples

### Delete a Chirp by ID

To delete a chirp with ID `123`, an authorized user would make the following request:

```bash
curl -X DELETE 'http://localhost:8080/api/chirps/123' \
-H 'Authorization: Bearer <jwt-token>'
```


Replace `<jwt-token>` with the JWT token of the user.

## Notes

- A valid JWT token must be included in the `Authorization` header.
- The `chirp_id` must correspond to a chirp that exists, and the authenticated user must be the author of the chirp. If these conditions are not met, appropriate error messages and status codes are returned.
- If the chirp is successfully deleted, a `200 OK` response is returned along with a message indicating success.
- The API response for a successful deletion does not include a content body, as the resource is no longer available. The HTTP status text for `200 OK` is used as the response to signify successful deletion.