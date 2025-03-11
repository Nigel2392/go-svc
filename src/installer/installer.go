package installer

type Installer interface {
	Install() error
	Installed() (bool, error)
	Remove() error
	Close() error
}
