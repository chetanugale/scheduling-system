## Scheduling System ##
Scheduling system to manage events and recommend ideal slots for maximum users to attend. 

```
scheduling-system
+--cmd
|   +--main.go
+--constants
|   +-- constants.go
+--Dockerfile
|   +-- Dockerfile
+--docs
+--handlers
|   +-- handlers.go
|   +-- handlers_test.go
+--mocker
|   +-- mock.go
+--models
|   +-- models.go
+--repository
|   +-- data.go
+--services
|   +-- services.go
```


### Event management
- **Create, Update, Delete** events

    Create Event:
    	`POST "/events"`    
        Payload : 
    ```
            {
                "title":"test event",
                "estimatedMins":30,
                "slots":[
                {   
                    "startTime":"2025-05-06T16:00:00Z",
                    "endTime":"2025-05-06T16:30:00Z"
                },
                {   
                    "startTime":"2025-05-06T14:00:00Z",
                    "endTime":"2025-05-06T14:30:00Z"
                }
                ]

            }
    ```  
	Get All Event:
        `GET "/events"`

    Get Event by ID:
	    `GET"/events/:id"` 

    Update Event:
	    `PUT "/events/:id"` 
    ```
            {
                "title":"test event",
                "estimatedMins":30,
                "slots":[
                {   
                    "startTime":"2025-05-06T16:00:00Z",
                    "endTime":"2025-05-06T16:30:00Z"
                },
                {   
                    "startTime":"2025-05-06T14:00:00Z",
                    "endTime":"2025-05-06T14:30:00Z"
                }
                ]

            }
    ``` 
    Delete Event:
	    `DELETE "/events/:id"`

### Availability management
- **Create, Update, Delete** availability of users

    Post Availability:
        `POST "/availability"`
    ```
            {
                "eventId":"681824b939e50f0b5f59eb7b",
                "slotId":"681824b939e50f0b5f59eb7a",
                "userId":"user1"
            }
    ```
    Get Availability based on EventID:
        `GET "/event/:id/availability"`

    Update Availability:    
        `PUT "/availability/:id"`
    ```
            {
                "eventId":"681824b939e50f0b5f59eb7b",
                "slotId":"681824b939e50f0b5f59eb7a",
                "userId":"user1"
            }
    ```
    Delete Availability:
        `DELETE "/availability/:id"`
### Recommendation
- Provide recommendation for probable event scheduling based on maximum user availability

    Get recommendations:
        `GET "/events/:id/recommend"`


### To run the system:
- Please use Make command : `make run`

### To run Unit-tests:
- Please use Make command : `make test`

 
Please follow Makefile to run other required commands.

Please follow directory structure : `/github.com/chetanugale/scheduling-system` to avoid wiring mismatch.