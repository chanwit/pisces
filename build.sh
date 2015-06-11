export GOPATH=$PWD
export GOBIN=$GOPATH/bin

gofmt -s -w src/github.com/chanwit/pisces/*.go
gofmt -s -w src/cmd/*.go

go install github.com/chanwit/pisces
go install src/cmd/pisces-build.go
go install src/cmd/pisces-up.go
go install src/cmd/pisces-scale.go
go install src/cmd/pisces-stop.go
go install src/cmd/pisces-start.go
go install src/cmd/pisces-kill.go
