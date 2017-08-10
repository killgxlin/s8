cd `dirname $0`
protoc3 -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gogoslick_out=. protos.proto 
protoc3 -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gograin_out=. protos.proto 
rm -r github_com google
go build
cd -
