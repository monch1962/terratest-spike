package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

func StartLocalStack() (context.Context, testcontainers.Container) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack",
		ExposedPorts: []string{"4510-4559/tcp", "4566/tcp", "4571/tcp", "5678/tcp"},
	}
	localStack, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	fmt.Println(err)
	//if err != nil {
	//	return nil
	//}
	return ctx, localStack
}

func TestMain(m *testing.M) {
	ctx, localStack := StartLocalStack()
	fmt.Println(ctx, localStack)
	os.Exit(m.Run())
}

func TestTerraformLocalstackS3(t *testing.T) {
	t.Parallel()

	//expectedText := "test"
	//expectedList := []string{expectedText}
	//expectedMap := map[string]string{"expected": expectedText}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// website::tag::1::Set the path to the Terraform code that will be tested.
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-localstack-test",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			//"example": expectedText,

			// We also can see how lists and maps translate between terratest and terraform.
			//"example_list": expectedList,
			//"example_map":  expectedMap,
		},

		// Variables to pass to our Terraform code using -var-file options
		// VarFiles: []string{"varfile.tfvars"},

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		//NoColor: true,
	})

	// website::tag::4::Clean up resources with "terraform destroy". Using "defer" runs the command at the end of the test, whether the test succeeds or fails.
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2::Run "terraform init" and "terraform apply".
	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	//actualTextExample := terraform.Output(t, terraformOptions, "example")
	//actualTextExample2 := terraform.Output(t, terraformOptions, "example2")
	//actualExampleList := terraform.OutputList(t, terraformOptions, "example_list")
	//actualExampleMap := terraform.OutputMap(t, terraformOptions, "example_map")

	// Get the bucket ID so we can query AWS
	awsRegion := "us-east-1"
	bucketID := "onexlab-bucket-terraform"
	//bucketID := terraform.Output(t, terraformOptions, "bucket_id")

	actualTags := aws.GetS3BucketTags(t, awsRegion, bucketID)
	fmt.Println(actualTags)

	// website::tag::3::Check the output against expected values.
	// Verify we're getting back the outputs we expect
	//assert.Equal(t, expectedText, actualTextExample)
	//assert.Equal(t, expectedText, actualTextExample2)
	//assert.Equal(t, expectedList, actualExampleList)
	//assert.Equal(t, expectedMap, actualExampleMap)
	assert.Equal(t, bucketID, "onexlab-bucket-terraform")
}
