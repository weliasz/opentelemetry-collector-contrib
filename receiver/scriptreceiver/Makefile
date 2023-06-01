.PHONY: generate-checksums
generate-checksums:
	sh ./generate_checksums.sh
	gofmt -s -w ./procpipe/checksums.go

include ../../Makefile.Common

.PHONY: common
common: lint generate-checksums
