#!/usr/bin/env bash

platforms=(
  # "android/arm" # Cross-compiling executables for Android requires the Android NDK
  # "darwin/386" # cmd/go: unsupported GOOS/GOARCH pair darwin/386
  "darwin/amd64"
  # "darwin/arm" # cmd/go: unsupported GOOS/GOARCH pair darwin/arm
  # "darwin/arm64" # ...\go\go1.15\pkg\tool\windows_amd64\link.exe: running gcc failed: exit status 1
  # ...\AppData\Local\Temp\go-link-763751623\go.o: file not recognized: file format not recognized
  "dragonfly/amd64"
  "freebsd/386"
  "freebsd/amd64"
  "freebsd/arm"
  "linux/386"
  "linux/amd64"
  "linux/arm"
  "linux/arm64"
  "linux/ppc64"
  "linux/ppc64le"
  "linux/mips"
  "linux/mipsle"
  "linux/mips64"
  "linux/mips64le"
  "netbsd/386"
  "netbsd/amd64"
  "netbsd/arm"
  "openbsd/386"
  "openbsd/amd64"
  "openbsd/arm"
  "plan9/386"
  "plan9/amd64"
  "solaris/amd64"
  "windows/386"
  "windows/amd64"
)

# iterate through platforms
for platform in "${platforms[@]}"; do
  # split platform string by /
  IFS="/" read -r -a platform_split <<<"$platform"

  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}

  output_name="build/rpdSix-$GOOS-$GOARCH"

  if [ "$GOOS" = "windows" ]; then
    output_name+=".exe"
  fi

  # use subshells to execute all build commands simultaneously
  (
    echo "Building for $platform"

    if ! env GOOS="$GOOS" GOARCH="$GOARCH" go build -o $output_name; then
      echo "An error has occurred! Aborting the script execution..."
      exit 1
    fi

    echo "Done building for $platform"
  ) &
done

wait

echo "Done."

exit 0
