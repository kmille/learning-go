all:
	CGO_ENABLED=0 go build ./polls.go 
install:
	ssh c mkdir -p /opt/polls/answeres
	scp -r questions.yaml html polls c:/opt/polls
