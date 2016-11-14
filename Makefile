fmt:
	goimports -w -l ./model
	goimports -w -l ./spider
	goimports -w -l ./engine
	goimports -w -l ./config
	goimports -w -l ./main.go
update:
	govendor add +e
	govendor remove +u
build:
	go build