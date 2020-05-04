
# Usage make proto SERVICE=xxx-service
proto:
	@echo "Making protobuf for ${$$SERVICE}"
	@mkdir -p ./$$SERVICE/src/genproto
	@protoc -I ./$$SERVICE/proto $$SERVICE/proto/*.proto --go_out=plugins=grpc:$$SERVICE/src/genproto

gateway_proto:
	@echo "Making protobuf for api-gateway service=${$$SERVICE}"
	@mkdir -p ./api-gateway/src/genproto/$$SERVICE
	@protoc -I ./$$SERVICE/proto $$SERVICE/proto/*.proto --go_out=plugins=grpc:api-gateway/src/genproto/$$SERVICE