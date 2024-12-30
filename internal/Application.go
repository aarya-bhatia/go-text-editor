package internal

const NORMAL_MODE int = 0
const INSERT_MODE int = 1
const COMMAND_MODE int = 2

type Application struct {
	QuitSignal bool
	Files      []*File
	StatusLine string
	Mode       int
}

func NewApplication() *Application {
	var this *Application = new(Application)
	this.QuitSignal = false
	this.Files = make([]*File, 0)
	this.StatusLine = ""
	return this
}

