package gojsondiff

type Iterator interface {
	SortObjectFields() bool
	EnterRoot(size int) bool
	ExitRoot() bool
	EnterObject(name string, size int) bool
	ExitObject(name string) bool
	EnterArray(name string, size int) bool
	ExitArray(name string) bool
	VisitSame(name string, d Same) bool
	VisitAdded(name string, d Added) bool
	VisitModified(name string, d Modified) bool
	VisitDeleted(name string, d Deleted) bool
}

type NullIterator struct{}

func (i *NullIterator) SortObjectFields() bool                     { return false }
func (i *NullIterator) EnterRoot(size int) bool                    { return false }
func (i *NullIterator) ExitRoot() bool                             { return false }
func (i *NullIterator) EnterObject(name string, size int) bool     { return false }
func (i *NullIterator) ExitObject(name string) bool                { return false }
func (i *NullIterator) EnterArray(name string, size int) bool      { return false }
func (i *NullIterator) ExitArray(name string) bool                 { return false }
func (i *NullIterator) VisitSame(name string, d Same) bool         { return false }
func (i *NullIterator) VisitAdded(name string, d Added) bool       { return false }
func (i *NullIterator) VisitModified(name string, d Modified) bool { return false }
func (i *NullIterator) VisitDeleted(name string, d Deleted) bool   { return false }
