#include <stdio.h>
#include "onlp.h"
#include <onlp/sys.h>
#include <onlp/fan.h>
#include <onlp/sfp.h>
#include <onlp/thermal.h>
#include <onlp/psu.h>
#include <onlp/led.h>
#include <dlfcn.h>

/* ONLP init */
int (*onlpInit)(void);
static void *handle;

/* System EEPROM */
int (*onlpSysInfoGet)(onlp_sys_info_t*);

/* FAN */
int (*onlpFanInfoGet)(onlp_oid_t, onlp_fan_info_t*);
static onlp_oid_t fan_oid_table[ONLP_OID_TABLE_SIZE];
static onlp_fan_info_t fan_info_table[ONLP_OID_TABLE_SIZE];
int flag[ONLP_OID_TABLE_SIZE] = {0};

/* SFP */
void (*onlpSfpBitmapInit)(onlp_sfp_bitmap_t*);
int (*onlpSfpBitmapGet)(onlp_sfp_bitmap_t*);
int (*onlpSfpIsPresent)(int);
int (*onlpSfpPresenceBitmapGet)(onlp_sfp_bitmap_t*);
int (*onlpSfpRxLosBitmapGet)(onlp_sfp_bitmap_t*);
int (*onlpSfpControlGet)(int, onlp_sfp_control_t, int*);
int (*onlpSfpEepromRead)(int, uint8_t**);

/* Thermal */
int (*onlpThermalInfoGet)(onlp_oid_t, onlp_thermal_info_t*);

/* PSU */
int (*onlpPsuInfoGet)(onlp_oid_t, onlp_psu_info_t*);

/* LED */
int (*onlpLedInfoGet)(onlp_oid_t, onlp_led_info_t*);

int
loadOnlpSym()
{
    char *error;

    handle = dlopen ("/lib/x86_64-linux-gnu/libonlp.so", RTLD_LAZY);
    if (!handle) {
        printf(dlerror(), stderr);
        return -1;
    }

	onlpInit = dlsym(handle, "onlp_init");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
        return -1;
    }

	onlpFanInfoGet = dlsym(handle, "onlp_fan_info_get");
	if ((error = dlerror()) != NULL)  {
		printf(error, stderr);
		return -1;
	}

    onlpSysInfoGet = dlsym(handle, "onlp_sys_info_get");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpSfpBitmapInit = dlsym(handle, "onlp_sfp_bitmap_t_init");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpSfpBitmapGet = dlsym(handle, "onlp_sfp_bitmap_get");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpSfpIsPresent = dlsym(handle, "onlp_sfp_is_present");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpSfpPresenceBitmapGet = dlsym(handle, "onlp_sfp_presence_bitmap_get");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpSfpRxLosBitmapGet = dlsym(handle, "onlp_sfp_rx_los_bitmap_get");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpSfpEepromRead = dlsym(handle, "onlp_sfp_eeprom_read");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpSfpControlGet = dlsym(handle, "onlp_sfp_control_get");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpThermalInfoGet = dlsym(handle, "onlp_thermal_info_get");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpPsuInfoGet = dlsym(handle, "onlp_psui_info_get");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }

    onlpLedInfoGet = dlsym(handle, "onlp_ledi_info_get");
    if ((error = dlerror()) != NULL)  {
        printf(error, stderr);
		return -1;
    }
    return 0;
}

int
Init()
{
	int ret = 0;
	int i = 0;
	onlp_oid_t* oidp;

    onlp_sys_info_t sysi;

	ret = loadOnlpSym();
    if (ret < 0) {
		printf("Error loading the ONLP symbols");
		return -1;
	}

	if (onlpInit == NULL) {
		printf("onlpInit ptr to onlp_init() is NULL\n");
		return -1;
	}
    (*onlpInit)();

    if (onlpSysInfoGet == NULL) {
		printf("onlpSysInfoGet ptr to onlp_sys_info_get is NULL\n");
		return -1;
	}
    ret = (*onlpSysInfoGet)(&sysi);
	if (ret < 0) {
		printf("onlp_sys_info_get() failed during init. Return Value:%d\n", ret);
		return ret;
	}

    ONLP_OID_TABLE_ITER_TYPE(sysi.hdr.coids, oidp, FAN) {
		fan_oid_table[i++] = *oidp;
	}

	return 0;
}

int
GetMaxNumOfFans()
{
	return ONLP_OID_TABLE_SIZE;
}

