# Terraform Provider Protogit

Run the following command to build the provider

```shell
make build
```


## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, get inside `examples/topic_schema` folder.

```shell
cd examples/topic_schema
```

Now, either set the `TF_VAR_git_password` env variable or set the `git_password` in the `terraform.tfvars` file.

```shell
export TF_VAR_git_password="<access_token>"
```

Finally, run the following commands to initialize the workspace and check the sample output.

```shell
terraform init
terraform plan
```
