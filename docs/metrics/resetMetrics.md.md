# POST /api/reset

This endpoint resets certain metrics within the Chirpy application, such as the number of hits to the file server.

## Request

`POST /api/reset`

## Response

### Success Response

- **Code**: `200 OK`
- **Content**: A plain text response indicating the operation was successful.

### Example Response

`OK`

## Examples

### Reset Application Metrics

To reset the file server hit counter, a user would make the following request:

```bash
curl -X POST 'http://localhost:8080/api/reset'
```


## Notes

- This is currently unauthenticated