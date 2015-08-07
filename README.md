#BUILDBULB

A server application to turn a LIFX light bulb into an extreme feedback device for multiple Jenkins projects.

## Install

```shell
go get -u github.com/kju2/buildbulb
```

The buildbulb command will be available at ${GOPATH}/bin/.

## Usage

```shell
./buildbulb 
```

Commandline parameters to configure the application are work in progress.

##Dependencies

- Go 1.4.+
- LIFX bulb (http://www.lifx.com/)

