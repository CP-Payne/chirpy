# GET /admin/metrics

This endpoint provides the admin with the current metrics of the Chirpy application, particularly the number of times Chirpy has been accessed.

## Request

`GET /admin/metrics`

## Response

### Success Response

- **Code**: `200 OK`
- **Content**: An HTML page displaying a welcome message to the Chirpy admin and the number of hits to the application.

### Example Response

The response will be an HTML document similar to the following:

```html
<html>   
	<body>     
		<h1>Welcome, Chirpy Admin</h1>     
		<p>Chirpy has been visited [number] times!</p>   
	</body> 
</html>
```

`[number]` will be replaced with the number of hits on the home page (`/app`)

### Error Responses

This endpoint is straightforward and  does not have error responses.

## Examples

### Retrieve Metrics

To access the metrics page, an admin would make the following request:

```bash
curl -X GET 'http://localhost:8080/admin/metrics'
```

## Notes

- As of now, the provided example does not include authentication logic.
- The metrics returned in this example are simple hit counts.
