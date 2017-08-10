cd `dirname $0`
protoc3 -I=. -I=$GOPATH/src --gogoslick_out=. protos.proto 
protoc3 -I=. -I=$GOPATH/src --gograin_out=. protos.proto 
go build
cd -
