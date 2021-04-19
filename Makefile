SHELL:=/bin/bash
+DESTDIR:=.
GOARCH:=$(shell go env GOARCH)
GOOS:=$(shell go env GOOS)
EXEC:=./bible
INPUTDIR=$(shell pwd)/input/one2one
OUTPUTDIR=$(shell pwd)/output/one2one

.PHONY: build
build:  ## build go binary
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOENVS) go build  -o bible ./app/main.go

# .PHONE: generate
# generate: build ## generate bible text
# 	@for f in $(shell ls ${INPUTDIR}); \
# 		do \
# 			$(EXEC) -config $(INPUTDIR)/$${f} > ${OUTPUTDIR}/$${f}.md; \
# 			$(EXEC) -config $(INPUTDIR)/$${f} -lang kor > ${OUTPUTDIR}/$${f}_kor.md; \
# 	done;
	

	# $(EXEC)
########
# Help #
########

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
