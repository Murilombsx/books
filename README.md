# books

## Overview
Books Golang API

## Deployment & Execution Instructions

### Deployment
Create Docker Image with our application:
```
docker build -t books-api -f Dockerfile .
```

Run a Docker container using our recently built Image:
```
docker run -dp 8080:8080 books-api
```

### Execution
Example with `title` set to **math**:
```
curl -s -X GET "localhost:8080/books?title=math&limit=2"
{
   "items":[
      {
         "id":"zRpNEAAAQBAJ",
         "volumeInfo":{
            "title":"Math in the Time of Corona",
            "language":"en"
         }
      },
      {
         "id":"pdZSGb5DqzsC",
         "volumeInfo":{
            "title":"All the Math That's Fit to Print",
            "language":"en"
         }
      }
   ],
   "count":2
}
```

