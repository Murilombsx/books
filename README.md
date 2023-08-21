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

## Kubernetes Instructions

First we need to tag the Docker Image created before:
```
docker tag books-api murilobarceloss/books-api
```

After that we need to push the Docker Image to the Docker Repositories:
```
docker push murilobarceloss/books-api
```

With the Image ready to be used we can apply the `deployment.yml` file:

```
kubectl apply -f deployment.yml
```

Once the **Deployment** and **Pods** are ready, we can register the Service and get the address(I used `minikube` locally):
```
minikube service books-api-service
```

With all up and running we can execute the same tests using the IP provided by `minikube`.
