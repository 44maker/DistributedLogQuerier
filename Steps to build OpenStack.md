# 本地搭建Openstack步骤

网络1 仅主机 192.168.10.0

网络2 仅主机 192.168.20.0



```bash
root
59561124
vi /etc/sysconfig/network-scripts/ifcfg-ens33
BOOTPROTO=static
ONBOTT=yes
IPADDR=192.168.10.10
NETMASK=255.255.255.0

systemctl restart network

vi /etc/sysconfig/network-scripts/ifcfg-ens34
IPADDR=192.168.20.10
NETMASK=255.255.255.0

systemctl stop firewalld
systemctl disable firewalld

setenforce 0
vi /etc/selinux/config
SELINUX=disable
vi /etc/hosts
192.168.10.10 controller
192.168.10.20 compute
echo '192.168.10.30 compute02
192.168.10.40 compute03
192.168.10.50 compute04' >> /etc/hosts

rm -rf /etc/yum.repos.d/*
```

关机、克隆

```bash
vi /etc/sysconfig/network-scripts/ifcfg-ens33
BOOTPROTO=static
ONBOTT=yes
IPADDR=192.168.10.20
NETMASK=255.255.255.0

systemctl restart network

vi /etc/sysconfig/network-scripts/ifcfg-ens34
IPADDR=192.168.20.20
NETMASK=255.255.255.0

vi /etc/sysconfig/network-scripts/ifcfg-ens33
IPADDR=192.168.10.30

systemctl restart network

vi /etc/sysconfig/network-scripts/ifcfg-ens34
IPADDR=192.168.20.30
NETMASK=255.255.255.0
```



```bash
#controller中
hostnamectl set-hostname controller
#compute
hostnamectl set-hostname compute
上传呢两个镜像
#controller
yum 源配置
vi /etc/yum.repos.d/local.repo
[centos]
name=centos
enabled=1
gpgcheck=0
baseurl=file:///opt/centos
[iaas]
name=iaas
enabled=1
gpgcheck=0
baseurl=file:///opt/iaas/iaas-repo

mkdir /opt/centos
mkdir /opt/iaas
```

```bash
#compute中
yum 源配置
rm -rf /etc/yum.repos.d/*
vi /etc/yum.repos.d/ftp.repo
[centos]
name=centos
enabled=1
gpgcheck=0
baseurl=ftp://192.168.10.10/centos
[iaas]
name=iaas
enabled=1
gpgcheck=0
baseurl=ftp://192.168.10.10/iaas/iaas-repo

```

```bash
#controller中
mount -o loop CentOS-7-x86_64-DVD-2009.iso /mnt/
cp -rvf /mnt/* /opt/centos
umount /mnt/

mount -o loop openstack.iso /mnt/
cp -rvf /mnt/* /opt/iaas/
umount /mnt/

yum repolist
rm -rf /etc/yum.repos.d/C*
yum install vsftpd -y
vi /etc/vsftpd/vsftpd.conf
anon_root=/opt

systemctl restart vsftpd
systemctl enable vsftpd
```



```bash
#controller&compute
yum install openstack-iaas -y
#controller
vi /etc/openstack/openrc.sh
ctrl v G d
HOST_IP=192.168.10.10
:%s/PASS=/PASS=000000/g


scp /etc/openstack/openrc.sh root@192.168.10.20:/etc/openstack/
```



```bash
#controller&compute
iaas-pre-host.sh
重新连接ssh
#controller
iaas-install-mysql.sh
iaas-install-keystone.sh
iaas-install-glance.sh
iaas-install-placement.sh
iaas-install-nova-controller.sh
iaas-install-neutron-controller.sh
iaas-install-dashboard.sh
#compute
iaas-install-nova-compute.sh
iaas-install-neutron-compute.sh
```



# 分布式日志查询

在所有机器上安装go

```bash
tar -zxf go1.17.2.linux-amd64.tar.gz -C /usr/local
vi /etc/profile
#go 环境变量
export GO111MODULE=on
export GOROOT=/usr/local/go
export GOPATH=/home/gopath
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

source /etc/profile
```



```bash
#在所有compute节点中
go run server.go
#在controller中
go build client.go
./client [query] [log file name]
If running for demo purpose, [log_file_name] should be vm.log

To Run Unit Test:
Generate all unit tests:

go run generate_testfiles.go
#在所有compute节点中
go run server.go
#在controller中
go test -v client_test.go
```







