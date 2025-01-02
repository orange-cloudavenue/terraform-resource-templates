# terraform-provider-resource-template

`terraform-provider-resource-template` is a command-line tool that generates Go code templates for Terraform data sources and resources.

## Installation

To install `terraform-provider-resource-template`, use the following command:

```bash
go install github.com/orange-cloudavenue/terraform-provider-resource-template@latest
```

## Usage

### Generating a data source template

To generate a template for a Terraform data source, create a file with the desired name (e.g. `my_name_datasource.go`) and the following content:

```go
//tfname: myPkgName_tf_datasource_name
package myPkgName
```

Here, `myPkgName` is the name of the Go package that will contain the data source implementation, and `my_tfname` is the name of the data source as it appears in the Terraform configuration.
See below real example.

You can now generate the file like this :

```bash
terraform-provider-resource-template -filename my_name_datasource.go
```

### Generating a resource template

To generate a template for a Terraform resource, create a file with the desired name (e.g. `my_name_resource.go`) and the following content:

```go
//tfname: myPkgName_tf_resource_name
package myPkgName
```

Here, `myPkgName` is the name of the Go package that will contain the resource implementation, and `my_tfname` is the type of resource as it appears in the Terraform configuration.
See below real example.

You can now generate the file like this :

```bash
terraform-templates -filename my_name_resource.go
```

### Real example

For creating resource `myprovider_iam_user` create file `iam_user_resource.go` with content:

```go
//tfname: iam_user
package iam
```

Then run:

```bash
terraform-templates -filename iam_user_resource.go
```

## Contributing

Pull requests are welcome! If you find a bug or would like to request a new feature, please open an issue.

Before submitting a pull request, please ensure that your changes are properly tested and that the documentation has been updated.
