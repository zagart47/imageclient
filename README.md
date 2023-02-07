# Test Task solution for "Tages".

Test task from Tages company.

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
go run api/client.go -o upload -f {FILEPATH}
```
Where {FILEPATH} is the path to the file relative to the client's root folder.

File upload example:
```go
go run api/client.go -o upload -f files/test.png
```
The "files" folder is located in the root folder of the client.

How to show uploaded files:
```go
go run api/client.go -o list
```


How to download uploaded file:
```go
go run api/client.go -o download -f {FILENAME}
```
Where {FILENAME} is the filename you chose from the table after ```go run client.go -o list```.

File download example:
```go
go run api/client.go -o download -f test.png
```
