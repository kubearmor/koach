module github.com/kubearmor/koach/koach

go 1.18

require go.uber.org/zap v1.21.0

require github.com/google/uuid v1.3.0 // indirect

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kubearmor/KubeArmor/protobuf v0.0.0-20220624172947-992a6abbeaf3
	github.com/kubearmor/koach/protobuf v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.12 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220621134657-43db42f103f7 // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gorm.io/driver/sqlite v1.3.4 // indirect
	gorm.io/gorm v1.23.6 // indirect
)

replace github.com/kubearmor/koach/protobuf => /home/nathaniel/Desktop/opensource/kubearmor/nathaniel-contrib/koach/protobuf
