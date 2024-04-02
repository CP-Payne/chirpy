## Serving Static Files with `/app` Endpoint

The `/app` endpoint is defined to serve static files, including an HTML file that displays a welcome message:

```html
<html> 
	<body>     
		<h1>Welcome to Chirpy</h1> 
	</body> 
</html>
```

### Example Usage

```bash
curl http://localhost:8080/app/
```

This will respond with the HTML content defined above.
