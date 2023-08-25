package container

import (
	"fmt"
	"strings"
)

// ExchangeName returns names with postfix for stage branch deployments.
func ExchangeName(exchangeName string, env Env, branch string) string {
	if env != Stage || branch == "master" {
		return exchangeName
	}

	return fmt.Sprintf("%s_%s", exchangeName, strings.ToLower(branch))
}
