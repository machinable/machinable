package mongo

// Project definition documents
func (d *Datastore) AddDefDocument(project, path string, fields map[string]interface{}) (string, error) {
	return "", nil
}

func (d *Datastore) ListDefDocuments(project, path string, limit, offset int, filter map[string]interface{}) ([]map[string]interface{}, error) {
	return nil, nil
}

func (d *Datastore) GetDefDocument(project, path, documentID string) (map[string]interface{}, error) {
	return nil, nil
}

func (d *Datastore) DeleteDefDocument(project, path, documentID string) error {
	return nil
}

func (d *Datastore) DropAllDefDocuments(project, path string) error {
	return nil
}