int
GetAllFanState(fan_info_t *info, int count)
{
	int i = 0;

    if (onlpFanInfoGet == NULL) {
		printf("onlpFanInfoGet Symbol not found \n");
		return -1;
	}

    for (i = 0; i < count; i++) {
		onlp_fan_info_t fi;
		int fid = ONLP_OID_ID_GET(fan_oid_table[i]);
		info[i].FanId = fid;
		info[i].valid = 0;

		if (fan_oid_table[i] == 0) {
			continue;
		}

        if ((*onlpFanInfoGet)(fan_oid_table[i], &fi) < 0) {
			printf("Failure retreiving status of Fan Id %d\n", fid);
			continue;
		}

        printf("Fan Id : %d Fan Status : %u\n", fid, fi.status);
		info[i].valid = 1;

        if (!(fi.status & 0x1)) {
			info[i].Status = FAN_STATUS_MISSING;
		} else if (fi.status & ONLP_FAN_STATUS_FAILED) {
			info[i].Status = FAN_STATUS_FAILED;
		} else if (fi.status & 0x1) {
			info[i].Status = FAN_STATUS_PRESENT;
		}

		info[i].Speed = fi.rpm;
		info[i].Direction = (fi.status & (1 << 2)) ? FAN_DIR_B2F : FAN_DIR_F2B;

        if (fi.rpm > 0) {
			info[i].Mode = FAN_MODE_ON;
		} else {
			info[i].Mode = FAN_MODE_OFF;
		}

		strncpy(info[i].Model, fi.model, DEFAULT_SIZE);
		strncpy(info[i].SerialNum, fi.serial, DEFAULT_SIZE);
	}
	return 0;
}

int
GetFanState(fan_info_t *info, int fanId)
{
	int i = 0, fid;

	if (onlpFanInfoGet == NULL) {
		printf("onlpFanInfoGet ptr to onlp_fan_info_get is NULL\n");
		return -1;
	}

    for (i = 0; i < sizeof(fan_oid_table)/sizeof(fan_oid_table[0]); i++) {
        onlp_fan_info_t fi;
		fid = ONLP_OID_ID_GET(fan_oid_table[i]);

        if (fid != fanId) {
			continue;
		}

        info[0].FanId = fid;
		info[0].valid = 0;

        if (fan_oid_table[i] == 0) {
			continue;
		}

        if ((*onlpFanInfoGet)(fan_oid_table[i], &fi) < 0) {
			printf("Failure retreiving status of Fan Id %d\n", fid);
			continue;
		}

        printf("Fan Id : %d Fan Status : %u\n", fid, fi.status);
		info[0].valid = 1;

        if (!(fi.status & 0x1)) {
			info[0].Status = FAN_STATUS_MISSING;
		} else if (fi.status & ONLP_FAN_STATUS_FAILED) {
			info[0].Status = FAN_STATUS_FAILED;
		} else if (fi.status & 0x1) {
			info[0].Status = FAN_STATUS_PRESENT;
		}

		info[0].Speed = fi.rpm;
		info[0].Direction = (fi.status & (1 << 2)) ? FAN_DIR_B2F : FAN_DIR_F2B;

        if (fi.rpm > 0) {
			info[0].Mode = FAN_MODE_ON;
		} else {
			info[0].Mode = FAN_MODE_OFF;
		}

		strncpy(info[0].Model, fi.model, DEFAULT_SIZE);
		strncpy(info[0].SerialNum, fi.serial, DEFAULT_SIZE);
		return 0;
	}
	return -1;
}

int
GetPlatformState(sys_info_t *info_p)
{
    int ret = 0;

    onlp_sys_info_t onlp_info;

    if (onlpSysInfoGet == NULL) {
		printf("onlpSysInfoGet ptr to onlp_sys_info_get is NULL\n");
		return -1;
	}

    ret = (*onlpSysInfoGet)(&onlp_info);
	if (ret < 0) {
		printf("onlp_sys_info_get() failed(%d)\n", ret);
		return ret;
	}

    strncpy(info_p->product_name, onlp_info.onie_info.product_name, DEFAULT_SIZE);
    strncpy(info_p->serial_number, onlp_info.onie_info.serial_number, DEFAULT_SIZE);
    strncpy(info_p->manufacturer, onlp_info.onie_info.manufacturer, DEFAULT_SIZE);
    strncpy(info_p->vendor, onlp_info.onie_info.vendor, DEFAULT_SIZE);
    strncpy(info_p->platform_name, onlp_info.onie_info.platform_name, DEFAULT_SIZE);
    strncpy(info_p->onie_version, onlp_info.onie_info.onie_version, DEFAULT_SIZE);
    strncpy(info_p->label_revision, onlp_info.onie_info.label_revision, DEFAULT_SIZE);

    return ret;
}

