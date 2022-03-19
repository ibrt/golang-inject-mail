# golang-inject-mail
[![Go Reference](https://pkg.go.dev/badge/github.com/ibrt/golang-inject-mail.svg)](https://pkg.go.dev/github.com/ibrt/golang-inject-mail)
![CI](https://github.com/ibrt/golang-inject-mail/actions/workflows/ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/ibrt/golang-inject-mail/branch/main/graph/badge.svg?token=BQVP881F9Z)](https://codecov.io/gh/ibrt/golang-inject-mail)

Mail module for the [golang-inject](https://github.com/ibrt/golang-inject) framework, with out-of-the-box SMTP and AWS 
SES integrations.

### Developers

Contributions are welcome, please check in on proposed implementation before sending a PR. You can validate your changes
using the `./test.sh` script.

```bash
$ MAILHOG_SMTP_PORT='4832' MAILHOG_UI_PORT='4833' MAILHOG_VERSION='v1.0.1' ./test.sh
```