all: init generate

init:
	@echo "Making tonic file server..."
#add special compiler options beyond -w later,
#maybe use a packer or some sort.
generate:
	CGO_ENABLED=0 go build -ldflags "-w" .

clean:
	@echo "done."
