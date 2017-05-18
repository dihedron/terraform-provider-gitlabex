BINARY=terraform-provider-gitlabx
TEST_ENV := GITLAB_TOKEN=<admin_token> GITLAB_BASE_URL=http://localhost/api/v3

.DEFAULT_GOAL: $(BINARY)

$(BINARY):
	go build -o bin/$(BINARY)

test:
	go test -v

docker_test:
	$(TEST_ENV) go test -v
