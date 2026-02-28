.PHONY: build run clean install

build:
	cargo build --release

run: build
	./target/release/nkrypt-desktop-bootstrapper

dev:
	cargo run

clean:
	cargo clean

install: build
	install -Dm755 target/release/nkrypt-desktop-bootstrapper $(DESTDIR)/usr/local/bin/nkrypt-desktop-bootstrapper
