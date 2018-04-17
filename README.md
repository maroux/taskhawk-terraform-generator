# Taskhawk Terraform Generator

.. image:: https://travis-ci.org/Automatic/taskhawk-terraform-generator.svg?branch=master
    :target: https://travis-ci.org/Automatic/taskhawk-terraform-generator

[Taskhawk](https://github.com/Automatic/taskhawk) is a replacement for celery that works on AWS SQS/SNS, while
keeping things pretty simple and straight forward. 

[Taskhawk-Terraform](https://github.com/Automatic/taskhawk-terraform) provides custom [Terraform](https://www.terraform.io/) modules for deploying Taskhawk infrastructure.

Taskhawk-Terraform-Generator is a utility that makes the process of managing Terraform modules easier by abstracting 
away details about Terraform.

## Usage 

### Installation

Download the latest version of the release from [Github releases](https://github.com/Automatic/taskhawk-terraform-generator/releases) - 
it's distributed as a zip containing a Go binary file.

### Configuration

Configuration is specified as a JSON file. Run 

```sh
./taskhawk-terraform-generator config-file-structure
```

to get the sample configuration file.

**Advanced usage**: The config *may* contain references to other terraform resources, as long as they resolve to 
an actual resource at runtime. 

### How to use

Run 

```sh
./taskhawk-terraform-generator apply-config <config file path>
```

to create Terraform modules. The module is named `taskhawk` by default in the current directory.

Re-run on any changes.

## Release Notes

[Github Releases](https://github.com/Automatic/taskhawk-terraform-generator/releases)

## How to publish


```sh
make clean build

cd bin/linux-amd64 && zip taskhawk-terraform-generator-linux-amd64.zip taskhawk-terraform-generator; cd -
cd bin/darwin-amd64 && zip taskhawk-terraform-generator-darwin-amd64.zip taskhawk-terraform-generator; cd -
```

Upload to Github and attach the zip files created in above step.