.PHONY: all clean build-linux docker-build

BINDIR := bin
BINDIR_LINUX := bin-linux
PROGRAMS := $(BINDIR)/describe $(BINDIR)/histogram $(BINDIR)/logregpredict $(BINDIR)/logregtrain $(BINDIR)/pairplot $(BINDIR)/scatterplot

INTERNAL_SOURCES := internal/hogwarts/dataset.go \
                    internal/logisticregression/model.go \
                    internal/stats/stats.go

all: $(PROGRAMS)

$(BINDIR):
	mkdir -p $(BINDIR)

$(BINDIR)/describe: cmd/describe/describe.go $(INTERNAL_SOURCES) | $(BINDIR)
	go build -o $@ ./cmd/describe

$(BINDIR)/histogram: cmd/histogram/histogram.go $(INTERNAL_SOURCES) | $(BINDIR)
	go build -o $@ ./cmd/histogram

$(BINDIR)/logregpredict: cmd/logregpredict/logreg_predict.go $(INTERNAL_SOURCES) | $(BINDIR)
	go build -o $@ ./cmd/logregpredict

$(BINDIR)/logregtrain: cmd/logregtrain/logreg_train.go $(INTERNAL_SOURCES) | $(BINDIR)
	go build -o $@ ./cmd/logregtrain

$(BINDIR)/pairplot: cmd/pairplot/pair_plot.go $(INTERNAL_SOURCES) | $(BINDIR)
	go build -o $@ ./cmd/pairplot

$(BINDIR)/scatterplot: cmd/scatterplot/scatter_plot.go $(INTERNAL_SOURCES) | $(BINDIR)
	go build -o $@ ./cmd/scatterplot

clean:
	rm -rf $(BINDIR)
	rm -rf $(BINDIR_LINUX)

build-linux:
	@echo "Creating output directory for binaries..."
	@mkdir -p $(BINDIR_LINUX)
	@echo "Building Docker image and compiling binaries..."
	@docker build -t dslx-builder .
	@echo "Extracting binaries to $(BINDIR_LINUX)/..."
	@docker run --rm -v "$$(pwd)/$(BINDIR_LINUX):/binaries" dslx-builder
	@echo ""
	@echo "Done! Linux binaries are available in $(BINDIR_LINUX)/"
	@ls -lh $(BINDIR_LINUX)/

docker-build: build-linux

