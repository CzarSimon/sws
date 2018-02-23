package proxy

import "fmt"

// Pair holds info about primary and secondary proxy.
type Pair struct {
	Primary   Proxy
	Secondary Proxy
}

// Alinged checks if primary and secondary proxy are identical in all but name.
func (pair Pair) Alinged() error {
	if pair.Primary.Port != pair.Secondary.Port {
		return fmt.Errorf("primary.Port=%d secondary.Port=%d",
			pair.Primary.Port, pair.Secondary.Port)
	}
	if pair.Primary.ID() != pair.Secondary.ID() {
		return fmt.Errorf("primary.ID=%s secondary.ID=%s",
			pair.Primary.ID(), pair.Secondary.ID())
	}
	return nil
}
