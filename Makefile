default: test

setup:
	go get github.com/tools/godep
	godep restore

getdumpfiles:
	mkdir data
	curl -O http://download.geonames.org/export/dump/cities1000.zip
	curl -O http://download.geonames.org/export/dump/alternateNames.zip
	unzip cities1000.zip -d data
	rm cities1000.zip
	unzip alternateNames.zip -d data
	rm alternateNames.zip
	rm data/iso-languagecodes.txt

test:
	godep go vet ./...
	godep go test ./... -cover
