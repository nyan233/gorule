package gorule

import "testing"

func TestGoRule(t *testing.T) {
	type A0 struct {
		Name   string
		UserId string
		Tags   []string
		Config map[string]bool
	}
	args := []Argument{
		{
			Name: "a0",
			Val: A0{
				Name:   "zbh255",
				UserId: "change",
				Tags: []string{
					"new-user", "inner",
				},
				Config: map[string]bool{
					"Show":     true,
					"DataSync": true,
				},
			},
		},
		{
			Name: "a2",
			Val:  true,
		},
	}
	t.Run("Base", func(t *testing.T) {
		ok, err := ExecuteSimpleBoolExpr("(a0.Name == \"zbh255\" || a0.UserId == \"change\") || a2 ", args...)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ok)
	})
	t.Run("FuncCall", func(t *testing.T) {
		ok, err := ExecuteSimpleBoolExpr("len(a0.Tags) == 2 && a0.Tags[1] == \"inner\"", args...)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ok)
	})
	t.Run("MapVisit", func(t *testing.T) {
		ok, err := ExecuteSimpleBoolExpr("len(a0.Config) == 2 && a0.Config[\"Show\"]", args...)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ok)
	})
}

func BenchmarkGoRule(b *testing.B) {
	type A0 struct {
		Name   string
		UserId string
		Tags   []string
	}
	args := []Argument{
		{
			Name: "a0",
			Val: A0{
				Name:   "zbh255",
				UserId: "change",
				Tags: []string{
					"new-user", "inner",
				},
			},
		},
		{
			Name: "a2",
			Val:  true,
		},
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ok, err := ExecuteSimpleBoolExpr("len(a0.Tags) == 2 && a0.Tags[1] == \"inner\"", args...)
		if err != nil {
			b.Fatal(err)
		}
		_ = ok
	}
}
