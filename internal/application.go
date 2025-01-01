package internal

const NORMAL_MODE int = 0
const INSERT_MODE int = 1
const COMMAND_MODE int = 2

type Application struct {
	QuitSignal bool
	Files      []*File
  CurrentFile *File
	StatusLine string
	Mode       int
}

func NewApplication() *Application {
	var this *Application = new(Application)
	this.QuitSignal = false
	this.Files = make([]*File, 0)
  this.CurrentFile = nil
	this.StatusLine = ""
  this.Mode = NORMAL_MODE
	return this
}

func (this *Application) OpenFile(filename string) error {
	var file *File = NewFile(filename)
	this.Files = append(this.Files, file)
	

	return nil
}

