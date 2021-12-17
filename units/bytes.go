package units

import "fmt"

const (
	Kilobyte int64 = 1024
	Megabyte       = Kilobyte * 1024
	Gigabyte       = Megabyte * 1024
	Terabyte       = Gigabyte * 1024
	Petabyte       = Terabyte * 1024
	Exabyte        = Petabyte * 1024
)

func FormatBytes(n int64) string {
	if n < Megabyte {
		return fmt.Sprintf("%.3f KiB", float64(n)/float64(Kilobyte))
	}

	if n < Gigabyte {
		return fmt.Sprintf("%.2f MiB", float64(n)/float64(Megabyte))
	}

	if n < Terabyte {
		return fmt.Sprintf("%.2f GiB", float64(n)/float64(Gigabyte))
	}

	if n < Petabyte {
		return fmt.Sprintf("%.2f TiB", float64(n)/float64(Terabyte))
	}

	if n < Exabyte {
		return fmt.Sprintf("%.2f PiB", float64(n)/float64(Petabyte))
	}

	return fmt.Sprintf("%.2f PiB", float64(n)/float64(Exabyte))
}
