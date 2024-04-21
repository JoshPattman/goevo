export GOEVO_VERSION=v0.4.1

nothing:
	@echo "Nothing to do"

release-new-version-for-parent-pkg:
	@git tag -a ${GOEVO_VERSION} -m "Release ${GOEVO_VERSION}"
	@git push --tags

release-new-version-for-sub-pkgs:
	@git tag -a geno/floatarr/${GOEVO_VERSION} -m "Release ${GOEVO_VERSION} for geno/floatarr"
	@git tag -a geno/neat/${GOEVO_VERSION} -m "Release ${GOEVO_VERSION} for geno/neat"
	@git tag -a pop/neatpop/${GOEVO_VERSION} -m "Release ${GOEVO_VERSION} for pop/neatpop"
	@git tag -a pop/simple/${GOEVO_VERSION} -m "Release ${GOEVO_VERSION} for pop/simple"
	@git tag -a pop/hillclimber/${GOEVO_VERSION} -m "Release ${GOEVO_VERSION} for pop/hillclimber"
	@git tag -a selec/tournament/${GOEVO_VERSION} -m "Release ${GOEVO_VERSION} for selec/tournament"
	@git tag -a selec/elite/${GOEVO_VERSION} -m "Release ${GOEVO_VERSION} for selec/elite"

	@git push --tags

upgrade-and-tidy-all: upgrade-all tidy-all

upgrade-all:
	@cd ./geno/floatarr && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./geno/neat && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./pop/neatpop && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./pop/simple && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./pop/hillclimber && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./selec/tournament && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./selec/elite && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./test && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}

tidy-all:
	@cd ./ && go mod tidy
	@cd ./geno/floatarr && go mod tidy
	@cd ./geno/neat && go mod tidy
	@cd ./pop/neatpop && go mod tidy
	@cd ./pop/simple && go mod tidy
	@cd ./pop/hillclimber && go mod tidy
	@cd ./selec/tournament && go mod tidy
	@cd ./selec/elite && go mod tidy
	@cd ./test && go mod tidy