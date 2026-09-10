<p align="center">
  <img src="logo.png" alt="Leaven" width="280">
</p>

# Leaven: Compile LLVM IR to Go

Leaven translates LLVM intermediate representation to Go. 
In theory, it should be able to transpile any language that has an LLVM-based compiler to Go.
But so far I’ve only used it for C.

Each LLVM instruction is translated to an equivalent statement in Go.
This produces very verbose code;
if you are looking for a tool that will convert a C codebase into maintainable Go,
Leaven isn’t it.

But it does allow you to call C code from Go without using CGo.
And I am hoping that it can produce a working Go translation of a program,
which will be a good starting point for incrementally re-translating it
(probably by hand) into idiomatic Go.

## Warning

This software is incomplete and experimental.
It does not support nearly all LLVM instructions.
It is 64-bit only (amd64, arm64, and other 8-byte-pointer hosts).

The transpiler at github.com/andybalholm/c2go produces much better results
(but it is not as automatic).

## Install

Download the archive for your OS/arch from [GitHub Releases](https://github.com/lewtec/leaven/releases), extract `leaven`, put it on `PATH`.

Builds cover the same targets as CI: linux, darwin, and windows, amd64 and arm64.

From source:

```bash
go install github.com/lewtec/leaven/cmd/leaven@latest
```

## Usage Example
(Translating `strcmp` from musl libc.)

	$ cat strcmp.c
	#include <string.h>

	int strcmp(const char *l, const char *r)
	{
		for (; *l==*r && *l; l++, r++);
		return *(unsigned char *)l - *(unsigned char *)r;
	}
	$ clang -S -emit-llvm -fno-discard-value-names strcmp.c
	$ go run ./cmd/leaven strcmp.ll
	$ clang -S -emit-llvm -fno-discard-value-names -o - strcmp.c | go run ./cmd/leaven > strcmp.go
	$ goimports -w strcmp.go
	$ cat strcmp.go
	package main

	import (
		"unsafe"

		"github.com/lewtec/leaven/libc"
	)

	func strcmp(l unsafe.Pointer, r unsafe.Pointer) int32 {
		var cmp, tobool, v6 bool
		var conv, conv1, conv3, conv5, conv6, sub int32
		var v1, v3, v5, v10, v12 byte
		var v0, v2, v4, v7, incdec_ptr, v8, incdec_ptr4, v9, v11 unsafe.Pointer
		var l_addr, r_addr unsafe.Pointer
		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = l_addr, r_addr, v0, v1, conv, v2, v3, conv1, cmp, v4, v5, conv3, tobool, v6, v7, incdec_ptr, v8, incdec_ptr4, v9, v10, conv5, v11, v12, conv6, sub

		l_addr_mem := libc.Alloca[unsafe.Pointer](1, int64(8))
		l_addr = libc.Ptr(l_addr_mem)
		defer libc.AllocaFree(libc.As[byte](libc.Ptr(l_addr_mem)))
		r_addr_mem := libc.Alloca[unsafe.Pointer](1, int64(8))
		r_addr = libc.Ptr(r_addr_mem)
		defer libc.AllocaFree(libc.As[byte](libc.Ptr(r_addr_mem)))
		*libc.As[unsafe.Pointer](l_addr) = l
		*libc.As[unsafe.Pointer](r_addr) = r
		goto for_cond

	for_cond:
		v0 = *libc.As[unsafe.Pointer](l_addr)
		v1 = *libc.As[byte](v0)
		conv = int32(int8(v1))
		v2 = *libc.As[unsafe.Pointer](r_addr)
		v3 = *libc.As[byte](v2)
		conv1 = int32(int8(v3))
		cmp = conv == conv1
		if cmp {
			goto land_rhs
		} else {
			v6 = false
			goto land_end
		}

	land_rhs:
		v4 = *libc.As[unsafe.Pointer](l_addr)
		v5 = *libc.As[byte](v4)
		conv3 = int32(int8(v5))
		tobool = conv3 != 0
		v6 = tobool
		goto land_end

	land_end:
		if v6 {
			goto for_body
		} else {
			goto for_end
		}

	for_body:
		goto for_inc

	for_inc:
		v7 = *libc.As[unsafe.Pointer](l_addr)
		incdec_ptr = libc.Ptr(libc.AddPointer[byte](libc.As[byte](v7), int(1)*1))
		*libc.As[unsafe.Pointer](l_addr) = incdec_ptr
		v8 = *libc.As[unsafe.Pointer](r_addr)
		incdec_ptr4 = libc.Ptr(libc.AddPointer[byte](libc.As[byte](v8), int(1)*1))
		*libc.As[unsafe.Pointer](r_addr) = incdec_ptr4
		goto for_cond

	for_end:
		v9 = *libc.As[unsafe.Pointer](l_addr)
		v10 = *libc.As[byte](v9)
		conv5 = int32(uint32(v10))
		v11 = *libc.As[unsafe.Pointer](r_addr)
		v12 = *libc.As[byte](v11)
		conv6 = int32(uint32(v12))
		sub = conv5 - conv6
		return sub
	}

## Release

[GoReleaser](https://goreleaser.com) + [svu](https://github.com/caarlos0/svu). Archives and checksums only (no Homebrew, Docker, or packages). Tags have no `v` prefix ([`.svu.yml`](.svu.yml)).

```bash
mise release          # next (svu) + goreleaser (needs GITHUB_TOKEN)
mise release patch    # or major | minor | next
```

CI: [`.github/workflows/autorelease.yml`](.github/workflows/autorelease.yml). Push/PR runs the six-target test matrix. `workflow_dispatch` with patch/minor/major runs that matrix, then tags and publishes if every cell passed.

