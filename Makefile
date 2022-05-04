all: install build

# Install dependencies
install: install-frontend

install-frontend:
	yarn install --pure-lockfile

build: build-frontend build-backend
build-frontend:
	yarn build

build-backend:
	mage -v build:windows
