build:
	docker build -t tent .

%.json:
	docker run --rm -i tent < $@

build-dev:
	docker build -t tent-dev -f dev.Dockerfile .

dev:
	docker run --rm -it -v $(shell pwd):/tents tent-dev CompileDaemon -command="$(command)"
