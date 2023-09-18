build:
	go build -o runway main.go

run:
	./main 09-15-2023 09-18-2023 -1 MSN DCA

composite:
	go run main.go 10-06-2023 10-09-2023 -1 DCA SFO default default default default Frontier