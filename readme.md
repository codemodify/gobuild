# demo
- https://codemodify.github.io/gobuild

# install
- `go install github.com/codemodify/gobuild@latest`

# build
- `cd MY-GO-PROJECT`
- `gobuild`
	- it understands standard folder layout with `go.mod` being in (https://github.com/golang-standards/project-layout)
		- the current folder
		- the `PKG` folder
		- the `SRC` folder

# customize
- `cd MY-GO-PROJECT`
- `gobuild gen`
	- generates `.gobuild-config`  file with platforms to remove / add and other tweaks
	- generates `.gobuild-version` file with the version to set (`-X main.version=` for `-ldflags`)
	- generates `.gobuild-binary`  file with the binary file name
	- these files can be added to the source control for custom builds

# purpose
- this is a tool designed for use in projects that
	- don't require custom build setups
	- fast prototyping and distribution for different platforms

- why to look elsewhere
	- if you need `docker/buildx` for cross builds in complex scenarios
	- if you use CGO that will use custom APIs (ex: Xlib for linux, and other things for other OSes)
	- custom C/C++ cross-compiler toolchains
