# Dockerhubbeat

Welcome to Dockerhubbeat.

DockerhubBeat is a small beat design to read and store informations on a specific docker images (directly read from https://hub.docker.com).

You just need to specify the repository to watch: `repository: traefik` or `repository: containous/traefik`.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/FriggaHel/dockerhubbeat`

## Getting Started with Dockerhubbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Dockerhubbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Dockerhubbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/FriggaHel/dockerhubbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Dockerhubbeat run the command below. This will generate a binary
in the same directory with the name dockerhubbeat.

```
make
```


### Run

To run Dockerhubbeat with debugging output enabled, run:

```
./dockerhubbeat -c dockerhubbeat.yml -e -d "*"
```


### Test

To test Dockerhubbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/dockerhubbeat.template.json and etc/dockerhubbeat.asciidoc

```
make update
```


### Cleanup

To clean  Dockerhubbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Dockerhubbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/FriggaHel/dockerhubbeat
cd ${GOPATH}/github.com/FriggaHel/dockerhubbeat
git clone https://github.com/FriggaHel/dockerhubbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
