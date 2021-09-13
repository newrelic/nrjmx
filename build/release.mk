.PHONY : package
package:
	@($(MAVEN_BIN) versions:set -DnewVersion=\$(subst v,,$(TAG)))
	@($(MAVEN_BIN) clean package -DskipTests)

.PHONY : sign
release/sign:
	@echo "=== [sign] signing packages"
	@bash $(CURDIR)/build/sign.sh

publish:
	@echo "=== [release/publish] publishing artifacts"
	@bash $(CURDIR)/build/upload_artifacts_gh.sh

release: package sign publish
	@echo "=== [release] full pre-release cycle complete for nix"