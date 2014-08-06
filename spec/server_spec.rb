require "build-bulb"
require "logger"
require "pp"
require "socket"
#require "unit_helper"

module BuildBulb

    ADDRESS = "0.0.0.0"
    PORT = 4712

    RSpec.describe Server do
        def start_server(address, port, logger)
            Thread.new {
            server = Server.new(address, port, double("Projects").as_null_object, double("Light").as_null_object, logger)
            server.listen_indefinitely}
        end

        def send_message(address, port, message)
            socket = TCPSocket.new(address, port)
            socket.write(message)
            socket.close
        end

        describe "#listen_indefinitely" do
            context "when an invalid message is sent to the server" do
                it "the server keeps on listening." do
                    logger = instance_double(Logger)
                    allow(logger).to receive(:debug)
                    allow(logger).to receive(:info)
                    expect(logger).to receive(:warn).exactly(4).times
                    #server = Server.new(ADDRESS, PORT, double("Projects").as_null_object, double("Light").as_null_object, LOGGER)
                    #server.listen_indefinitely
                    server = start_server(ADDRESS, PORT, logger)
                    send_message(ADDRESS, PORT, "test")
                    expect(server.alive?).to be true
                    send_message(ADDRESS, PORT, "test")
                    expect(server.alive?).to be true
                    server.join(0.001)
                    server.kill
                end
            end
        end
    end

end

