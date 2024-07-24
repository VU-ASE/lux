# template: Go Service

This repository serves as a template repository which can be cloned to get you set up quickly. It can be used to write utilities, modules and packages in Go. The following will be set up for you when you use this template:

- A *Makefile* following the spec of our multi-repo setup using Meta (with the essential targets: `build`, `start`, `test`, `clean`)
- Github CI workflows and branch protections
    - Including Google's release-please action to release conveniently
    - On every push, the linter and all automated tests are run
    - Direct pushes to main are restricted, a PR with passing status checks (lint + test) is enforced
- A correct gitignore
- A base Go module definition with zerolog installed (see *src/main.go*)

## Getting started

- Press the top right button in Github to use this repository as a template
- After cloning, make sure to do the following:
    - Open *go.mod* and update the `module` field. Replace `ENTERYOURMODULENAMEHERE` with the name of your module
    - Open the *Makefile* and update the `BINARY_NAME` field with the name of your module
