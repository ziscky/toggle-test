# Toggle-Test

## Building
```
git clone https://github.com/ziscky/toggle-test.git
go build
```

## Running
#### Default options
| HTTP Address |DB  |
|--|--|
| :8080 | :memory: |

```
./toggle-test
```
#### Custom Options
`ADDR=:8080 DB=/path/to/db ./toggle-test`

## API 
### Swagger file
 `api.swagger.yaml`
### Routes
- Create shuffled deck with specific cards
`POST /deck?cards=AS,2S,3D&shuffled=true`
- Create default unshuffled deck
`POST /deck`
- Open deck
`GET /deck?deck_id=a251071b-662f-44b6-ba11-e24863039c59`
- Draw 3 cards from deck
`POST /deck/draw?count=3&deck_id=a251071b-662f-44b6-ba11-e24863039c59`
- Draw 1 card from deck
`POST /deck/draw?deck_id=a251071b-662f-44b6-ba11-e24863039c59`

##  Tests
Test coverage is provided in the file `cover.out`
To run all tests and view coverage:
```
go test -coverprofile cover.out ./...
go tool cover -html=cover.out
```

## Documentation

### Folder structure
 - **api/**
	Contains the API structs and HTTP handlers
- **internal/**
	- **games/**
	Contains functions to generate playing cards and shuffle them. Other game specific functions can be added here.
	- **models/**
	Contains the database Models
	- **persist/**
	Contains the interface required to implement database operations
	- **sql/**
	Contains the implementation of the persist interface and database migrations
- **test/**
    Contains test mocks and test data helper functions
    
 ### Dependencies
 - github.com/glebarez/go-sqlite
 Chosen because it is a PureGo implementation that does not need CGO
 - github.com/sirupsen/logrus
 Structured logging library
 - github.com/stretchr/testify
 Generation of mocks and test helpers
 - github.com/pressly/goose/v3
 To perform programmatic database migrations
 - gorm.io/gorm
 The ORM 
