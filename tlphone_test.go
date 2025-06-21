package tlphone_test

import (
	"testing"

	tlphone "github.com/deepakpadukone20/tlphone"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		input      string
		expectKey0 string
		expectKey1 string
		expectKey2 string
	}{
		{"ತುಂಬಾ", "03B", "03B", "053B"},
		{"ಮಕ್ಕಳು", "MKL", "MKL1", "MK2L15"},
		{"ಬಂಗಾರಾ", "B3KR", "B3KR", "B3KR"},
		{"ಅನುಗ್ರಹ", "ANKRH", "ANKRH", "AN5KRH"},
		{"ವೃತ್ತಿ", "VR0", "VR0", "VR04"},
		{"ಅಧ್ಯಕ್ಷ", "A0YKS", "A0YKS1", "A0YKS1"},
	}

	p := tlphone.New()
	for _, test := range tests {
		k0, k1, k2 := p.Encode(test.input)
		if k0 != test.expectKey0 {
			t.Errorf("Key0 mismatch for input '%s': got=%s want=%s", test.input, k0, test.expectKey0)
		}
		if k1 != test.expectKey1 {
			t.Errorf("Key1 mismatch for input '%s': got=%s want=%s", test.input, k1, test.expectKey1)
		}
		if k2 != test.expectKey2 {
			t.Errorf("Key2 mismatch for input '%s': got=%s want=%s", test.input, k2, test.expectKey2)
		}
	}
}
