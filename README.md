<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/go-nf?layout=left&font=1.1&bg=geometric&icon=dinkie-icons:file-font-filled&icon.height=88&icon.color=%2306BAFE)

---

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">
  <a href="https://github.com/lrstanley/go-nf/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/go-nf?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/go-nf/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/go-nf?style=flat-square">
  </a>



  <a href="https://github.com/lrstanley/go-nf/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/actions/workflow/status/lrstanley/go-nf/test.yml?branch=master&label=test&style=flat-square">
  </a>


  <a href="https://codecov.io/gh/lrstanley/go-nf">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/go-nf/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/go-nf">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/go-nf?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/go-nf">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/go-nf?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/go-nf/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/go-nf/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/go-nf/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/go-nf/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/go-nf/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/go-nf?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/go-nf/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [Features](#sparkles-features)
  - [Usage](#gear-usage)
  - [Examples](#clap-examples)
    - [Simple example](#simple-example)
    - [Package manager example](#package-manager-example)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [License](#balance_scale-license)
<!-- template:end:toc -->

`go-nf` is a Go package that provides a type-safe way to access Nerd Font glyphs.

## :sparkles: Features

> [!WARNING]
> The scope and design of this module is still a work in progress, and may change
> before it is considered stable, and a `v1.0.0` release is made.

- :heavy_check_mark: Reference Nerd Font glyphs by their name, rather than hardcoding
  them as strings in your project. Ensures type safety, validation, and better
  documentation.
- :heavy_check_mark: Split into multiple packages, per Nerd Font class, to reduce
  binary size and improve organization.
- :heavy_check_mark: Detect if Nerd Fonts are installed on the system, to prevent
  using glyphs that are not available.
  - Note that this is not a perfect solution. Just because Nerd Fonts are installed,
    doesn't mean the font in question is actively being used in the associated
    terminal emulator.
- :heavy_check_mark: Helpers for iterating over all glyphs, by class, ID, and more.

---

## :gear: Usage

<!-- template:begin:goget -->
<!-- do not edit anything in this "template" block, its auto-generated -->
```console
go get -u github.com/lrstanley/go-nf@latest
```
<!-- template:end:goget -->

## :clap: Examples

### Simple example

- Static glyph references and dynamic lookup by file extension.
- [Example source](./_examples/simple/main.go)

### Package manager example

- Bubble Tea TUI simulating package installation with glyphs per package type.
- [Example source](./_examples/package-manager/main.go)

![package-manager example](./_examples/package-manager/demo.gif)

---

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :raising_hand_man: Support & Assistance

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
     guidelines on ensuring everyone has the best experience interacting with
     the community.
* :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
     guidelines for tips on how to ask the right questions.
* :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/go-nf/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :handshake: Contributing

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
     on ensuring everyone has the best experience interacting with the
    community.
* :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
     issues/a guide on submitting pull requests and helping out.
* :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/go-nf/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :balance_scale: License

```
MIT License

Copyright (c) 2026 Liam Stanley <liam@liam.sh>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_
<!-- template:end:license -->
