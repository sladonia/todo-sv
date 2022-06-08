gen:
	cd api; buf generate

test:
	go test -p 1 -count 1 ./...

.PHONY: gen test
