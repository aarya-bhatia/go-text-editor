package internal

func (this *Application) OpenFile(filename string) error {
	var file *File = NewFile(filename)
	this.Files = append(this.Files, file)
	

	return nil
}
