# Uploader
Uploader is a microservice that processes and converts Excel (.xlsx) files for company budgets to a database.  
Various types of data can then be queried via the UI and are output as JSON objects.

## About 
When running, users can go to `localhost:8080/` to access the menu. 
From there, users can upload an Excel file `.xlsx` or a `.zip` file with correct format.
Uploader will then convert and merge the files with information from a CSV file it reads on startup containing budget info.
Users can also query data by `vessel name`, `vessel+year`, and `vessel+year+budget code`.


## Setup
* Install [Go](https://golang.org/)
* Install [MongoDB](https://www.mongodb.com/download-center?jmp=homepage#community)
* Add $GOPATH ("C:\Users\\$USER\go") and MongoDB to PATH
* `go get github.com/ysaliens/uploader` to get project files and all dependencies. This can take a while for a new install.
* `mkdir $GOPATH/src/github.com/ysaliens/uploader/db` <--- This is where database files will be stored
* `cd $GOPATH/src/github.com/ysaliens/uploader/`

## Build
`go build` from root of project.

If missing dependencies (or used git clone), run `go get -v ./...` to get all dependencies

## Run 
* `mongod --dbpath "$GOPATH/src/github.com/ysaliens/uploader/db` to start database (or use provided `docker-compose.yml`)
* `./uploader.exe` to start server (separate window)
* Navigate to `http://localhost:8080`  

![UI](/files/config/UI.PNG)

## Architecture
Uploader uses MongoDB for a database. It is written in Golang utilizing the high-performance Gin-Gonic framework.
On startup, the server will attempt to read `./files/config/budget_codes.csv` which contains opex, category, budget code, and budget description information.
This information is then stored as a hash map that is passed to the various handlers. 
This is done for speed as information read from the Excel files needs to be spliced with the information from the budget codes CSV file.  

When a user uploads an Excel file, information is converted to the new format. 
The information is then checked against existing records in the database prior to being saved to avoid duplicates.    
All new records are then saved to the database.


All files with the correct extension are first saved to `./files/cache`. Zip files are extracted into their own folder and processed.
Upon process completion, temp files are deleted.

## Performance
The code is written to concurrently insert data into the database via go routines.
It doesn't currently do that as performance and concurrency should be balanced against server load. 
Even without concurrency, service proceses an Excel file every 2-5 seconds.
The database is a bottleneck during full-speed inserts, even if the same connection is kept for a sheet/file.
![UI](/files/config/Output.PNG)

## Docker
Uploader is ready to be containerized. Provided `Dockerfile` will build the service.
Use `docker-compose-yml` to get a MongoDB instance.


## TO-DOs
Uploader was written in <3 days - TO-DOs are below:
* __Automated unit testing__ Add more cases. Add proper mock db and tests for model layer
* __Split long functions__ into smaller, modular functions - mostly in uploader.go
* __Optimize Docker images__ Size can be optimized, combine into a single service with service dependent on db.
* __All the @TO-DOs__ scattered across the code.
* __UI Improvements__ The UI could use a lot of love. Since Golang cannot run client-side code, adding some Javascript and a new design would really help.

