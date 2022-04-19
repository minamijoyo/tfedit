package schema

type Resource map[string]interface{}
type ImportIDFunc func(r Resource) string

type Dictionary struct {
	importIDMap map[string]ImportIDFunc
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		importIDMap: make(map[string]ImportIDFunc),
	}
}

func (d *Dictionary) RegisterImportIDFunc(resourceType string, f ImportIDFunc) {
	d.importIDMap[resourceType] = f
}

func (d *Dictionary) ImportID(resourceType string, r Resource) string {
	f := d.importIDMap[resourceType]
	return f(r)
}

var defaultDictionary = NewDictionary()

func RegisterImportIDFunc(resourceType string, f ImportIDFunc) {
	defaultDictionary.RegisterImportIDFunc(resourceType, f)
}

func ImportID(resourceType string, r Resource) string {
	return defaultDictionary.ImportID(resourceType, r)
}