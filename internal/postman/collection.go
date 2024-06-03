package postman

func NewCollection(name string) CollectionJson {
	return CollectionJson{
		Info: Info{
			Name:   name,
			Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item: make([]interface{}, 0),
	}
}

func (c *CollectionJson) AddItem(item ItemGroup) {
	c.Item = append(c.Item, item)
}
