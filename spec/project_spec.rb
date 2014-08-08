require "build-bulb"
require "json"
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

    RSpec.describe Projects do
        describe "#update" do

            project = Project.new("existing project", "success", "unknown", 0)
            projects = Projects.new(Ignore, ProjectsMemoryMarshaller.new({project.id => project}))

            it "when the project is nil then raise an exception." do
                expect{projects.update(nil, "success")}.to raise_error(KeyError, "Project \"\" not found.")
            end

            it "when the project is unknown then raise an exception." do
                expect{projects.update("not_existing", "success")}.to raise_error(KeyError, "Project \"not_existing\" not found.")
            end

            it "when the projects exists, but the status is nil then raise an exception." do
                expect{projects.update("existing project", nil)}.to raise_error(KeyError)
            end

            it "when the projects exists, but the status is empty then raise an exception." do
                expect{projects.update("existing project", "")}.to raise_error(KeyError)
            end

            it "when the projects exists, but the status is unknown then raise an exception." do
                expect{projects.update("existing project", "blub")}.to raise_error(KeyError)
            end

            it "when project and status are valid then update the project status." do
                projects.update("existing project", "unstable")
                expect(projects["existing project"].actual_status).to eq(Status::UNSTABLE)
            end
        end
    end

    RSpec.describe ProjectsFileMarshaller do
        filename = "projects.json"
        marshaller = ProjectsFileMarshaller.new(Ignore, filename=filename)

        it "#dump" do
            project = Project.new("existing project", "success", "unknown", 0)

            marshaller.dump({project.id => project})
            expect(File.file?(filename)).to be true 
        end

        it "#load" do
            expect(File.file?(filename)).to be true
            expect(marshaller.load).not_to be_empty
        end

    end

    RSpec.describe Status do
        describe "#[]" do
            it "when nil as a status key is given then raise an exception." do
                expect{Status[""]}.to raise_error(KeyError)
            end

            it "when an empty status key is given then raise an exception." do
                expect{Status[""]}.to raise_error(KeyError)
            end
            
            it "when an unknown status key is given then raise an exception." do
                expect{Status["abort"]}.to raise_error(KeyError)
            end

            it "when a known status key is given in lower case letters then return the status object." do
                expect(Status["success"]).to eq(Status::SUCCESS)
            end

            it "when a known status key is given in mixed case letters then return the status object." do
                expect(Status["Unknown"]).to eq(Status::UNKNOWN)
            end

            it "when a known status key is given in upper case letters then return the status object." do
                expect(Status["FAILURE"]).to eq(Status::FAILURE)
            end

        end
    end
end
