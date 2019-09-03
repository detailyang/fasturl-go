# Background color
GREEN  				:= $(shell tput -Txterm setaf 2)
YELLOW 				:= $(shell tput -Txterm setaf 3)
BLUE 				:= $(shell tput -Txterm setaf 4)
MAGENTA             := $(shell tput -Txterm setaf 5)
WHITE  				:= $(shell tput -Txterm setaf 7)
RESET  				:= $(shell tput -Txterm sgr0)
TARGET_MAX_CHAR_NUM := 20

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET} ${MAGENTA}[variable=value]${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)	

.PHONY: fuzz-url
## fuzz url decode and encode
fuzz-url:
	cd fuzz; go-fuzz-build && go-fuzz -func FuzzURL --workdir url 

.PHONY: fuzz-query
## fuzz query decode and encode
fuzz-query:
	cd fuzz; go-fuzz-build && go-fuzz -func FuzzQuery --workdir query 

.PHONY: test
## test everything
test:
	go test -race -v -run="^Test" github.com/detailyang/fasturl-go/fasturl

.PHONY: bench
## benchmark everything
bench:
	go test -v -benchmem -run="^$$" github.com/detailyang/fasturl-go/fasturl -bench Benchmark
