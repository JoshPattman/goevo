export GOEVO_VERSION=v0.4.0

nothing:
	@echo "Nothing to do"

upgrade-and-tidy-all: upgrade-all tidy-all

upgrade-all:
	@cd ./geno/floatarr && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./geno/neat && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./pop/neatpop && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./pop/simple && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./selec/tournament && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}
	@cd ./test && go get github.com/JoshPattman/goevo@${GOEVO_VERSION}

tidy-all:
	@cd ./ && go mod tidy
	@cd ./geno/floatarr && go mod tidy
	@cd ./geno/neat && go mod tidy
	@cd ./pop/neatpop && go mod tidy
	@cd ./pop/simple && go mod tidy
	@cd ./selec/tournament && go mod tidy
	@cd ./test && go mod tidy