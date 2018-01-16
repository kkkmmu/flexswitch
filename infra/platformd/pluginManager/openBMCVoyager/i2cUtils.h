#ifndef _I2C_UTILS_H
#define _I2C_UTILS_H

int i2cSet(int i2cBusNum, int chipAddr, int dataAddr, int val);
int i2cGet(int i2cBusNum, int chipAddr, int dataAddr);
#endif
