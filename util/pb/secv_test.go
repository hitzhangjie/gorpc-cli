package pb

import "testing"

func TestLocateSECVProto(t *testing.T) {
	fp, err := LocateSECVProto()
	if err != nil {
		t.Fatalf("locate validate.proto error = %v", err)
	}
	t.Logf("locate validate.proto in = %s", fp)
}
