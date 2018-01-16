namespace go sysdInt
typedef i32 int
service SYSDINTServices {
	oneway void PeriodicKeepAlive(1:string Name);
}
