

compose-up: 
	docker-compose -f ./docker/docker-compose.dev.yaml up -d 

compose-down: 
	docker-compose -f ./docker/docker-compose.dev.yaml down


.PHONY: compose-up compose-down
