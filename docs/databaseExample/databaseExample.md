# Structure of `database.json`

The `database.json` file organizes data in a structured JSON format. Here is an example of how data is structured within the file:

```json
{   
"chirps": {     
		"1": {
		       "author_id": 1,       
		       "body": "I'm the one who knocks!",       
		       "id": 1,       
		       "created_at": "2024-04-01T21:55:08.111276848+02:00"     
		       },     
		"2": {     
		  "author_id": 1,       
		  "body": "Gale!",       
		  "id": 2,       
		  "created_at": "2024-04-01T21:55:08.118521674+02:00"     }    
		   // Additional chirps...   
		  },   

"users": {     
		"1": {
		"id": 1,       
		"email": "walt@breakingbad.com",       
		"Password": "...",       
		"is_chirpy_red": false     
		}, 
		    // Additional users...   },   
"tokens": {     // Token data...   
			} 
}
```