int
GetSfpCnt()
{
    int rt;
    onlp_sfp_bitmap_t bMap;

    (*onlpSfpBitmapInit)(&bMap);
    if (onlpSfpBitmapGet == NULL) {
		printf("onlpSfpBitmapGet NULL\n");
		return 0;
	}

    rt = (*onlpSfpBitmapGet)(&bMap);

    if (rt < 0)
        return 0;

    return AIM_BITMAP_COUNT(&bMap);
}

SFP_RET
GetSfpState(sfp_info_t *sfpInfo, int sfpId)
{
    int rt, rval;
    unsigned char *eeprom_p = NULL;

    memset(sfpInfo, 0, sizeof(sfp_info_t));
    if (!sfpId) {
        return SFP_MISSING;
    }

    rt = (*onlpSfpIsPresent)(sfpId);
    if (rt < 0) {
		printf("%s :(%d) INVALID SFP(%d) \n",__FUNCTION__, rt, sfpId);
		return SFP_ERROR;
    }

    if (!rt)
        return SFP_MISSING;

    sfpInfo->sfp_id = sfpId;
    sfpInfo->sfp_present = 1;

    /* RX LOS */
    (*onlpSfpControlGet)(sfpId, ONLP_SFP_CONTROL_RX_LOS, &rval);
    sfpInfo->sfp_los = rval?1:0;

    rt = (*onlpSfpEepromRead)(sfpId, &eeprom_p);
    if (rt < 0) {
		printf("%s :(%d) Eeprom read failed for SFP(%d)\n",__FUNCTION__, rt, sfpId);
		return SFP_ERROR;
    }
    memcpy(sfpInfo->eeprom, eeprom_p, 256);
    return SFP_OK;
}

SENSOR_RET
GetThermalState(thermal_info_t *info_p, int sensor_id)
{
    int rt;
    onlp_oid_t id;
    onlp_thermal_info_t t_info;

    id = ONLP_THERMAL_ID_CREATE(sensor_id + 1);

    if (!onlpThermalInfoGet)
        return SENSOR_ERROR;

    rt = (*onlpThermalInfoGet)(id, &t_info);

    if (rt < 0)
        return SENSOR_MISSING;

    memcpy(info_p->description, t_info.hdr.description, sizeof(t_info.hdr.description));
    info_p->sensor_id = sensor_id;
    info_p->status = t_info.status;
    info_p->caps = t_info.caps;
    info_p->temp = t_info.mcelsius;
    info_p->threshold_warning = t_info.thresholds.warning;
    info_p->threshold_error = t_info.thresholds.error;
    info_p->threshold_shutdown = t_info.thresholds.shutdown;

    return SENSOR_OK;
}

PSU_RET
GetPsuState(psu_info_t *info_p, int psu_id)
{
    int rt;
    onlp_oid_t id;
    onlp_psu_info_t p_info;

    id = ONLP_PSU_ID_CREATE(psu_id);

    if (!onlpPsuInfoGet)
        return PSU_ERROR;

    rt = (*onlpPsuInfoGet)(id, &p_info);


    if (rt < 0)
        return PSU_MISSING;

    info_p->psu_id = psu_id;
    info_p->status = (p_info.status & ONLP_PSU_STATUS_PRESENT);
    if (info_p->status) {
        info_p->mvin = p_info.mvin;
        info_p->mvout = p_info.mvout;
        info_p->miin = p_info.miin;
        info_p->miout = p_info.miout;
        info_p->mpin = p_info.mpin;
        info_p->mpout = p_info.mpout;
    }
    return PSU_OK;
}

LED_RET
GetLedState(led_info_t *info_p, int led_id)
{
    int rt;
    onlp_oid_t id;
    onlp_led_info_t led_info;

    id = ONLP_LED_ID_CREATE(led_id + 1);

    if (!onlpLedInfoGet)
        return LED_ERROR;

    rt = (*onlpLedInfoGet)(id, &led_info);


    if (rt < 0)
        return LED_MISSING;

    info_p->led_id = led_id;
    info_p->status = (led_info.status & ONLP_LED_STATUS_PRESENT);

    if (info_p->status) {
        memcpy(info_p->description, led_info.hdr.description, sizeof(led_info.hdr.description));
        if (led_info.mode != ONLP_LED_MODE_OFF)
            sprintf(info_p->color, "ON");
    } else {
        sprintf(info_p->description, "LED NOT PRESENT");
    }
    return LED_OK;
}

int
DeInit() {
	dlclose(handle);
}
