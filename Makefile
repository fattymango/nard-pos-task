#################### GRPC ####################
.PHONY: init-proto
init-proto:
	@go get -u google.golang.org/protobuf@v1.26.0 
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest	
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	


# Generate the proto files
.PHONY: proto
proto:
	@protoc --go_out=./proto/ --go-grpc_out=./proto/ proto/*.proto



#################### Testing ####################
.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: bench
bench:
	@for d in $$(go list ./...); do \
		# go test -bench=.  -benchmem -cpuprofile=$(HOME)/vfx12/vfxmarket/pprof/cpu/cpu_$$(basename $$d).pprof -memprofile=$(HOME)/vfx12/vfxmarket/pprof/mem/mem_$$(basename $$d).pprof $$d; \
		go test -bench=.  -benchmem $$d; \
	done
	@$(MAKE) clean-test
	

.PHONY: clean-test
clean-test:
	@rm ./*.test


#################### Linting ####################
.PHONY: format
fmt , format:
	@go fmt ./...



#################### RUN ####################
	
# Run the multi_tenant multi_tenant
.PHONY: run
run:
	@ $(MAKE) build-quick
	@ $(shell source exports.sh)
	@ ./bin/multi_tenant



# Race Detector
.PHONY: race
race:
	@CGO_ENABLED=1 go run -race cmd/multi_tenant/*.go


.PHONY: migrate
migrate:
	@$(MAKE) build-migrate-quick
	@./bin/migrate

#################### Profiling ####################
.PHONY: pprof-cpu
pprof-cpu:
	@go tool pprof -http=":8000" pprofbin ./cpu.pprof

.PHONY: pprof-mem
pprof-mem:
	@go tool pprof -http=":8000" pprofbin ./mem.pprof 


#################### Build Executable ####################
# Build amd64	for alpine
.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w -s' -o ./bin/multi_tenant cmd/multi_tenant/*.go

# Build depending on the OS
.PHONY: build-quick
build-quick:
	@go build  -o ./bin/multi_tenant cmd/multi_tenant/*.go


.PHONY: build-migrate
build-migrate:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w -s' -o ./bin/migrate cmd/migrate/*.go

.PHONY: build-migrate-quick
build-migrate-quick:
	@go build  -o ./bin/migrate cmd/migrate/*.go




#################### Docker Compose ####################

# Stack
.PHONY: up
up:
	@docker compose up -d --build --force-recreate --remove-orphans

.PHONY: down
down:
	@docker compose down 

.PHONY: top
top:
	@docker stats



