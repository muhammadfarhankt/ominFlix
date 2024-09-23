build:
	docker-compose build --no-cache  # Rebuild images without cache
	
up:
	docker-compose up -d

down:
	docker-compose down -v  # Remove volumes

logs:
	docker-compose logs -f app

psql:
	docker exec -it omniflix-indexer-db-1 psql -U omniflix_user