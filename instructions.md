Develop a Go application (`tf_each`) that significantly improves Terraform code maintainability and reduces redundancy by automatically refactoring resource blocks to leverage `for_each` and variables. The application should operate on an input Terraform file (`input.tf`) and generate a refactored Terraform file (`convert/generated.tf`) and corresponding `.tfvars` files (`convert/{resource_type}.tfvars`).

**Detailed Specifications:**
Golang version 1.24, package hashicorp/hcl latest version
1.  **Terraform Parsing and Analysis:**
    * Implement a robust Terraform parser that can accurately extract resource blocks, attributes, and nested structures. Consider using an existing Terraform parser library for better accuracy (e.g., hashicorp/hcl)
    * Develop a resource grouping mechanism that identifies and groups resources based on their `resource_type` (e.g., `aws_instance`, `aws_security_group`).
    * Ensure the parser can handle common Terraform syntax variations and potential inconsistencies.

2.  **Refactoring Logic:**
    * Implement the core logic to convert grouped resources into `for_each` constructs.
    * Generate variable definitions (lists of objects) that represent the original resource configurations.
    * Modify resource attribute references to use `each.value.attribute_name` syntax.
    * Handle resource attributes of various Terraform types (strings, numbers, booleans, lists, maps).
    * If a resource type only appears once, do not attempt to convert it, and skip it.
    * Create a map to store the original values of each resource, that will be used to create the .tfvars file.

3.  **Output Generation:**
    * Create a new directory named `convert/` to store the generated files.
    * Generate a refactored Terraform file (`convert/{resource_type}.tf`) containing the modified resource definitions.
    * Generate a `.tfvars` file for each resource type that is converted.
    * Populate the `.tfvars` files with the original resource configurations, formatted correctly for the generated variables.
    * Ensure the data types stored in the .tfvars file match the data types of the variables created.
    * The .tfvars files should be created in the `convert/` directory, and be named `{resource_type}.tfvars`.

4.  **Input and Usage:**
    * The application should accept the input Terraform file path as a command-line argument.
    * Provide clear usage instructions if no input file is provided or if an invalid file path is given.
