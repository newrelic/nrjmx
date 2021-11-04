all: clean thrift-gen build-java

.PHONY: clean
clean:
	@printf "===  Cleanup...\n"
	@rm -rf ./generated

.PHONY: thrift-gen
thrift-gen:
	@printf "===  Generating thrift code...\n"
	@mkdir ./generated
	@thrift -out ./generated -r --gen go:package_prefix=generated/jmx ./thrift/jmx.thrift
	@thrift -out ./nrjmx/src/main/java/ -r --gen java ./thrift/jmx.thrift

.PHONY: build-java
build-java:
	@printf "===  Build nrjmx...\n"
	cd ./nrjmx; mvn clean package -P \!deb,\!rpm,\!tarball,\!test -DskipTests

.PHONY: install
install:
	sudo cp -r ./nrjmx/bin/nrjmx /usr/local/bin/nrjmx
	sudo cp -r ./nrjmx/bin/nrjmx.jar /usr/local/bin/nrjmx.jar

.PHONY: test
test:
	@printf "===  Running tests...\n"
	cd tests && docker-compose up -d && sleep 30s
	@printf "===  Waiting for service...\n"
	go clean -testcache
	go test -v .; ret=$$?; \
	cd tests && docker-compose down; \
	exit $$ret