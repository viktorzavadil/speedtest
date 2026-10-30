DURATION ?= 20
STREAMS  ?= 10
TARGET   ?= http://target:8080
API_KEY  ?= sp33dt3st

up:
	docker compose up -d --build target

run-test:
	docker compose run --build --rm client \
		--duration=$(DURATION)s \
		--streams=$(STREAMS) \
		--target=$(TARGET)

down:
	docker compose down