package stst

// PsqlTypemap is an implementation of Typemap for PostgreSQL.
type PsqlTypemap struct{}

// GetGoType .
func (p PsqlTypemap) GetGoType(databaseTypeName string) (string, bool) {
	m := map[string]string{
		"FLOAT4":  "float32",
		"FLOAT8":  "float64",
		"NUMERIC": "float64",
	}
	v, ok := m[databaseTypeName]
	return v, ok
}
