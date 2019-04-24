package network

import (
	"fmt"
	"net"
	"strings"
	"time"
	"github.com/vishvananda/netlink"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

type BridgeNetworkDriver struct {
}

func (d *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (d *BridgeNetworkDriver) Create(subnet string, name string) (*Network, error) {
	// 通过net包中的net.ParseCIDR方法, 取到网段的字符串中的网关IP地址和网络IP段
	ip, ipRange, _ := net.ParseCIDR(subnet)
	ipRange.IP = ip
	// 初始化网络对象
	n := &Network {
		Name: name,
		IpRange: ipRange,
		Driver: d.Name(),
	}
	// 配置Linux Bridge 
	err := d.initBridge(n)
	if err != nil {
		log.Errorf("error init bridge: %v", err)
	}
	// 返回配置好的网络
	return n, err
}

func (d *BridgeNetworkDriver) Delete(network Network) error {
	// 网络名即Linux Bridge的设备名
	bridgeName := network.Name
	// 通过netlink库的LinkByName找到网络对应的设备
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}
	// 删除网络对应的Linux Bridge设备
	return netlink.LinkDel(br)
}


// 连接一个网络和网络端点
func (d *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {
	// 获取网络名, 即Linux Bridge的接口名
	bridgeName := network.Name
	// 通过接口名获取到Linux Bridge接口的对象和接口属性
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}

	// 创建Vet接口的配置
	la := netlink.NewLinkAttrs()
	// 由于Linux接口名的限制, 名字取endpoint ID的前五位
	la.Name = endpoint.ID[:5]
	// 通过设置Veth接口的master属性, 设置这个Veth的一端挂载到网络对应的Linux Bridge上
	la.MasterIndex = br.Attrs().Index

	// 创建Veth对象, 通过PeerName配置Veth另一端的接口名
	// 配置Veth另一端的名字cif-(endpoint ID的前五位)
	endpoint.Device = netlink.Veth{
		LinkAttrs: la,
		PeerName:  "cif-" + endpoint.ID[:5],
	}

	// 调用netlink的LinkAdd方法创建出这个Veth接口
	// 因为上面指定了link的MasterIndex是网络对应的Linux Bridge
	// 所以Veth的一端就已经挂载到了网络对应的Linux Bridge上
	if err = netlink.LinkAdd(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Add Endpoint Device: %v", err)
	}

	// 调用netlink的LinkSetUp方法, 设置Veth启动
	// 相当于设置ip link set xxx up 命令
	if err = netlink.LinkSetUp(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Add Endpoint Device: %v", err)
	}
	return nil
}

func (d *BridgeNetworkDriver) Disconnect(network Network, endpoint *Endpoint) error {
	return nil
}



func (d *BridgeNetworkDriver) initBridge(n *Network) error {
	// 1. 创建Bridge网络设备
	bridgeName := n.Name
	if err := createBridgeInterface(bridgeName); err != nil {
		return fmt.Errorf("Error add bridge： %s, Error: %v", bridgeName, err)
	}

	// 2. 设置Bridge设备的地址和路由
	gatewayIP := *n.IpRange
	gatewayIP.IP = n.IpRange.IP

	if err := setInterfaceIP(bridgeName, gatewayIP.String()); err != nil {
		return fmt.Errorf("Error assigning address: %s on bridge: %s with an error of: %v", gatewayIP, bridgeName, err)
	}

	// 3. 启动Bridge设备
	if err := setInterfaceUP(bridgeName); err != nil {
		return fmt.Errorf("Error set bridge up: %s, Error: %v", bridgeName, err)
	}

	// 4. 设置iptables的SNAT规则
	if err := setupIPTables(bridgeName, n.IpRange); err != nil {
		return fmt.Errorf("Error setting iptables for %s: %v", bridgeName, err)
	}

	return nil
}

// deleteBridge deletes the bridge
func (d *BridgeNetworkDriver) deleteBridge(n *Network) error {
	bridgeName := n.Name

	// get the link
	l, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("Getting link with name %s failed: %v", bridgeName, err)
	}

	// delete the link
	if err := netlink.LinkDel(l); err != nil {
		return fmt.Errorf("Failed to remove bridge interface %s delete: %v", bridgeName, err)
	}

	return nil
}

// 创建Bridge网络设备
func createBridgeInterface(bridgeName string) error {
	// 先检查是否已经存在这个同名的Bridge设备
	_, err := net.InterfaceByName(bridgeName)
	// 如果已经存在或者报错则返回创建错误
	if err == nil || !strings.Contains(err.Error(), "no such network interface") {
		return err
	}

	// 初始化一个netlink的Link基础对象, Link的名字即Bridge虚拟设备的名字
	la := netlink.NewLinkAttrs()
	la.Name = bridgeName

	// 使用刚才创建的Link属性创建netlink的Bridge对象
	br := &netlink.Bridge{LinkAttrs: la}
	// 调用netlink的LinkAdd方法, 创建Bridge网络设备
	if err := netlink.LinkAdd(br); err != nil {
		return fmt.Errorf("Bridge creation failed for bridge %s: %v", bridgeName, err)
	}
	return nil
}

// 设置网络接口为UP状态
func setInterfaceUP(interfaceName string) error {
	// 找到需要设置的网络接口
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return fmt.Errorf("Error retrieving a link named [ %s ]: %v", iface.Attrs().Name, err)
	}
	// 通过"netlink"的"LinkSetUp"方法设置接口状态为"UP"状态
	// 等价于ip link set xxx up命令
	if err := netlink.LinkSetUp(iface); err != nil {
		return fmt.Errorf("Error enabling interface for %s: %v", interfaceName, err)
	}
	return nil
}


// 设置一个网络接口的IP地址, 例如setInterfaceIP("testbridge", "192.168.0.1./24")
func setInterfaceIP(name string, rawIP string) error {
	retries := 2
	var iface netlink.Link
	var err error
	// 通过netlink的LinkByName方法找到需要设置的网络接口
	for i := 0; i < retries; i++ {
		iface, err = netlink.LinkByName(name)
		if err == nil {
			break
		}
		log.Debugf("error retrieving new bridge netlink link [ %s ]... retrying", name)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("Abandoning retrieving the new bridge link from netlink, Run [ ip link ] to troubleshoot the error: %v", err)
	}
	// 由于netlink.ParseIPNet是对net.ParseCIDR的一个封装, 因此可以将net.ParseCIDR的返回值中的IP和net整合
	// 返回值中的ipNet既包含了网段的信息, 192.168.0.0/24, 也包含了原始的ip 192.168.0.1
	ipNet, err := netlink.ParseIPNet(rawIP)
	if err != nil {
		return err
	}
	// 通过netlink.AddrAdd给网路接口配置地址, 相当于ip addr add xxx 命令
	// 同时如果配置了地址所在的网段信息, 例如192.168.0.0/24
	// 还会配置路由表192.168.0.0/24转发到这个testbridge的网路接口上
	addr := &netlink.Addr{IPNet: ipNet, Peer: ipNet, Label: "", Flags: 0, Scope: 0}
	return netlink.AddrAdd(iface, addr)
}

// 设置iptables对应bridge的MASQUERADE规则
// 通过直接执行iptables命令, 创建SNAT规则, 只要从这个网桥上出来的包, 都会对其进行源IP的转换, 保证了
// 容器经过宿主机访问到达宿主机外部网络请求的包转换成机器的IP
func setupIPTables(bridgeName string, subnet *net.IPNet) error {
	iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	// 执行iptables命令配置SNAT规则
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables Output, %v", output)
	}
	return err
}
