APP_NAME=go-monorepo
SERVICE_AUTH_IMAGE=somprasongd/go-monorepo/auth
SERVICE_AUTH_VERSION=0.0.1
SERVICE_TODO_IMAGE=somprasongd/go-monorepo/todo
SERVICE_TODO_VERSION=0.0.1

export APP_NAME
export SERVICE_AUTH_IMAGE
export SERVICE_AUTH_VERSION
export SERVICE_TODO_IMAGE
export SERVICE_TODO_VERSION

dev-up:
	@echo "---Start Dev $(APP_NAME) Environtment---"
	docker-compose -p monorepo-dev -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.dev.yml up -d

dev-down:
	@echo "---Stop Dev $(APP_NAME) Environtment---"
	docker-compose -p monorepo-dev -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.dev.yml down

dev-down-v:
	@echo "---Stop Dev $(APP_NAME) Environtment---"
	docker-compose -p monorepo-dev -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.dev.yml down -v

build-auth:
	@echo "---Build $(APP_NAME) $(SERVICE_AUTH_IMAGE):$(SERVICE_AUTH_VERSION)---"
	docker build -t $(SERVICE_AUTH_IMAGE):$(SERVICE_AUTH_VERSION) -f ../apps/auth-api/Dockerfile .

# build with debug
build-auth-d:
	@echo "---Build $(APP_NAME) $(SERVICE_AUTH_IMAGE):$(SERVICE_AUTH_VERSION)---"
	docker build --progress plain -t $(SERVICE_AUTH_IMAGE):$(SERVICE_AUTH_VERSION) -f ../apps/auth-api/Dockerfile .

build-todo:
	@echo "---Build $(APP_NAME) $(SERVICE_TODO_IMAGE):$(SERVICE_TODO_VERSION)---"
	docker build -t $(SERVICE_TODO_IMAGE):$(SERVICE_TODO_VERSION) -f ../apps/todo-api/Dockerfile .

# build with debug
build-todo-d:
	@echo "---Build $(APP_NAME) $(SERVICE_TODO_IMAGE):$(SERVICE_TODO_VERSION)---"
	docker build --progress plain -t $(SERVICE_TODO_IMAGE):$(SERVICE_TODO_VERSION) -f ../apps/todo-api/Dockerfile .

prod-up:
	@echo "---Start Prod $(APP_NAME)---"
	docker-compose -p task-api-prod -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.prod.yml up -d

prod-down:
	@echo "---Stop Prod $(APP_NAME)---"
	docker-compose -p task-api-prod -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.prod.yml down

prod-down-v:
	@echo "---Stop Prod $(APP_NAME)---"
	docker-compose -p task-api-prod -f ./deploy/docker-compose.yml -f ./deploy/docker-compose.prod.yml down -v
