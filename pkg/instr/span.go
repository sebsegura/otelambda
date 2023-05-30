package instr

import (
	"fmt"
	"os"
)

const (
	_lambdaName    = "AWS_LAMBDA_FUNCTION_NAME"
	_lambdaVersion = "AWS_LAMBDA_FUNCTION_VERSION"
)

func GetTracerName() string {
	return fmt.Sprintf("%s:%s", os.Getenv(_lambdaName), os.Getenv(_lambdaVersion))
}
