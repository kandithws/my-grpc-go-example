
# Usage make proto PACKAGE=xxx-service
proto:
	@echo "Making protobuf for ${$$PACKAGE}"
	@mkdir -p ./$$PACKAGE/src/genproto
	@protoc -I ./$$PACKAGE/proto $$PACKAGE/proto/*.proto --go_out=plugins=grpc:$$PACKAGE/src/genproto