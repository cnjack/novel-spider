fmt:
	goimports -w -l ./model
	goimports -w -l ./spider
	goimports -w -l ./job
	goimports -w -l ./tool
	goimports -w -l ./config
	goimports -w -l ./httpserv
	goimports -w -l ./main.go
update:
	govendor add +e
	govendor remove +u
build:
	go build