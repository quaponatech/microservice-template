# A microservice template
This projects targets to provide a template for a go based microservice. Things to be adapted are marked with a `TODO`
within this repository.

## Building
The following dependencies are needed to use the full capacities of the template:

- [acbuild](https://github.com/containers/build)
- [curl](https://github.com/curl/curl)
- [docker](https://github.com/docker/docker)
- [go](https://github.com/golang/go)
- [gometalinter](https://github.com/alecthomas/gometalinter)
- [make](https://www.gnu.org/software/make/)
- [mockgen](https://github.com/golang/mock)
- [protobuf go](https://github.com/golang/protobuf)
- [protobuf](https://github.com/google/protobuf)
- [rkt](https://github.com/coreos/rkt)
- [unzip](http://www.info-zip.org/UnZip.html)

Building the main service library can be done with a simple `make`. For verbose builds run `make VERBOSE=1`. The main
command line client can be build with `go build main.go` within the root directory of the repository.

### Common make targets
Make targets used for development are:

- `make`: Builds the statically linked microservice into `./deploy/main`
- `make doc`: Run source and package documentation server and open it in the browser
- `make utest`: Run the unit tests
- `make mtest`: Run the module tests
- `make itest`: Run the integration tests
- `make bench`: Run the benchmarks
- `make docker`: Build a local docker image as tarball in `./deploy/microservice-template.tar`
- `make dockerload`: Load the local docker image into the local running docker server
- `make rk`: Build a local rkt image as aci in `./deploy/microservice.aci`
- `make clean`: Cleans the whole working directory

## Installation
All binaries can be easily installed with: `make install`, which should install the library and the executable `service`
to your `$GOBIN` directory. The `service` executable is the main interface of the service. For more information simply
execute `service -h` or `go run main.go -h`.

## Testing
The unit tests within this template should be fully mocked. If you want to run the integration tests the following
dependencies needs to be fulfilled.

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [kubernetes](https://kubernetes.io/docs/getting-started-guides/kubeadm/)
