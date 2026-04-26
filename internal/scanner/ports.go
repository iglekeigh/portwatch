package scanner

import (
	"fmt"
	"strconv"
	"strings"
)

// CommonPorts is a curated list of well-known ports to scan by default.
var CommonPorts = []int{
	21, 22, 23, 25, 53, 80, 110, 143, 443, 465,
	587, 993, 995, 3306, 3389, 5432, 5900, 6379,
	8080, 8443, 8888, 27017,
}

// ParsePortRange parses a port expression such as "80", "80,443", or "8000-8080".
func ParsePortRange(expr string) ([]int, error) {
	var ports []int
	for _, part := range strings.Split(expr, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			start, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("invalid port range start %q: %w", bounds[0], err)
			}
			end, err := strconv.Atoi(bounds[1])
			if err != nil {
				return nil, fmt.Errorf("invalid port range end %q: %w", bounds[1], err)
			}
			if start > end {
				return nil, fmt.Errorf("port range start %d > end %d", start, end)
			}
			for p := start; p <= end; p++ {
				ports = append(ports, p)
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port %q: %w", part, err)
			}
			ports = append(ports, p)
		}
	}
	return ports, nil
}

// ValidatePort checks whether p is a valid TCP/UDP port number (1–65535).
func ValidatePort(p int) error {
	if p < 1 || p > 65535 {
		return fmt.Errorf("port %d out of valid range (1-65535)", p)
	}
	return nil
}
