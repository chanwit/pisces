#!/bin/bash

INTEGRATION_ROOT=$(dirname "$(readlink -f "$BASH_SOURCE")")

# unit tests
echo "Unit testing ..."

sudo bats "${INTEGRATION_ROOT}/go"

if [[ $? -eq 0 ]]; then

	# recompile
	echo "Compiling ..."
	(cd $GOPATH/src/github.com/chanwit/pisces && go install . )

	GO_RESULT=$?

	if [[ GO_RESULT -eq 0 ]]; then
	  #
	  echo "Integration testing ..."
	  # run tests
	  sudo PISCES_BINARY=$GOBIN/pisces bats $INTEGRATION_ROOT
	fi

fi