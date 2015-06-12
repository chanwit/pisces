
# recompile
(cd $GOPATH/src/github.com/chanwit/pisces && go install . )

GO_RESULT=$?

if [[ GO_RESULT -eq 0 ]]; then
  # run tests
  sudo PISCES_BINARY=$GOBIN/pisces bats .
fi
