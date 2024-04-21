idltrans:
	kitex -service Transport -module github.com/llsw/ikunnet -gen-path ./internal/kitex_gen  ./internal/idl/transport.proto
	go mod tidy