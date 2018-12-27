protoc:
	docker run -v `pwd`:/defs namely/protoc-all -d proto/src -l go -o .
