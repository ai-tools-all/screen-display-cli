# Multi-Platform Build Strategy

## 1. Cross-Compilation Matrix
Go supports cross-compilation out-of-the-box using `GOOS` and `GOARCH` environment variables.

| Platform | GOOS | GOARCH |
| :--- | :--- | :--- |
| Linux (x64) | linux | amd64 |
| Linux (ARM) | linux | arm64 |
| macOS (Intel) | darwin | amd64 |
| macOS (Silicon) | darwin | arm64 |
| Windows (x64) | windows | amd64 |

## 2. Implementation: Makefile
Automate the build process using a `Makefile` to iterate through the platform matrix.

```makefile
BINARY_NAME=myapp
OUT_DIR=bin

.PHONY: build-all
build-all:
	@mkdir -p $(OUT_DIR)
	@for platform in "linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64"; do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output_name=$(OUT_DIR)/$(BINARY_NAME)-$$os-$$arch; \
		if [ $$os = "windows" ]; then output_name=$$output_name.exe; fi; \
		echo "Building $$output_name..."; \
		GOOS=$$os GOARCH=$$arch go build -o $$output_name ./cmd/app; \
	done
```

## 3. GitHub Actions Integration
For CI/CD, use a matrix strategy to build in parallel:
```yaml
strategy:
  matrix:
    os: [linux, darwin, windows]
    arch: [amd64, arm64]
    exclude:
      - os: windows
        arch: arm64
```
