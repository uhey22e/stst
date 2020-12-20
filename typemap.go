package stst

// Typemap is an interface
type Typemap interface {
	// GetGoType returns a golang type
	GetGoType(databaseTypeName string) (string, bool)
}
