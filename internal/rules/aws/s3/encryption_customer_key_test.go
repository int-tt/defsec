package s3

import (
	"testing"

	defsecTypes "github.com/aquasecurity/defsec/pkg/types"

	"github.com/aquasecurity/defsec/pkg/state"

	"github.com/aquasecurity/defsec/pkg/providers/aws/s3"
	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckEncryptionCustomerKey(t *testing.T) {
	tests := []struct {
		name     string
		input    s3.S3
		expected bool
	}{
		{
			name: "S3 Bucket missing KMS key",
			input: s3.S3{
				Buckets: []s3.Bucket{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Encryption: s3.Encryption{
							Metadata: defsecTypes.Metadata{},
							Enabled:  defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
							KMSKeyId: defsecTypes.String("", defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "S3 Bucket with KMS key",
			input: s3.S3{
				Buckets: []s3.Bucket{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Encryption: s3.Encryption{
							Metadata: defsecTypes.Metadata{},
							Enabled:  defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
							KMSKeyId: defsecTypes.String("some-sort-of-key", defsecTypes.NewTestMetadata()),
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
			testState.AWS.S3 = test.input
			results := CheckEncryptionCustomerKey.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckEncryptionCustomerKey.Rule().LongID() {
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
