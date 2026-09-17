package designsystem

import "testing"

// TestTokensCSSLiefertKopie belegt, dass ein Aufrufer, der das Ergebnis von
// TokensCSS in-place verändert, den eingebetteten Puffer nicht beschädigt —
// sonst wäre jeder weitere Aufruf im selben Prozess betroffen, nicht nur der
// des Aufrufers, der die Änderung vorgenommen hat.
func TestTokensCSSLiefertKopie(t *testing.T) {
	vorher := string(TokensCSS())

	kopie := TokensCSS()
	if len(kopie) == 0 {
		t.Fatal("TokensCSS() liefert leeren Puffer")
	}
	kopie[0] = 'X'

	nachher := string(TokensCSS())
	if nachher != vorher {
		t.Error("TokensCSS() gibt den eingebetteten Puffer selbst heraus — eine Änderung am Ergebnis wirkt sich auf spätere Aufrufe aus")
	}
}

// TestBaseCSSLiefertKopie prüft dieselbe Eigenschaft für BaseCSS.
func TestBaseCSSLiefertKopie(t *testing.T) {
	vorher := string(BaseCSS())

	kopie := BaseCSS()
	if len(kopie) == 0 {
		t.Fatal("BaseCSS() liefert leeren Puffer")
	}
	kopie[0] = 'X'

	nachher := string(BaseCSS())
	if nachher != vorher {
		t.Error("BaseCSS() gibt den eingebetteten Puffer selbst heraus — eine Änderung am Ergebnis wirkt sich auf spätere Aufrufe aus")
	}
}
