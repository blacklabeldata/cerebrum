.PHONY: test

# CWD=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

default: test

test:
	@echo "------------------"
	@echo " test"
	@echo "------------------"
	@godep go test -v -coverprofile=coverage.out -covermode=set
	# @godep go test -v -coverprofile=$(CWD)/coverage.out -covermode=set

bench:
	@echo "------------------"
	@echo " benchmark"
	@echo "------------------"
	@godep go test -test.bench "^Bench*" -benchmem

html:
	@echo "------------------"
	@echo " html report"
	@echo "------------------"
	# @godep go tool cover -html=$(CWD)/coverage.out -o $(CWD)/coverage.html
	@godep go tool cover -html=coverage.out -o $(CWD)/coverage.html
	@open coverage.html

detail:
	@echo "------------------"
	@echo " detailed report"
	@echo "------------------"
	@gocov convert coverage.out | gocov report
	# @gocov convert $(CWD)/coverage.out | gocov report

report: test detail html
