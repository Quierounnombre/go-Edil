package nobuildtag

const (
	GO_BUILD       = "//go:build"
	PLUS_BUILD     = "// +build"
	GO_SUFFIX      = ".go"
	GO_MOD         = "go.mod"
	NAME_SEP       = "_"
	HIDDEN_DOT     = "."
	VENDOR_DIR     = "vendor"
	MIN_PARTS      = 2
	OS_LIST        = "aix android darwin dragonfly freebsd hurd illumos ios js linux nacl netbsd openbsd plan9 solaris wasip1 windows zos"
	ARCH_LIST      = "386 amd64 amd64p32 arm armbe arm64 arm64be loong64 mips mipsle mips64 mips64le mips64p32 mips64p32le ppc ppc64 ppc64le riscv riscv64 s390 s390x sparc sparc64 wasm"
	LIST_SEP       = " "
	MSG_CONSTRAINT = "build constraints hide code from the linter; remove the //go:build line"
	MSG_SUFFIX     = "GOOS/GOARCH file name suffix hides code from the linter; do not use platform-specific files"
	MSG_IGNORED    = "file is excluded from this build, so it is never checked; delete it or remove its constraint"
	MSG_HIDDEN     = "directory %q holds Go files that ./... never checks; remove it or move it into a normal package"
)
