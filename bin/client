require 'socket'

project = ARGV[0]
status = ARGV[1]
port = ARGV[2] ? ARGV[2] : 4712

socket = TCPSocket.new('127.0.0.1', port)
data = <<eos
{
    "name": "#{project}",
    "url": "job/#{project}/",
    "build": {
        "full_url": "http://localhost:8080/job/#{project}/48/",
        "number": 48,
        "phase": "FINALIZED",
        "status": "#{status}",
        "url": "job/#{project}/48/",
        "scm": {
            "url": "git@github.com:jenkinsci/#{project}.git",
            "branch": "origin/master",
            "commit": "4886d1ff4821879410f4f4a93168e6cc179a8eb3"
        },
        "artifacts": {
            "test.jar": [
                "http://localhost:8080/job/test/48/artifact/target/test.jar"
            ],
            "test.hpi": [
                "http://localhost:8080/job/test/48/artifact/target/test.hpi"
            ]
        }
    }
}
eos
socket.write data
socket.close
