.PHONY: go
go:
	rm  -f *go*.tar.gz
	curl -OL https://golang.org/dl/go1.22.5.linux-amd64.tar.gz
	rm -rf /usr/local/go
	tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
	export PATH=$$PATH:/usr/local/go/bin
	echo "export PATH=$$PATH:/usr/local/go/bin" >> ~/.profile
	go version
	rm  -f *go*.tar.gz

.PHONY: migrations
migrations:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=sqlite/dump.db \
	goose -dir migrations up

.PHONY: run
run: go migrations
	cd ~/mntk
	go build cmd/main.go
