.PHONY : release/package
release/package: release/package-fips release/package-non-fips

.PHONY : release/package-fips
release/package-fips:
	@echo "=== [package-fips] Creating FIPS-compliant package"
	@(export MAVEN_OPTS="$(MAVEN_FIPS_OPTS)"; $(MAVEN_BIN) versions:set -DnewVersion=$(subst v,,$(TAG)) -f pom-fips.xml)
	@(export MAVEN_OPTS="$(MAVEN_FIPS_OPTS)"; $(MAVEN_BIN) clean package -DskipTests -f pom-fips.xml -P tarball-linux,deb,\!tarball-windows)
	@mkdir -p $(CURDIR)/dist
	@find target -name "*.jar" -o -name "*.tar.gz" -o -name "*.deb" | xargs -I {} cp {} $(CURDIR)/dist/

.PHONY : release/package-non-fips
release/package-non-fips:
	@echo "=== [package-non-fips] Creating non-FIPS package"
	@($(MAVEN_BIN) versions:set -DnewVersion=$(subst v,,$(TAG)) -f pom.xml)
	@($(MAVEN_BIN) clean package -DskipTests -f pom.xml -P tarball-linux,tarball-windows,deb,rpm)
	@mkdir -p $(CURDIR)/dist
	@find target -name "*.jar" -o -name "*.tar.gz" -o -name "*.rpm" -o -name "*.deb" -o -name "*.zip" | xargs -I {} cp {} $(CURDIR)/dist/

.PHONY : release/sign
release/sign:
	@echo "=== [sign] signing packages"
	@bash sign.sh

.PHONY : release/sign-rpm-fips
release/sign-rpm-fips:
	@echo "=== [sign-rpm-fips] signing FIPS RPM packages"
	@bash sign_rpm_fips.sh

.PHONY : release/publish
release/publish:
	@echo "=== [release/publish] publishing artifacts"
	@bash $(CURDIR)/build/upload_artifacts_gh.sh

.PHONY : release/publish-rpm-fips
release/publish-rpm-fips:
	@echo "=== [release/publish-rpm-fips] publishing FIPS RPM artifacts only"
	@bash $(CURDIR)/build/upload_fips_rpm_gh.sh

.PHONY : release
release: release/package release/sign release/publish
	@echo "=== [release] full pre-release cycle complete for nix"

.PHONY : release/rpm-fips
release/rpm-fips:
	@echo "=== [package-rpm-fips] Creating FIPS-compliant RPM package"
	@(export MAVEN_OPTS="$(MAVEN_FIPS_OPTS)"; $(MAVEN_BIN) versions:set -DnewVersion=$(subst v,,$(TAG)) -f pom-fips.xml)
	@(export MAVEN_OPTS="$(MAVEN_FIPS_OPTS)"; $(MAVEN_BIN) clean package -DskipTests -f pom-fips.xml -P tarball-linux,rpm)
	@mkdir -p $(CURDIR)/dist
	@find target -name "*.jar" -o -name "*.rpm" | xargs -I {} cp {} $(CURDIR)/dist/


.PHONY : release/package-fips/rpm
release/package-fips/rpm: release/rpm-fips release/sign-rpm-fips release/publish-rpm-fips
