package util

import "testing"

func TestNameUUIDFromBytes(t *testing.T) {
	var tests = []struct {
		name, want string
	}{
		{"R4G3_BABY", "70fb6ba4-a868-32c6-8dce-43e0c4462196"},
		{"Dinnerbone", "4d258a81-2358-3084-8166-05b9faccad80"},
		{"Notch", "b50ad385-829d-3141-a216-7e7d7539ba7f"},
		{"jeb_", "a762f560-4fce-3236-812a-b80efff0b62b"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NameUUIDFromBytes([]byte("OfflinePlayer:" + test.name))
			if got.String() != test.want {
				t.Errorf("UUID was incorrect, got: %s, want: %s.", got, test.want)
			}
		})
	}
}
