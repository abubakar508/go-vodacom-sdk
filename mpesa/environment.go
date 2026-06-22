package mpesa

// Environment selects the M-Pesa OpenAPI URL path and default platform public key.
type Environment string

const (
	EnvironmentSandbox Environment = "sandbox"
	EnvironmentOpenAPI Environment = "openapi"

	// EnvironmentProduction is an alias for the live OpenAPI environment.
	EnvironmentProduction = EnvironmentOpenAPI
)

const defaultHost = "openapi.m-pesa.com"

// BasePath returns the first OpenAPI path segment used by this environment.
func (e Environment) BasePath() string {
	if e == EnvironmentOpenAPI {
		return "openapi"
	}
	return "sandbox"
}

// Valid reports whether the environment is supported by this SDK.
func (e Environment) Valid() bool {
	return e == EnvironmentSandbox || e == EnvironmentOpenAPI
}
