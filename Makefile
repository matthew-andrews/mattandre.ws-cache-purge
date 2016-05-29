test:
	go test -v -race ./...

deploy:
	@apex deploy \
		-s CLOUDFLARE_API_KEY=$(CLOUDFLARE_API_KEY) \
		-s CLOUDFLARE_IDENTIFIER=$(CLOUDFLARE_IDENTIFIER) \
		-s CLOUDFLARE_AUTH_EMAIL="matt@mattandre.ws"
