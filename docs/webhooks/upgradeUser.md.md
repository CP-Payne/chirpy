# POST /api/polka/webhooks

This endpoint is designed to handle webhook requests from external services, specifically for upgrading a user's status based on events received from Polka.

## Request

`POST /api/polka/webhooks`

### Headers

- `X-API-Key`: The API key required for authentication to ensure that the request is authorized.

### Body Parameters

The request body must include an event type and the relevant data associated with that event.

- `event`: A string indicating the type of event. This endpoint specifically handles the `user.upgraded` event.
- `data`: An object containing the user's details. For this endpoint, it includes:
    - `user_id`: An integer representing the unique identifier of the user to be upgraded.

### Example Request Body

```json
{   
	"event": "user.upgraded",   
	"data": {    
		"user_id": 123   
	} 
}
```

## Response

### Success Response

- **Code**: `200 OK`
- **Content**: A message indicating that the user has been successfully upgraded.

### Example Success Response

`"User upgraded"`

### Error Responses

- **Unauthorized (API Key Error)**:
    - **Code**: `401 UNAUTHORIZED`
    - **Content**: `{"error": "error with apikey"}`
- **Internal Server Error (Decoding Parameters)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Couldn't decode parameters"}`
- **Not Found (User Not Found)**:
    - **Code**: `404 NOT FOUND`
    - **Content**: `{"error": "User not found"}`
- **Internal Server Error (General)**:
    - **Code**: `500 INTERNAL SERVER ERROR`
    - **Content**: `{"error": "Something went wrong"}`

## Examples

### Processing a User Upgrade Event

To process a user upgrade event from Polka, you would make the following request:

```bash
curl -X POST 'http://<your-domain>/api/polka/webhooks' \
-H 'X-API-Key: <polka-api-key>' \
-H 'Content-Type: application/json' \
-d '{"event": "user.upgraded", "data": {"user_id": 123}}'
```

Replace `<polka-api-key>` with the actual API key.

## Notes

- This endpoint requires a valid API key to be included in the `X-API-Key` header for authentication purposes. Requests without a valid API key or with an incorrect key will be rejected.
- The endpoint is designed to handle specific events from Polka. In this case, it processes events related to user upgrades. It verifies the event type before proceeding with the user upgrade logic.
- The response for non-matching event types is a successful HTTP status without performing any action, implying that the endpoint is designed to silently ignore unrelated events.