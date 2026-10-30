TARGET   ?= http://target:8080
API_KEY  ?= sp33dt3st
STREAMS  ?= 10
IDLE_PHASE_DURATION ?= 5
DOWNLOAD_PHASE_DURATION ?= 10
UPLOAD_PHASE_DURATION ?= 10

run-target:
	docker compose up -d --build target

run-client:
	docker compose run --build --rm client \
		--idle-phase-duration=$(IDLE_PHASE_DURATION)s \
		--download-phase-duration=$(DOWNLOAD_PHASE_DURATION)s \
		--upload-phase-duration=$(UPLOAD_PHASE_DURATION)s \
		--streams=$(STREAMS) \
		--target=$(TARGET)

down:
	docker compose down