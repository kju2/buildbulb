require "build-bulb"
require "pp"

module BuildBulb

    ID = "id"
    SUCCESS = "SUCCESS"
    THREE_DAYS_AGO = (Time.now - 3 * 24 * 3600).to_i
    ONE_DAY_AGO = (Time.now - 1 * 24 * 3600).to_i

    RSpec.describe Project do
        describe ".new" do
            context "given the stale update protection is inactive" do
                context "and the last update was three days ago" do
                    it "returns actual status." do
                        project = Project.new(ID, SUCCESS, SUCCESS, nil)
                        project.actual_status = Status::SUCCESS
                        expect(project.status).to eq(Status::SUCCESS)
                    end
                end

                context "and the last update was one day ago" do
                    it "returns actual status." do
                        project = Project.new(ID, SUCCESS, SUCCESS, nil)
                        project.actual_status = Status::SUCCESS
                        expect(project.status).to eq(Status::SUCCESS)
                    end
                end
            end

            context "given the stale update protection is active " do
                context "and the last update was three days ago" do
                    it "returns status unknown." do
                        project = Project.new(ID, SUCCESS, SUCCESS, THREE_DAYS_AGO)
                        expect(project.status).to eq(Status::UNKNOWN)
                    end
                end

                context "and the last update was one day ago" do
                    it "returns actual status." do
                        project = Project.new(ID, SUCCESS, SUCCESS, ONE_DAY_AGO)
                        expect(project.status).not_to eq(Status::UNKNOWN)
                    end
                end
            end

        end
        describe "#status" do
            context "when created with the same expected and actual status" do
                it "returns success." do
                    project = Project.new("id", "SUCCESS", "SUCCESS", Time.now.to_i)
                    expect(project.status).to eq(Status::SUCCESS)
                end
            end
            context "when created with different expected and actual status" do
                it "returns the actual status." do
                    project = Project.new("id", "SUCCESS", "UNKNOWN", Time.now.to_i)
                    expect(project.status).to eq(Status::UNKNOWN)
                end
            end
        end
    end
end
