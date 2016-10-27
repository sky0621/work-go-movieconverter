package movieconverter

import "testing"

func TestInputPath(t *testing.T) {
	res := joinPath("A", "B", "C")
	if "A\\B\\C" != res {
		t.Fatal("Res: ", res)
	}
}
