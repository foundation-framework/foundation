package units

type Byte uint64

const (
	Kilobyte = 1024
	Megabyte = Kilobyte * 1024
	Gigabyte = Megabyte * 1024
	Terabyte = Gigabyte * 1024
	Petabyte = Terabyte * 1024
	Exabyte  = Petabyte * 1024
)
