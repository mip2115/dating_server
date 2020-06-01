package scripts

func LoadDB() error {
	err := CreateUsers()
	if err != nil {
		return err
	}
	err = CreateImages()
	if err != nil {
		return err
	}
	return nil
}

func TearDB() error {
	err := DropCollectionUsers()
	if err != nil {
		return err
	}
	return nil
}
