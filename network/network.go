package network

import (
	"path/filepath"
	"encoding/json"
	"fmt"
	"path"
	"os"
	"text/tabwriter"
	"github.com/sirupsen/logrus"
	"net"
	"github.com/vishvananda/netlink"
)

var (
	defaultNetworkPath = "/var/run/paddle/network/network"
	drivers = map[string]NetworkDriver{}
	networks = map[string]*Network{}
)

type Endpoint struct {
	ID	string `json:"id"`
	Device	netlink.Veth `json:"dev"`
	IPAddress	net.IP `json:"ip"`
	MacAddress	net.HardwareAddr `json:"mac"`
	PortMapping []string `json:"portmapping"`
	Network	*Network
}

type Network struct {
	Name string
	IpRange *net.IPNet
	Driver string
}

type NetworkDriver interface {
	Name() string
	Create(subnet string, name string) (*Network, error)
	Delete(network Network) error
	Connect(network *Network, endpoint *Endpoint) error
	Disconnect(network Network, endpoint *Endpoint) error
}

func CreateNetwork(driver, subnet, name string) error {
	// 将网段字符串转换成net.IPNet的对象
	// type IPNet struct {
	// 	IP   IP     // network number
	// 	Mask IPMask // network mask
	// }
	_, cidr, _ := net.ParseCIDR(subnet)
	// 通过IPAM分配网关ip, 获取到网段中第一个IP作为网关的IP
	gatewayIp, err := ipAllocator.Allocate(cidr)
	if err != nil {
		return err
	}
	cidr.IP = gatewayIp
	/* 调用指定的网络驱动创建网络, 这里的drivers字典是各个网络驱动的实例字典, 通过调用网络驱动的Create方法创建网络 */
	nw, err := drivers[driver].Create(cidr.String(), name)
	if err != nil {
		return err
	}
	// 保存网络信息, 将网络的信息保存在文件系统中, 以便查询和网络上连接网络端点
	return nw.dump(defaultNetworkPath)
}

func (nw *Network) dump(dumpPath string) error {
	// 检查保存的目录是否存在, 不存在则创建
	if _, err := os.Stat(dumpPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dumpPath, 0644)
		} else {
			return err
		}
	}
	// 保存的文件名是网络的名字
	nwPath := path.Join(dumpPath, nw.Name)
	// 打开保存的文件用于写入
	// os.O_RDONLY：只读
	// os.O_WRONLY：只写
	// os.O_CREATE：创建：如果指定文件不存在，就创建该文件。
	// os.O_TRUNC：截断：如果指定文件已存在，就将该文件的长度截为0
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		logrus.Errorf("error: ", err)
		return err
	}
	defer nwFile.Close()

	// 通过json的库序列化网络对象到json的字符串
	nwJson, err := json.Marshal(nw)
	if err != nil {
		logrus.Errorf("error:", err)
		return err
	}
	// 将网络配置的json字符串写入到文件中
	_, err = nwFile.Write(nwJson)
	if err != nil {
		logrus.Errorf("error:", err)
		return err
	}
	return nil
}

func (nw *Network) load(dumpPath string) error {
	// 打开配置文件
	nwConfigFile, err := os.Open(dumpPath)
	defer nwConfigFile.Close()
	if err != nil {
		return err
	}
	// 从配置文件中读取网络的配置json字符串
	nwJson := make([]byte, 2000)
	n, err := nwConfigFile.Read(nwJson)
	if err != nil {
		return err
	}
	// 通过json字符串反序列化出网络
	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		logrus.Errorf("Error load nw info", err)
		return err
	}
	return nil
}

func connect(networkName string, cinfo *container.ContainerInfo) error {
	// 从networks字典中取到容器连接的网络的信息, networks字典中保存了当前以及创建的网络
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("No Such Network: %s", networkName)
	}
	// 通过调用IPAM从网络的网段中获取可用的IP作为容器IP地址
	ip, err := ipAllocator.Allocate(network.IpRange)
	if err != nil {
		return err
	}
	// 创建网络端点
	ep := &Endpoint{
		ID: fmt.Sprintf("%s-%s", cinfo.Id, networkName),
		IPAddress: ip,
		Network: network,
		PortMapping: cinfo.PortMapping,
	}
	// 调用网络驱动的"Connect"方法去连接和配置网络端点
	if err = drivers[network.Driver].Connect[network.ep]; err != nil {
		return err
	}
	// 进入到容器的网络Namespace配置容器网络设备的ip地址和路由
	if err = configEndpointIpAddressAndRoute(ep, cinfo); err != nil {
		return err
	}
	// 配置容器到宿主机的端口映射
	return configPortMapping(ep, cinfo)
}

func Init() error {
	// 加载网络驱动
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver
	// 判断网络的配置目录是否存在, 不存在则创建
	if _, err := os.Stat(defaultNetworkPath); err != nil {
		if os.IsNotExist(err) {
		os.MkdirAll(defaultNetworkPath, 0644)
		} else {
			return err
		}
	}
	// 检查网络配置目录中所有文件
	filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		// 加载文件名作为网络名
		_, nwName := path.Split(nwPath)
		nw := &Network{
			Name: nwName,
		}
		// 加载网络配置信息
		if err := nw.load(nwPath); err != nil {
			logrus.Errorf("error load network: %s", err)
		}

		// 将网络的配置信息加入到networks字典中
		networks[nwName] = nw
		return nil
	})
	return nil
}

func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	// 遍历网络信息
	for _, nw := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", 
			nw.Name,
			nw.IpRange.String(),
			nw.Driver,)
	}
	// 输出到标准输入输出
	if err := w.Flush(); err != nil {
		logrus.Errorf("Flush error %v", err)
		return
	}
}

/*
删除网络
1.删除网络网关ip
2.删除网络对应网络设备
3.删除网络配置文件
*/

func DeleteNetwork(networkName string) error {
	// 查找网络是否存在
	nw, ok := Network[networkName]
	if !ok {
		return fmt.Errorf("No such Network: %s", networkName)
	}
	//调用IPAM实例ipAllocator释放网络网关的ip
	if err := ipAllocator.Release(nw.IpRange, &nw.IpRange.IP); err != nil {
		return fmt.Errorf("Error Remove Network gateway ip: %s", err)
	}

	/*调用网络驱动删除网络创建的设备与配置*/
	if err := drivers[nw.Driver].Delete(*nw); err != nil {
		return fmt.Errorf("Error Remove Network DriverError: %s", err)
	}

	// 从网络的配置目录中删除该网络对应的配置文件
	return nw.remove(defaultNetworkPath)
}

// Network.remove会从网络配置目录中删除网络的配置文件
func (nw *Network) remove(dumpPath string) error {
	// 网络对应的配置文件, 即配置文件下的网络名文件
	if _, err := os.Stat(path.Join(dumpPath, nw.Name)); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
				return err
			} 
		} else {
			// 调用 os.Remove 删除这个网络对应的配置文件
			return os.Remove(path.Join(dumpPath, nw.Name))
	}
}

