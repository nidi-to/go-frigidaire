package frigidaire

import (
	"os"
	"testing"
	"time"
)

func Test_NewService(t *testing.T) {
	username := os.Getenv("FRIGIDAIRE_USERNAME")
	password := os.Getenv("FRIGIDAIRE_PASSWORD")

	sess, err := NewSession(username, password)

	if err != nil {
		t.Fatalf("Could not create service: %v", err)
	}

	if sess.Appliances == nil {
		t.Fatal("No appliances found")
	}

	if sess.Appliances == nil {
		t.Fatal("No appliances found")
	}

	if sess.Expires.After(time.Now()) {
		t.Fatal("Expired token")
	}
}
