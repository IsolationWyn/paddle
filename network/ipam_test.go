package network

import(
	"testing"
	"net"
)

func TestAllocate(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.0.0/24")
	ip, _ := ipAllocator.Allocate(ipnet)
	t.Logf("alloc ip: %v\n", ip)
}

func TestRelease(t *testing.T) {
	// 在192.168.0.0/24网段中释放刚才分配的192.168.0.1的IP
	ip, ipnet, _ := net.ParseCIDR("192.168.0.0/24")
	ipAllocator.Release(ipnet, &ip)
}