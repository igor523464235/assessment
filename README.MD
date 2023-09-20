
## Architecture
Application consists of the bff service, that connects to the database and storage servers. The service handles client requests and implements the logic of splitting files into parts. Storage services just store that parts.
![Arch pic](/arch.jpg)

## How to run
To build&run applications, use:
```bash
make build-run
```
By default, 1 instance of the BFF application, one database, and 6 instances of the storage will be launched.

You can specify amount of storages:
```bash
make build-run STORAGES=3
```

## How to test
Testing module written as golang application that uploads picture test/picture.jpg, then downloads it and compares hashes of the original and downloaded file.
You need to run applications (see at the previous step), then: 

```bash
cd test/
go run test/main.go
```

### Functionality
- [x] Upload
- [x] Download
- [ ] Remove