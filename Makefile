BINARY_NAME=mailapp
DSN="host=localhost port=5432 user=postgres password=password dbname=database sslmode=disable"


start:
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BINARY_NAME} ./main
	@echo "Built!"
	@env DSN=${DSN}  ./${BINARY_NAME} &
	@echo "Started!"


clean:
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"


stop:
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"


restart: stop start
