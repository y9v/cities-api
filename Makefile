default: test

setup:
	go get github.com/mattn/gom
	gom install

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
	go vet ./...
	gom test ./... -cover
