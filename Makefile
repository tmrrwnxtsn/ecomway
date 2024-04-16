.PHONY: compose-up
compose-up:
	docker-compose -f deployments/docker-compose.yml up

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: proto-engine
proto-engine:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/engine/*.proto

.PHONY: proto-integration
proto-integration:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/integration/*.proto

.PHONY: proto-shared
proto-shared:
	protoc --go_out=. --go_opt=paths=source_relative api/proto/shared/*.proto

.PHONY: proto
proto: proto-shared proto-engine proto-integration

.PHONY: swag-fmt
swag-fmt:
	@swag fmt -d internal/services/gateway/api/v1

.PHONY: swag-init
swag-init:
	@swag init -g handler.go -d internal/services/gateway/api/v1 -o api/swagger/gateway/v1 --outputTypes go,json

.DEFAULT_GOAL := compose-up
