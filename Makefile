runserver:
	@docker build -t server -f nginx/Dockerfile nginx
	@docker run -it --rm -p 8080:80 server

slowattack:
	@go run main.go