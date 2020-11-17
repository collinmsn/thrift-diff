#!/bin/bash

goimports -w *.go
go build .
