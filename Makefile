all:
	go build -buildvcs=false

test:
	go test ./vdb
	go test ./ollm

clean:
	rm ./git-llama
