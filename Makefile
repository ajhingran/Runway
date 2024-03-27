build:
	go build -o runway driver.go

run:
	./runway 09-15-2023 09-18-2023 -1 MSN DCA

composite:
	go run driver.go 04-11-2024 04-15-2024 -1 MSN-ORD DCA default default default default Frontier 700

clean:
	rm -f runway