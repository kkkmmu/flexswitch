SnapRoute ASIC Daemon 
========================

More documentation is available at 
[Product Overview] (http://opensnaproute.github.io/docs/)

Complete system architecture can be found 
[Here](http://opensnaproute.github.io/docs/architecture.html) 


![alt text](https://github.com/OpenSnaproute/asicd/blob/master/Asic_Daemon.jpg "High level architecture diagram")

##Overview
Snaproute's asic daemon serves as a hardware abstraction layer (HAL). A common northbound API interface is exposed to all protocol daemons. This interface allows provisioning a range of packet processing ASICs. Support for software simulation on a linux host OS is also provided.

##Software Architecture
####Northbound interface :
The asic daemons northbound interface is implemented using thrift RPC. This is the interface that is used by users/protocols to apply configuration.

####Core resource managers :
The core infrastructure within Asicd is distributed across multiple resource managers, e.x. portMgr.go, routeMgr.go, neighborMgr.go etc. Each of these individual resource managers support Create/Retrieve/Update and Delete operations on the corresponding resource. These resource managers also maintain any relevant state data for each corresponding resource.

####Plugins :
The asic daemon uses a plugin based approach to effectively abstract differences between ASICs from multiple vendors. The following plugins and asic vendors are currently supported
- Vendor SDK 
- SAI 
- Softswitch (Linux host)

####Events handling :
The asicd daemon supports signaling/notification of asynchronous events. The notification engine employs a nano message based publisher. Notifications for the following events are supported
- Port operational state changes
- Vlan/Lag interface creation/deletion
- IP interface operational state changes

Extending ASICd to support a new ASIC
=====================================

Snaproute's ASIC daemon, currently has support to provision multiple vendor Asic's. ASICd can easily be ported over to a new vendor Asic as documented below.

The following steps detail how ASICd can be ported over to support a new silicon vendor's chip.

#####Step 1:
Provide implementations for all methods defined in ASICd's thrift interface. The handler functions for the thrift interface are located in the following files
- [rpc/l2Services.go](https://github.com/OpenSnaproute/asicd/blob/master/rpc/l2Services.go) 
- [rpc/l3Services.go](https://github.com/OpenSnaproute/asicd/blob/master/rpc/l3Services.go) 
- [rpc/vlanServices.go](https://github.com/OpenSnaproute/asicd/blob/master/rpc/vlanServices.go) 
- [rpc/portServices.go](https://github.com/OpenSnaproute/asicd/blob/master/rpc/portServices.go) 

#####Step 2:
Compile ASICd by running 'make BUILD_TARGET=custom'

#####Step 3:
Build a flexswitch package, install and run.
