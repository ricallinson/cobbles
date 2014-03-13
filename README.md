cobbles
=======

Multidimensional Configuration for Golang

## Testing

The following should all be executed from the `cobbles` directory _$GOPATH/src/github.com/ricallinson/cobbles/_.

### Install

    go get github.com/ricallinson/cobbles

### Run

    go test

### Generate Code Coverage

    go test -cover

To view the coverage run;
    
    go test -coverprofile=coverage.out
    go tool cover -html=coverage.out