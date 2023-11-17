package cmd_test

import (
	"bytes"
	"testing"

	"github.com/abemedia/appcast/pkg/cmd"
	"github.com/abemedia/appcast/pkg/crypto/dsa"
	"github.com/abemedia/appcast/pkg/crypto/ed25519"
	"github.com/abemedia/appcast/pkg/crypto/pgp"
	"github.com/abemedia/appcast/pkg/secret"
)

func TestKeysCreateCmd(t *testing.T) {
	t.Setenv("APPCAST_PATH", t.TempDir())

	for _, s := range []string{"dsa", "ed25519", "pgp"} {
		_, err := secret.Get(s + "_key")
		if err == nil {
			t.Fatalf("should not have %s key: %s", s, err)
		}
	}

	err := cmd.Execute("", []string{"keys", "create", "--name", "test", "--email", "test@example.com"})
	if err != nil {
		t.Fatal(err)
	}

	dsaKey, err := secret.Get("dsa_key")
	if err != nil {
		t.Fatalf("should have created dsa key: %s", err)
	}
	_, err = dsa.UnmarshalPrivateKey(dsaKey)
	if err != nil {
		t.Fatalf("should be valid dsa key: %s", err)
	}

	edKey, err := secret.Get("ed25519_key")
	if err != nil {
		t.Fatalf("should have created ed25519 key: %s", err)
	}
	_, err = ed25519.UnmarshalPrivateKey(edKey)
	if err != nil {
		t.Fatalf("should be valid ed25519 key: %s", err)
	}

	pgpKey, err := secret.Get("pgp_key")
	if err != nil {
		t.Fatalf("should have created pgp key: %s", err)
	}
	_, err = pgp.UnmarshalPrivateKey(pgpKey)
	if err != nil {
		t.Fatalf("should be valid pgp key: %s", err)
	}

	// Run again to ensure existing keys aren't overwritten.
	err = cmd.Execute("", []string{"keys", "create"})
	if err != nil {
		t.Fatal(err)
	}

	if k, _ := secret.Get("dsa_key"); !bytes.Equal(dsaKey, k) {
		t.Fatal("should not have regenerated dsa key")
	}

	if k, _ := secret.Get("ed25519_key"); !bytes.Equal(edKey, k) {
		t.Fatal("should not have regenerated ed25519 key")
	}

	if k, _ := secret.Get("pgp_key"); !bytes.Equal(pgpKey, k) {
		t.Fatal("should not have regenerated pgp key")
	}
}
