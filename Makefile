idltrans:
	kitex -service Transport -module github.com/llsw/ikunet ./internal/idl/transport.proto
	grep -rl 'ikunet/kitex_gen' ./kitex_gen | xargs sed -i "" "s#ikunet/kitex_gen#ikunet/internal/kitex_gen#g"
	rm -rf internal/kitex_gen
	mv kitex_gen internal/kitex_gen
	go mod tidy