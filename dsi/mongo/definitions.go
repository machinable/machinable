package mongo

import (
	"bitbucket.org/nsjostrom/machinable/dsi/errors"
)

// Project definition documents
func (d *Datastore) AddDefDocument(project, path string, fields map[string]interface{}) (string, *errors.DatastoreError) {
	return "", nil
}

func (d *Datastore) ListDefDocuments(project, path string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, *errors.DatastoreError) {
	return nil, nil
}

func (d *Datastore) GetDefDocument(project, path, documentID string) (map[string]interface{}, *errors.DatastoreError) {
	return nil, nil
}

func (d *Datastore) DeleteDefDocument(project, path, documentID string) *errors.DatastoreError {
	return nil
}

func (d *Datastore) DropAllDefDocuments(project, path string) *errors.DatastoreError {
	return nil
}
