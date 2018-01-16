#include <sys/ioctl.h>
#include <errno.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <linux/i2c.h>
#include <linux/i2c-dev.h>
#include <linux/types.h>
#include <sys/stat.h>
#include <fcntl.h>

#define MISSING_FUNC_FMT "Error: Adapter does not have %s capability\n"

__s32 i2c_smbus_access(int file, char read_write, __u8 command,
                       int size, union i2c_smbus_data *data)
{
        struct i2c_smbus_ioctl_data args;
        __s32 err;

        args.read_write = read_write;
        args.command = command;
        args.size = size;
        args.data = data;

        err = ioctl(file, I2C_SMBUS, &args);
        if (err == -1) 
                err = -errno;
        return err;
}



__s32 i2c_smbus_read_byte_data(int file, __u8 command)
{
        union i2c_smbus_data data;
        int err;

        err = i2c_smbus_access(file, I2C_SMBUS_READ, command,
                               I2C_SMBUS_BYTE_DATA, &data);
        if (err < 0)
                return err;

        return 0x0FF & data.byte;
}

__s32 i2c_smbus_write_byte_data(int file, __u8 command, __u8 value)
{
        //printf("Addr: %d Value %d\n", command, value);
        union i2c_smbus_data data;
        data.byte = value;
        return i2c_smbus_access(file, I2C_SMBUS_WRITE, command,
                                I2C_SMBUS_BYTE_DATA, &data);
}


/*
 *  * Parse an I2CBUS command line argument and return the corresponding
 *   * bus number, or a negative value if the bus is invalid.
 *    */
int lookup_i2c_bus(int i2cbus)
{
        if (i2cbus > 0xFFFFF) {
                fprintf(stderr, "Error: I2C bus out of range!\n");
                return -2; 
        }

        return i2cbus;
}

/*
 *  * Parse a CHIP-ADDRESS command line argument and return the corresponding
 *   * chip address, or a negative value if the address is invalid.
 *    */
int parse_i2c_address(int address)
{
        if (address < 0x03 || address > 0x77) {
                fprintf(stderr, "Error: Chip address out of range "
                        "(0x03-0x77)!\n");
                return -2;
        }

        return address;
}

int open_i2c_dev(int i2cbus, char *filename, size_t size)
{
        int file;

        snprintf(filename, size, "/dev/i2c/%d", i2cbus);
        filename[size - 1] = '\0';
        file = open(filename, O_RDWR);

        if (file < 0 && (errno == ENOENT || errno == ENOTDIR)) {
                sprintf(filename, "/dev/i2c-%d", i2cbus);
                file = open(filename, O_RDWR);
        }

        if (file < 0) {
                if (errno == ENOENT) {
                        fprintf(stderr, "Error: Could not open file "
                                "`/dev/i2c-%d' or `/dev/i2c/%d': %s\n",
                                i2cbus, i2cbus, strerror(ENOENT));
                } else {
                        fprintf(stderr, "Error: Could not open file "
                                "`%s': %s\n", filename, strerror(errno));
                        if (errno == EACCES)
                                fprintf(stderr, "Run as root?\n");
                }
        }

        return file;
}

int set_slave_addr(int file, int address)
{
        /* With force, let the user read from/write to the registers
 *            even when a driver is also running */
        if (ioctl(file, I2C_SLAVE, address) < 0) {
                fprintf(stderr,
                        "Error: Could not set address to 0x%02x: %s\n",
                        address, strerror(errno));
                return -errno;
        }

        return 0;
}

static int check_funcs(int file)
{
	unsigned long funcs;

	/* check adapter functionality */
	if (ioctl(file, I2C_FUNCS, &funcs) < 0) {
		fprintf(stderr, "Error: Could not get the adapter "
			"functionality matrix: %s\n", strerror(errno));
		return -1;
	}

	if (!(funcs & I2C_FUNC_SMBUS_WRITE_BYTE_DATA)) {
		fprintf(stderr, MISSING_FUNC_FMT, "SMBus write byte");
		return -1;
	}

	return 0;
}

int i2cSet(int i2cBusNum, int chipAddr, int dataAddr, int val)
{
	int res, i2cbus, address, file;
	int value, daddress;
	char filename[20];

	i2cbus = lookup_i2c_bus(i2cBusNum);
	if (i2cbus < 0)
		return -1;

	address = parse_i2c_address(chipAddr);
	if (address < 0)
		return -1;

	daddress = dataAddr;
	if (daddress < 0 || daddress > 0xff) {
		fprintf(stderr, "Error: Data address invalid!\n");
		return -1;
	}

	/* read values from command line */
	value = val;
	if (value > 0xff) {
		fprintf(stderr, "Error: Data value out of range!\n");
		return -1;
	}

	//printf("I2CBus = %d Address = %d\n", i2cbus, address);
	file = open_i2c_dev(i2cbus, filename, sizeof(filename));
	if (file < 0
	 || check_funcs(file)
	 || set_slave_addr(file, address))
		return -1;

	res = i2c_smbus_write_byte_data(file, daddress, value);
	if (res < 0) {
		//fprintf(stderr, "Error: Write failed\n");
		close(file);
		return -1;
	}

	close(file);

	return 0;
}

int i2cGet(int i2cBusNum, int chipAddr, int dataAddr)
{
        int res, i2cbus, address, file;
        int daddress;
        char filename[20];

        i2cbus = lookup_i2c_bus(i2cBusNum);
        if (i2cbus < 0)
		return -1;

        address = parse_i2c_address(chipAddr);
        if (address < 0)
		return -1;

	daddress = dataAddr; 
	if (daddress < 0 || daddress > 0xff) {
		fprintf(stderr, "Error: Data address invalid!\n");
		return -1;
	}

        file = open_i2c_dev(i2cbus, filename, sizeof(filename));
        if (file < 0
         || check_funcs(file)
         || set_slave_addr(file, address))
		return -1;

        res = i2c_smbus_read_byte_data(file, daddress);
        close(file);

        if (res < 0) {
                fprintf(stderr, "Error: Read failed\n");
                return -1;
        }

        //printf("0x%0*x\n", 2, res);

        return res;
}

