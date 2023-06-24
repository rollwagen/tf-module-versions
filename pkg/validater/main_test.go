package validater

import "testing"

func Test_determineLatestVersionString(t *testing.T) {
	tests := []struct {
		name                 string
		availableVersionsRaw []string
		want                 string
		wantErr              bool
	}{
		{
			name:                 "Latest version retrieved",
			availableVersionsRaw: []string{"1.0.0", "0.9.9", "1.1.0"},
			want:                 "1.1.0",
			wantErr:              false,
		},
		{
			name:                 "Latest version retrieved with v prefix",
			availableVersionsRaw: []string{"v1.0", "v0.9.9", "v1.1.0"},
			want:                 "v1.1.0",
			wantErr:              false,
		},
		{
			name:                 "No versions returns nil",
			availableVersionsRaw: []string{},
			want:                 "nil",
			wantErr:              false,
		},
		{
			name:                 "Invalid version returns error",
			availableVersionsRaw: []string{"1.0.0", "invalid"},
			want:                 "1.0.0",
			wantErr:              false,
		},
		{
			name:                 "Mix of version strings",
			availableVersionsRaw: []string{"v1.0", "0.9.9", "v3"},
			want:                 "v3",
			wantErr:              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := determineLatestVersionString(tt.availableVersionsRaw)
			if (err != nil) != tt.wantErr {
				t.Errorf("determineLatestVersionString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("determineLatestVersionString() = %v, want %v", got, tt.want)
			}
		})
	}
}
