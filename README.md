# buildbulb

A server application to turn a LIFX light bulb into an extreme feedback device for multiple Jenkins projects.

## Install

```shell
go get -u github.com/kju2/buildbulb
```

The buildbulb command will be available at ${GOPATH}/bin/.

## Usage

```shell
./buildbulb --bulbName=BuildBulb --port=8080 --jobsFilePath=/path/to/load/and/persist/jobs --apiKey=<your-lifx-api-key>
```

All parameters are optional:
- Default bulbName: BuildBulb
- Default port: 8080
- If jobsFilePath isn't set, the job status is lost when the program closes.
- if you provide no LIFX API key, the application will try to find your lamp in the local network - otherwise the LIFX cloud services are used for controlling the light.

## Jenkins Setup

- Install the [https://wiki.jenkins-ci.org/display/JENKINS/Notification+Plugin](Jenkins Notification Plugin).
- Configure Jenkins Notification Plugin to send the build status as JSON objects over HTTP the the server, e.g.
  - Format: JSON
  - Protocol: HTTP
  - Event: Finalized
  - URL: http://serverip:port/notify
- Build the job.

##Dependencies

- Go 1.4.+
- LIFX bulb (http://www.lifx.com/)

