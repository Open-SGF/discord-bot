package infra

import (
	"discord-bot/pkg/infra/customconstructs"
	"discord-bot/pkg/shared/resource"
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AppStackProps struct {
	awscdk.StackProps
	AppEnv string
}

func NewStack(scope constructs.Construct, id string, props *AppStackProps) awscdk.Stack {
	stackName := resource.NewNamer(props.AppEnv, id)

	stack := awscdk.NewStack(scope, jsii.String(stackName.FullName()), &props.StackProps)

	commonEnvVars := map[string]*string{
		"LOG_LEVEL": jsii.String("debug"),
		"LOG_TYPE":  jsii.String("json"),
	}

	workerFunctionName := resource.NewNamer(props.AppEnv, "Worker")

	workerSSMPath := "/discord-bot/" + workerFunctionName.FullName()

	workerFunction := customconstructs.NewGoLambdaFunction(stack, jsii.String(workerFunctionName.Name()), &customconstructs.GoLambdaFunctionProps{
		CodePath:     jsii.String("./cmd/worker"),
		FunctionName: jsii.String(workerFunctionName.FullName()),
		Environment: mergeMaps(commonEnvVars, map[string]*string{
			"SSM_PATH": jsii.String(workerSSMPath),
		}),
	})

	//nolint:staticcheck
	workerFunction.Function.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("ssm:GetParameter", "ssm:GetParametersByPath"),
		Resources: jsii.Strings(fmt.Sprintf("arn:aws:ssm:%s:%s:parameter%s*", *awscdk.Aws_REGION(), *awscdk.Aws_ACCOUNT_ID(), workerSSMPath)),
	}))

	workerScheduleRule := awsevents.NewRule(stack, jsii.String("WorkerEventBridgeRule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Expression(jsii.String("cron(0 15 * * ? *)")), // every 2 hours
	})

	workerScheduleRule.AddTarget(awseventstargets.NewLambdaFunction(
		workerFunction.Function,
		&awseventstargets.LambdaFunctionProps{},
	))

	return stack
}

func mergeMaps[M ~map[K]V, K comparable, V any](maps ...M) *M {
	merged := make(M)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return &merged
}
