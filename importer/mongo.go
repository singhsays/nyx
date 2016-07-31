package importer

import "gopkg.in/mgo.v2"

// MongoImporter imports models into a Mongodb DB.
type MongoImporter struct {
	session    *mgo.Session
	collection *mgo.Collection
}

// NewMongoImporter returns an instantiated MongoImporter based on the given
// database address, name and collection name.
func NewMongoImporter(address, name, collection string) (*MongoImporter, error) {
	session, err := mgo.Dial(address)
	if err != nil {
		return nil, err
	}
	return &MongoImporter{
		session:    session,
		collection: session.DB(name).C(collection),
	}, nil
}

// Close closes the underlying mongo session.
func (m *MongoImporter) Close() {
	m.session.Close()
}

// Import persists the given data model in the database.
// Based on the overwrite parameter, it either inserts a new document or
// upserts into the document matching the given query.
func (m *MongoImporter) Import(data interface{}, overwrite bool, query interface{}) error {
	var err error
	if overwrite {
		_, err = m.collection.Upsert(query, data)
	} else {
		err = m.collection.Insert(data)
	}
	return err
}
