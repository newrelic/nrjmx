.PHONY : build
build:
	@($(MAVEN_BIN) clean package -DskipTests -P \!deb,\!rpm,\!tarball,\!test)

.PHONY : package
package:
	@($(MAVEN_BIN) versions:set -DnewVersion=\$(subst v,,$(TAG)))
	@($(MAVEN_BIN) clean -DskipTests package)

.PHONY : test
test:
	@($(MAVEN_BIN) clean test -P test)

.PHONY : sign
release/sign:
	@echo "=== [sign] signing packages"
	@bash $(CURDIR)/build/sign.sh
