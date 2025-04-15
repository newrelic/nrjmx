.PHONY : release/package
release/package:
	@($(MAVEN_BIN) versions:set -DnewVersion=\$(subst v,,$(TAG)))
	@($(MAVEN_BIN) clean package -DskipTests)

.PHONY : release/sign
release/sign:
	@echo "=== [sign] signing packages"
	@bash $(CURDIR)/build/sign.sh

.PHONY : release/publish
release/publish:
	@echo "=== [release/publish] publishing artifacts"
	@bash $(CURDIR)/build/upload_artifacts_gh.sh

.PHONY : release
release: release/package release/sign release/publish
	@echo "=== [release] full pre-release cycle complete for nix"