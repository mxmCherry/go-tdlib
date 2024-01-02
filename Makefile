.PHONY: default
default: third-party/td/tdlib

# aggressive
.PHONY: clean
clean:
	rm -rf third-party

# https://tdlib.github.io/td/build.html?language=Go

third-party/td:
	git clone https://github.com/tdlib/td.git \
		--branch master \
		--single-branch \
		--depth 1 \
		$@

third-party/td/tdlib: third-party/td
	rm -rf $(@D)/build
	mkdir -p $(@D)/build
	cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=../tdlib -S $(@D) -B $(@D)/build
	cmake --build $(@D)/build --target install
