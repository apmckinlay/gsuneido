package stor

import (
	. "gsuneido/util/hamcrest"
	"os"
	"testing"
)

func TestAlloc(t *testing.T) {
	hs := HeapStor(64)
	adr := hs.Alloc(12)
	Assert(t).That(adr, Equals(Adr(1)))
	adr = hs.Alloc(8)
	Assert(t).That(adr, Equals(Adr(3)))
	adr = hs.Alloc(8)
	Assert(t).That(adr, Equals(Adr(4)))
	adr = hs.Alloc(48) // requires new chunk
	Assert(t).That(adr, Equals(Adr(9)))
}

func TestData(t *testing.T) {
	hs := HeapStor(64)
	hs.Alloc(12)
	adr := hs.Alloc(12)
	buf := hs.Data(adr)
	Assert(t).That(len(buf), Equals(48)) // to end of chunk
	for i := 0; i < 12; i++ {
		buf[i] = byte(i)
	}
}

func TestMmapRead(t *testing.T) {
	ms, _ := MmapStor("stor_test.go", READ)
	buf := ms.Data(Adr(1))
	Assert(t).That(string(buf[:12]), Equals("package stor"))
	ms.Close()
}

func TestMmapWrite(t *testing.T) {
	ms, _ := MmapStor("stor_test.tmp", CREATE)
	const N = 100
	buf := ms.Data(ms.Alloc(N))
	for i := 0; i < N; i++ {
		buf[i] = byte(i)
	}
	ms.Close()

	ms, _ = MmapStor("stor_test.tmp", READ)
	buf = ms.Data(Adr(1))
	for i := 0; i < N; i++ {
		Assert(t).That(buf[i], Equals(byte(i)))
	}
	ms.Close()

	os.Remove("stor_test.tmp")
}
