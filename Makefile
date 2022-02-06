dev:
	docker-compose up
build:
	docker build -t app-prod . --target production
start:
	docker run -e DB_HOST='localhost' -p 8000:8000 --name app-prod app-prod
db:
	docker-compose up -d db