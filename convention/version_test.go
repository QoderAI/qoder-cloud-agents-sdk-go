package convention

import "testing"

// Both version reports must come from packageVersion: the previous hardcoded
// literal drifted from the released tag without any test noticing.
func TestVersionHeadersReportPackageVersion(t *testing.T) {
	if got, want := getDefaultHeaders()["User-Agent"], "qca-go/"+packageVersion; got != want {
		t.Errorf("User-Agent = %q, want %q", got, want)
	}
	if got := getPlatformProperties()["X-Qoder-Package-Version"]; got != packageVersion {
		t.Errorf("X-Qoder-Package-Version = %q, want %q", got, packageVersion)
	}
}
