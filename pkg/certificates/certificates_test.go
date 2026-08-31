package certificates

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestEmbeddedCertificateBundle(t *testing.T) {
	const expectedSHA256 = "f66dff1bdf8f96060b8177976f8b7d9254bc89bc4db933d769f7384d28480bc9"
	actualSHA256 := fmt.Sprintf("%x", sha256.Sum256(certificateBundle))
	if actualSHA256 != expectedSHA256 {
		t.Fatalf("certificate bundle SHA-256 = %s", actualSHA256)
	}
	pool, err := Pool()
	if err != nil {
		t.Fatal(err)
	}
	if rootCount := len(pool.Subjects()); rootCount != 121 {
		t.Fatalf("embedded root count = %d", rootCount)
	}
}
