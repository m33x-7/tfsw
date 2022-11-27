# Terraform Switch

A small command line utility to manage Terraform versions written in [Go](https://go.dev/). It allows you to quickly switch, add, and remove versions as needed. It should support all operating systems and architectures supported by Terraform.

## Installation

### Uninstall Terraform

If you get output similar to below Terraform has already been installed and is on your $PATH. You will need to remove any existing installations beofre installing `tfsw`:

```sh
command -v terraform
/usr/local/bin/terraform
```

### Install from binary

You can grab binary releases from https://github.com/m33x-7/tfsw/releases and add them into your `$PATH` in a place like `~/bin`.

For example, if you're on Linux, you can do this:

```sh
curl https://github.com/m33x-7/tfsw/releases/download/0.0.0/tfsw_linux_amd64.gz | gunzip -c > ~/bin/tfsw
```

### Install from source

Install Go 1.17 [manually](https://go.dev/doc/install) or via your operating systems package manager. You will also need to install [GNU Make](https://www.gnu.org/software/make/).

Then run the following in the repo:

```
make && make install
```

This will put `tfsw` in `~/bin/tfsw` which you will need to add to your `$PATH`:

```
export "${HOME}/bin:${PATH}"
```
