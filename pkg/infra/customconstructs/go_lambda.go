package customconstructs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoLambdaFunction struct {
	awslambda.Function
}

type GoLambdaFunctionProps struct {
	CodePath     *string
	FunctionName *string
	MemorySize   *float64
	Timeout      *float64
	Handler      *string
	Environment  *map[string]*string
}

func NewGoLambdaFunction(
	scope constructs.Construct,
	id *string,
	props *GoLambdaFunctionProps,
) *GoLambdaFunction {
	if props == nil {
		panic("props is required for GoLambdaFunction")
	}
	if props.CodePath == nil || *props.CodePath == "" {
		panic("CodePath is required for GoLambdaFunction")
	}

	construct := constructs.NewConstruct(scope, id)

	memory := float64(128)
	timeout := float64(60)
	handler := "main"

	if props.MemorySize != nil {
		memory = *props.MemorySize
	}
	if props.Timeout != nil {
		timeout = *props.Timeout
	}
	if props.Handler != nil {
		handler = *props.Handler
	}

	functionName := *id
	if props.FunctionName != nil {
		functionName = *props.FunctionName
	} else {
		if stack, ok := scope.(awscdk.Stack); ok {
			functionName = *stack.StackName() + "-" + *id
		}
	}

	functionProps := &awslambda.FunctionProps{
		FunctionName: jsii.String(functionName),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		MemorySize:   jsii.Number(memory),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(timeout)),
		Code: awslambda.AssetCode_FromAsset(jsii.String("."), &awss3assets.AssetOptions{
			Bundling: &awscdk.BundlingOptions{
				Image: awscdk.DockerImage_FromBuild(jsii.String("./pkg/infra"), nil),
				Command: jsii.Strings(
					"bash", "-c",
					"go build -o /asset-output/bootstrap "+*props.CodePath+"/main.go",
				),
				Volumes: &[]*awscdk.DockerVolume{
					{
						ContainerPath: jsii.String("/cache"),
						HostPath:      jsii.String("cache"),
					},
				},
			},
		}),
		Handler: jsii.String(handler),
	}

	if props.Environment != nil {
		functionProps.Environment = props.Environment
	}

	lambdaFunction := awslambda.NewFunction(construct, id, functionProps)

	goLambda := &GoLambdaFunction{
		Function: lambdaFunction,
	}

	return goLambda
}
