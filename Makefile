.PHONY: all build binaries

all: run

run:
	echo "Doing nothing"

build:
	glide up
	# remove vendor from terraform (FIX for https://github.com/coreos/etcd/issues/9357)
	@rm -rf vendor/github.com/hashicorp/terraform/vendor
	go build