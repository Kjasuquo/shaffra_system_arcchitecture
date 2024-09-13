# Part 1: GoLang Problem-Solving Task

## About
Architecturally, This was set up in 3 layers -> Repository, Service and Controller layer. With the User model/domain travelling through the 3 layers

## Technology Used
- Router: Gin-gonic
- Database: Postgres (GORM used to migrate the tables and queries)

## Installation
Before installing, download and install go version 1.22.1 and above from **[golang's official website](https://go.dev/doc/install)**
Ensure you have docker if you are running it on docker.
You can also set up postgres DB locally

## Install make
for MacOS run
```brew install make```

### Make available the correct set of dependencies
```go mod tidy```

## Test
- The database is mocked using `github.com/golang/mock`.
- Tests done for all endpoints, which calls mocked db methods
- Test also done for the function in `internal/pkg/web/response.go` used in returning response
- To run tests with coverage run `make test`

### Starting server
- create a `.env` file in the base directory with sample fields from `.env.sample`
- run `make postgres` to provision a postgres DB on docker
- run `make run` to run the service with the provisioned DB locally.

You can also start up the app with the docker build using `make docker-up`. 
This will build the app in the dockerfile present in the base directory.
- NOTE: You will need to change the DB_HOST in your env to `postgres`.
- After starting the service on `port: 5053`, use `http://localhost:5053/ping` to test the service


### Sample Payload for POST and PUT user
```
{
    "name": "Joseph A",
    "email": "okoasuquo@gmail.com",
    "age": 24
}
```


_______________________________________________________________________________________________________________________________________________________________________________________________


# Part 2: Debugging and Problem-Solving (Code Review Task)

## Initial Step taken
- Started the service with `go run buggy_project.go` - It started but database connection error was not handled
- Tried using `http://localhost:8080/create?name=Joseph` to create a user, but it threw an error that connection was refused, and it returned status code 200 despite the error

## Improvements
1. Provided database in the makefile connected it appropriately. All you have to run `make buggy-db` to spin it up before starting the service with `go run buggy_project.go`
2. Handled error while connecting to database
3. Create the 'users' table if it doesn't exist. This ideally should be in the makefile while provisioning db
4. Handled error and return with the write status code in the getUsers endpoint
5. In the createUser, handled error if username is not present in query
6. Also re-wrote the insert db query to guard against SQL injection attacks
7. And, also handled the error of that db query with the right status code



_______________________________________________________________________________________________________________________________________________________________________________________________


# Part 3: Systems Design and Architecture

## Service Design
Each of the major services will be written in 3 layers: Controller, Service and Repository for easy test and scaling of each codebase
- Controller Layer: For Routing requests
- Service Layer: For the business logic
- Repository Layer: For database integration and methods

The 5 microservices I would consider are:
1. Queuing and Logging Service: I prefer to use a message broker to communicate between the different microservices instead of http request for scalability and fault tolerance.
- Used to communicate between all the microservices through RabbitMQ
- Log and measure all calls on other microservices using prometheus and grafana

2. User microservice: (For Admin and Buyer (users))
- Handles user's account creation and login
- Authentication and Authorization middleware
- All user related duties such as user profile (set up and update)
Database: PostgresDB because it is well suited for structured data and relatively fast

3. Product Service: The real marketplace
- Handles all product operations such as create, update and delete products by admin User.
- View available product by users
- Keeping track of products
Database: MongoDB because it is well suited for large unstructured data

4. Order Service:
- Used by the user to Creates and updates orders
- Keeps track of orders/transaction history
- Handles payment integration
Database: PostgresDB because it is well suited for structured transactional data and relatively fast

5. Notification Service:
- Processes all notifications including sending successful/unsuccessful mails
- Order Receipts
- Can be scaled later for other mode of notifications

## DATABASE MANAGEMENT
This will be backed up regularly to avoid data loss

## SCALING
I will deploy the microservices in GKE Kubernetes Cluster with autoscaling which will help it scale horizontally depending on the load. The Load balancer will distribute incoming network traffic across various instances. 
This is to ensure that no single instance is overwhelmed with too many requests and helps in managing the load effectively.

Consul can be integrated to monitor service.


## CI/CD PIPELINE FOR EACH MICROSERVICES
1. GitHub will be used as version control system to store source code: Will set up 2 major environments `develoment and master`
2. Development environment is for testing
3. GitHub workflow set up to run tests on push/merge to development and master
4. Write a Dockerfile to build a docker image for each service and push to GCP container registry
5. Write a multistage Dockerfile to minimize image size 
6. Define kubernetes deployment manifest for each service and deploy in such a way that it deploys the latest build whenever it is pushed