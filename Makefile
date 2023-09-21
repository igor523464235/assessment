ifndef STORAGES
STORAGES := 6
endif

export STORAGES

DOCKER_COMPOSE_FILE := docker-compose.yml

define BEGINNING_OF_FILE
version: "3.1"

networks:
  backend_network:
    driver: bridge

volumes:
  postgres_data:
    driver: local

services:
  db:
    image: postgres
    container_name: db
    ports:
      - 5432:5432
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: db
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - backend_network

  bff:
    container_name: bff
    image: bff
    command: 
      - "/app/wait-for-it.sh"
      - "db:5432"
endef

export BEGINNING_OF_FILE

define MIDDLE_OF_FILE
      - "--"
      - "/app/myapp"
    ports:
      - "8000:8000"
    depends_on:
      - db
endef

export MIDDLE_OF_FILE

define END_OF_FILE
    networks:
      - backend_network
    environment:
      - STORAGES=NUM

endef

export END_OF_FILE

define STORAGE_DESCR
  chunk_storage_NUM:
    container_name: chunk_storage_NUM
    volumes:
      - ./data/dataNUM:/data
    image: chunk_storage
    networks:
      - backend_network

endef

export STORAGE_DESCR

build-docker-compose:
	echo "$${BEGINNING_OF_FILE}" > $(DOCKER_COMPOSE_FILE); \

	@for i in $(shell seq 1 $(STORAGES)); do \
		echo "      - \"chunk_storage_NUM:8000\"" | sed "s/NUM/$$i/g" >> $(DOCKER_COMPOSE_FILE); \
	done

	echo "$${MIDDLE_OF_FILE}" >> $(DOCKER_COMPOSE_FILE); \

	@for i in $(shell seq 1 $(STORAGES)); do \
		echo "      - chunk_storage_NUM" | sed "s/NUM/$$i/g" >> $(DOCKER_COMPOSE_FILE); \
	done

	echo "$${END_OF_FILE}" | sed "s/NUM/$$STORAGES/g" >> $(DOCKER_COMPOSE_FILE); \

	@for i in $(shell seq 1 $(STORAGES)); do \
		mkdir ./data/data$$i -p; \
		echo "$${STORAGE_DESCR}" | sed "s/NUM/$$i/g" >> $(DOCKER_COMPOSE_FILE); \
	done

build-start: build-docker-compose
	docker build -t bff -f bff/Dockerfile .
	docker build -t chunk_storage -f chunk_storage/Dockerfile .
	docker-compose up --build

clean:
	sudo rm -R -f ./data
	sudo docker rm bff -f
	sudo docker rm db -f
	sudo docker volume rm str_postgres_data -f

start:
	docker-compose up