package importer

type Importer interface {
	Import(interface{}, bool, interface{}) error
	Close()
}
