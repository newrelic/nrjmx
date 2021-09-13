.PHONY : build
build:
	@($(MAVEN_BIN) clean package -DskipTests -P \!deb,\!rpm,\!tarball,\!test)

.PHONY : test
test:
	@($(MAVEN_BIN) clean test -P test)

