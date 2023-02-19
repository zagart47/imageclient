# Image Client.

The client is written to work with the [server](https://github.com/zagart47/imageserver) as part of a test job for the company.

## Installation

Use the git cli to clone repository.

```bash
git clone https://github.com/zagart47/imageclient.git
```

## Usage

### Client
```bash
cd imageclient
```

Server and client supports 3 actions:
```upload```, ```list``` and ```download```.

How to upload files:
```go
go run cmd/app/main.go -ul {FILEPATH}
```
Where {FILEPATH} is the path to the file relative to the client's root folder.

File upload example:
```go
go run cmd/app/main.go -ul files/test.png
```
The "files" folder is located in the root folder of the client.

How to show uploaded files:
```go
go run cmd/app/main.go -ls
```


How to download uploaded file:
```go
go run cmd/app/main.go -dl {FILENAME}
```
Where {FILENAME} is the filename you chose from the table after ```go run cmd/app/main.go -ls```.

File download example:
```go
go run cmd/app/main.go -dl test.png
```
