#!/bin/bash
echo "deepsource.sh running...."

# Run your tests and generate coverage report
go test -coverprofile=cover.out

# Install 'deepsource CLI'
curl https://deepsource.io/cli | sh

# Set DEEPSOURCE_DSN env variable from repository settings page
export DEEPSOURCE_DSN=$DEEPSOURCE_DSN

# From the root directory, run the report coverage command
./bin/deepsource report --analyzer test-coverage --key go --value-file ./cover.out
