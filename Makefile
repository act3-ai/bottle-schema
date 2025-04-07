CONTROLLER_GEN_VERSION?=v0.14.0
CRD_REF_DOCS_VERSION?=v0.0.12
GOLANGCILINT_VERSION?=v1.57.2
GOMD2MAN_VERSION?=v2.0.4

all: build

.PHONY: generate
generate: tool/controller-gen
	go generate ./...

.PHONY: build
build: generate
	go build ./...

.PHONY: test
test: test-go
test-go:
	go test ./...

.PHONY: test-lint
test: test-lint
test-lint: tool/golangci-lint
	tool/golangci-lint run

.PHONY: cover
cover:
	go clean -testcache
	- rm coverage.txt
	go test ./... -coverprofile coverage.txt -coverpkg=$(shell go list)/...
	./filter-coverage.sh < coverage.txt > coverage.txt.filtered
	go tool cover -func coverage.txt.filtered
	# ./filter-coverage.sh < coverage.txt | sponge coverage.txt
	# go tool cover -func coverage.txt

.PHONY: clean
clean:
	# go clean -cache


tool/controller-gen: tool/.controller-gen.$(CONTROLLER_GEN_VERSION)
	GOBIN=$(PWD)/tool go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

tool/.controller-gen.$(CONTROLLER_GEN_VERSION):
	@rm -f tool/.controller-gen.*
	@mkdir -p tool
	touch $@


tool/crd-ref-docs: tool/.crd-ref-docs.$(CRD_REF_DOCS_VERSION)
	GOBIN=$(PWD)/tool go install github.com/elastic/crd-ref-docs@$(CRD_REF_DOCS_VERSION)

tool/.crd-ref-docs.$(CRD_REF_DOCS_VERSION):
	@rm -f tool/.crd-ref-docs.*
	@mkdir -p tool
	touch $@


tool/golangci-lint: tool/.golangci-lint.$(GOLANGCILINT_VERSION)
	GOBIN=$(PWD)/tool go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION)

tool/.golangci-lint.$(GOLANGCILINT_VERSION):
	@rm -f tool/.golangci-lint.*
	@mkdir -p tool
	touch $@


tool/go-md2man: tool/.go-md2man.$(GOMD2MAN_VERSION)
	GOBIN=$(PWD)/tool go install github.com/cpuguy83/go-md2man/v2@$(GOMD2MAN_VERSION)

tool/.go-md2man.$(GOMD2MAN_VERSION):
	@rm -f tool/.go-md2man.*
	@mkdir -p tool
	touch $@


.PHONY: tool
tool: tool/controller-gen tool/crd-ref-docs tool/golangci-lint tool/go-md2man

.PHONY: apidoc
apidoc: $(addsuffix .md, $(addprefix docs/apis/data.act3-ace.io/, v1alpha2 v1alpha3 v1alpha4 v1alpha5 v1beta1 v1))
docs/apis/%.md: tool/crd-ref-docs $(wildcard pkg/apis/$*/*_types.go) 
	@mkdir -p $(@D)
	tool/crd-ref-docs --config=apidocs.yaml --renderer=markdown --source-path=pkg/apis/$* --output-path=$@


# tool/conversion-gen:
# 	@mkdir -p tool
# 	GOBIN=$(PWD)/tool go install k8s.io/code-generator/cmd/conversion-gen@v0.24.0

# tool/defaulter-gen:
# 	@mkdir -p tool
# 	GOBIN=$(PWD)/tool go install k8s.io/code-generator/cmd/defaulter-gen@v0.24.0

# .PHONY: gen_conversion
# gen_conversion: tool/controller-gen 
# 	tool/conversion-gen -i github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha2,github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha3,github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha4,github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha5,github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1beta1 --go-header-file ./header.txt -p pkg/api -v=2 -O zz_generated.conversion

# .PHONY: gen_defaults
# gen_defaults: tool/controller-gen 
# 	tool/defaulter-gen -i github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha2,github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha3,github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha4,github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1alpha5,github.com/act3-ai/bottle-schema/pkg/apis/data.act3-ace.io/v1beta1 --go-header-file ./header.txt -v=6

# The above command output files to the wrong path because they are only meant to be run in the GOPATH.  This can be hacked around by using the script https://github.com/kubernetes/kubernetes/blob/master/hack/run-in-gopath.sh
