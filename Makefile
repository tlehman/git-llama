all:
	go build

test:
	go test ./vdb
	go test ./ollm

clean:
	rm ./git-llama
