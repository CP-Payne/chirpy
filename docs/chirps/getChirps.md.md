
# GET /api/chirps

This endpoint retrieves a list of chirps. A chirp is a short message or status update similar to a tweet on Twitter.

## Request

`GET /api/chirps`

### Query Parameters

- `author_id` (optional): Filters the chirps to those created by a specific author, identified by their unique author ID.
- `sort` (optional): Specifies the order in which the chirps should be returned. (`asc` or `desc`)

## Response

Upon a successful request, the API responds with a JSON array of chirp objects.

### Success Response

- **Code**: `200 OK`
- **Content**:

```json
[
  {
    "id": 1,
    "body": "This is a chirp message",
    "author_id": 10
  },
  {
    "id": 2,
    "body": "Another chirp message",
    "author_id": 15
  }
  // ... additional chirp objects
]
````

Each `ChirpResponse` object contains:

- `id`: An integer representing the unique identifier for the chirp.
- `body`: A string containing the chirp's message content.
- `author_id`: An integer representing the ID of the author who posted the chirp.

### Error Responses

If the query includes an invalid `author_id`, or if there is an internal error when retrieving chirps:

- **Invalid Author ID**:
    - **Code**: `400 BAD REQUEST`
    - **Content**: `{"error": "invalid author id"}`
- **Internal Server Error**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Couldn't retrieve chirps"}`

## Examples
### Get All Chirps

```bash
curl -X GET 'http://localhost:8080/api/chirps' -H 'Accept: application/json'
```

### Get Chirps by Author ID

```bash
curl -X GET 'http://localhost:8080/api/chirps?author_id=10' -H 'Accept: application/json'
```

### Get Chirps with Sorting
#### Sorting Ascending Order
```bash
curl -X GET 'http://localhost:8080/api/chirps?sort=asc' -H 'Accept: application/json'
```
#### Sorting Descending Order
```bash
curl -X GET 'http://localhost:8080/api/chirps?sort=desc' -H 'Accept: application/json'
```
### Get Chirps by Author ID and with Sorting
```bash
curl -X GET 'http://localhost:8080/api/chirps?author_id=10&sort=asc' -H 'Accept: application/json'
```

### Examples of valid URLs
- `GET http://localhost:8080/api/chirps?sort=asc&author_id=2`
- `GET http://localhost:8080/api/chirps?sort=asc`
- `GET http://localhost:8080/api/chirps?sort=desc`
- `GET http://localhost:8080/api/chirps`
## Notes

- The `author_id` parameter should be an integer. The API will respond with a `400 BAD REQUEST` if a non-integer value is provided.
- The `sort` parameter only excepts `asc` or `desc`, any other value will result in the sorting of ascending order.
- If no parameters are provided, the API will return all chirps in ascending order.