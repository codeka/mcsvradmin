

all:
	@$(MAKE) -C static
	go build -o mcsvradmin

