container_runtime := $(shell which podman || which docker)

$(info using ${container_runtime})

up: down
	${container_runtime} compose up --pull never --build -d

down:
	${container_runtime} compose down
