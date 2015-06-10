export GOPATH=$PWD
export GOBIN=$GOPATH/bin
go install github.com/chanwit/pisces
go install src/cmd/pisces-build.go
go install src/cmd/pisces-up.go
go install src/cmd/pisces-scale.go
go install src/cmd/pisces-stop.go
go install src/cmd/pisces-start.go
