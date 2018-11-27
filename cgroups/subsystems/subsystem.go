package subsystems

//用于传递资源限制配置的结构体, 包含内存限制, CPU时间片权重, CPU核心数
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

//Subsystem接口, 每个Subsystem可以实现下面的4个接口
//这里将cgroup抽象成了path, 原因是cgroup在hierarchy的路径, 便是虚拟文件系统中的虚拟路径
type Subsystem interface {
	Name() string
	Set(path string, res *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

//通过不同的subsystem初始化实例创建资源处理链数组
var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)