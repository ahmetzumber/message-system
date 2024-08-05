run:
	go build -o message-system . && APP_ENV="local" ./message-system

create-mocks:
	mockgen -source=app/service/service.go -destination=mocks/mock_clients.go -package=mocks
	mockgen -source=app/handler/handler.go -destination=mocks/mock_service.go -package=mocks


unit-test:
	go test -v ./... -tags=unit