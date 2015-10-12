package types

import (
	"github.com/ng-vu/graphql-go/ql"
)

type QLDirective struct {
	Name        string
	Description string
	Args        []*QLArgument
	OnOperation bool
	OnFragment  bool
	OnField     bool
}

var QLIncludeDirective = &QLDirective{
	Name:        "include",
	Description: `Directs the executor to include this field or fragment only when the "if" argument is true.`,
	Args: []*QLArgument{
		{
			Name:        "if",
			Type:        NewQLNonNull(NewQLType(ql.Boolean)),
			Description: "Included when true.",
		},
	},
	OnOperation: false,
	OnFragment:  true,
	OnField:     true,
}

var QLSkipDirective = &QLDirective{
	Name:        "skip",
	Description: `Directs the executor to skip this field or fragment when the "if" argument is true.`,
	Args: []*QLArgument{
		{
			Name:        "if",
			Type:        NewQLNonNull(NewQLType(ql.Boolean)),
			Description: "Included when true.",
		},
	},
	OnOperation: false,
	OnFragment:  true,
	OnField:     true,
}
