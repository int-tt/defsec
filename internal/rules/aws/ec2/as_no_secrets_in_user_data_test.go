package ec2

import (
	"testing"

	defsecTypes "github.com/aquasecurity/defsec/pkg/types"

	"github.com/aquasecurity/defsec/pkg/state"

	"github.com/aquasecurity/defsec/pkg/providers/aws/ec2"
	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/stretchr/testify/assert"
)

func TestASCheckNoSecretsInUserData(t *testing.T) {
	tests := []struct {
		name     string
		input    ec2.EC2
		expected bool
	}{
		{
			name: "Launch template with sensitive info in user data",
			input: ec2.EC2{
				LaunchTemplates: []ec2.LaunchTemplate{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Instance: ec2.Instance{
							Metadata: defsecTypes.NewTestMetadata(),
							UserData: defsecTypes.String(`
							export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
							export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
							export AWS_DEFAULT_REGION=us-west-2
							`, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Launch template with no sensitive info in user data",
			input: ec2.EC2{
				LaunchTemplates: []ec2.LaunchTemplate{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Instance: ec2.Instance{
							Metadata: defsecTypes.NewTestMetadata(),
							UserData: defsecTypes.String(`
							export GREETING=hello
							`, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.EC2 = test.input
			results := CheckASNoSecretsInUserData.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckASNoSecretsInUserData.Rule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
