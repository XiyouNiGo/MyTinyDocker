# MyTinyDocker
## 项目简述
+ MyTinyDocker是采用Go语言编写，实现run、exec、log、volume、network等命令的小型容器。
+ 项目采用Linux Namespace、PrivotRoot实现了容器之间资源的隔离；采用Cgroups实现了对单个容器的资源限制、控制；采用AUFS为每个容器构造独立的roofts，利用写时复制实现高效容器存储管理；采用Cgo配合setns和环境变量实现exec命令；采用Linux Veth、Bridge和iptables路由，配合IPAM构建容器网络模型。
## 技术要点
### Linux Namespace
+ 容器的基本是隔离性，而Linux namespaces是对全局系统资源的一种封装隔离，使得处于不同namespace的进程拥有独立的全局系统资源，改变一个namespace中的系统资源只会影响当前 namespace 里的进程，对其他 namespace中的进程没有影响。
### Cgroups
+ Cgroups是Linux内核提供的一种可以限制、记录、隔离进程组所使用的物理资源的机制。可以对一组进程及将来的子进程的资源的限制、控制和统计的能力，这些资源包括CPU，内存，存储，网络等。
+ Cgroups主要由三个组件构成。其中，cgroup是对进程分组管理的一种机制，subsystem是一组资源控制的模块，hierarchy的功能是把一组cgroup串成一个树形的结构。
+ 系统默认创建的hierarchy可以在/proc/self/meminfo中获得。
### AUFS
+ AUFS是一种Union FS，是文件级的存储驱动。AUFS能透明覆盖一或多个现有文件系统的层状文件系统，把多层合并成文件系统的单层表示。简单来说就是支持将不同目录挂载到同一个虚拟文件系统下的文件系统。这种文件系统可以一层一层地叠加修改文件。无论底下有多少层都是只读的，只有最上层的文件系统是可写的。
+ 当需要修改一个文件时，AUFS创建该文件的一个副本，使用CoW将文件从只读层复制到可写层进行修改，结果也保存在可写层。在Docker中，底下的只读层就是image，可写层就是Container。
## run命令实现
+ run命令运行后，首先调用NewParentProcess函数创建对应exec.Cmd对象，具体调用/proc/self/exe，相当于调用自己，并传入init参数，使子进程能够进入init函数来初始化容器内部资源，例如挂载proc文件系统等等。
+ 为了使容器中的PID为1的进程是容器进程，由于PID为1的进程是不能被kill掉的，我们需要使用syscall.Exec来覆盖当前进程的镜像、数据和堆栈等信息。
+ 为了解决传参过程中的缺陷，使用无名管道在父子进程间传参。
## 构建镜像实现
+ 创建容器文件系统时，首先创建只读层、容器读写层，创建挂载点并并把只读层和读写层挂载到挂载点上，最后如果存在用户指定的volume，执行MountVolume来挂载数据卷。而卸载时，需要先卸载volume挂载点的文件系统，再卸载整个容器文件系统的挂载点，否则会卸载失败。
+ 构建镜像，为了改变当前的root文件系统，由于chroot只针对某个进程，我们需要调用privot_root来改变整个系统的root文件系统。
### exec命令实现
+ 要使容器进入某个Namespace，调用setns系统调用，传入对应的PID即可。但Go每启动一个程序就会进入多进程状态，因此需要借助Cgo调用C代码，并通过环境变量指定PID和命令，以及进行区分容器进程。
+ 对于指定了环境变量的容器，执行exec命令的进程跟执行run命令的进程不是一个进程，因此无法同步环境变量，对于新进程，我们需要到/proc/\<PID\>/environ获取，再加入exec.Cmd对象。
### network命令实现
+ 构建容器网络主要使用了Veth和Linux Bridge。Veth是成对出现的虚拟网络设备，发送到Veth一端虚拟设备的请求会从另一端的虚拟设备中发出。Bridge虚拟设备是用来桥接的网络设备，相当于交换机，当请求到达Bridge设备时，会通过报文中的Mac地址进行广播或转发。
+ iptables 是对 Linux 内核的 netfilter 模块进行操作和展示的工具，用来管理包的流动和转送。 iptables 定义了一套链式处理的结构，在网络包传输的各个阶段可以使用不同的策略对包进行加工、传送或丢弃。在容器虚拟化的技术中，经常会用到两种策略 MASQUERADE和DNAT ，用于容器和宿主机外部的网络通信。
+ 为保证网络中IP地址唯一，我们借助bitmap算法对网关IP和容器的网络端点IP进行分配，它对大规模连续且少状态的数据处理有很高的效率。
## 相关项目
[https://github.com/xianlubird/mydocker](https://github.com/xianlubird/mydocker)
