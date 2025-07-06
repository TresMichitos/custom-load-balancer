// IP Hashing load balancing algorithm

package lbalgorithms

// Struct to implement serverpool.LbAlgorithm interface
type IpHashing struct {
	index int
}
