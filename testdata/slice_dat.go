package testdata

// @Copyable
type C struct {
	M map[string]string
}

// @Copyable
// @PtrRecv true
type S struct {
	Age  int
	Name string
	Sex  *int

	S *[]*C
}

// @Copyable
type List []S
