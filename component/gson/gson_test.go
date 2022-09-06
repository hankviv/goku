package gson

import "testing"

type SimpleStruct struct {
	A int
	B int
}

func Test_populateStructReflect(t *testing.T) {

	var m SimpleStruct
	populateStructReflect(&m)
}

func BenchmarkPopulateReflect(b *testing.B) {
	b.ReportAllocs()
	var m SimpleStruct
	for i := 0; i < b.N; i++ {
		if err := populateStructReflect(&m); err != nil {
			b.Fatal(err)
		}
		if m.B != 42 {
			b.Fatalf("unexpected value %d for B", m.B)
		}
	}
}
