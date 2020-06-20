all:
	@echo "------------- [gorpc create] --------------"
	gorpc create -protofile=helloworld.proto -protocol=gorpc

.PHONY: clean

clean:
	rm -rf greeter

compile:
	@echo "------------- [compile gorpc] -------------"
	cd ../../../hitzhangjie-cmdline/
	make clean && make && make install
	@cd -
