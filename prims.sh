#!/bin/bash

# if in the VM we can guarantee where the project code is at /opt/prims-mst
# in the vm, we also add the go binaries to our path
if [ "$HOSTNAME" = prims-mst ]; then
  cd /opt/prims-mst
  GOPATH=/opt/go
  PATH=/opt/go/bin:$PATH
fi

# run tests
go test -v .

# run benchmarks 
go test -bench=.
