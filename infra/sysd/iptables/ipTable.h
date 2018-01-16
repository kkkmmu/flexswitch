#ifndef IP_TABLES_H
#define IP_TABLES_H
#include <stdint.h>
#include <stdbool.h>

#define IP_ADDR_MIN_LENGTH 8
#define ERROR_STRING_SIZE 512
#define MAX_PORT_NUM 0xFFFF
#define ICMP_PORT_MAX 0xFF
#define INPUT_CHAIN "INPUT"
#define RULE_NAME_SIZE 64

typedef enum operation_s {
    UNKNOWN = 0,
    ADD_RULE = 1,
    DELETE_RULE = 2,
}rule_operation_t;

typedef struct rule_entry_s {
    char *Name; 
    char *PhysicalPort; 
    char *Action; 
    char *IpAddr; 
    char *Protocol; 
    uint16_t  Port;
    int  PrefixLength;
    bool Restart;
}rule_entry_t;

typedef struct ipt_config_s {
    char   name[RULE_NAME_SIZE];
    int    err_num;
    struct ipt_entry *entry;
}ipt_config_t;

typedef struct err_s {
    char err_string[ERROR_STRING_SIZE];
}err_t;

// ADD RULE
int add_iptable_tcp_rule(rule_entry_t *config, ipt_config_t *rc);
int add_iptable_udp_rule(rule_entry_t *config, ipt_config_t *rc);
int add_iptable_icmp_rule(rule_entry_t *config, ipt_config_t *return_config_p);

// DELETE RULE
int del_iptable_rule(ipt_config_t *config);

// Error Info
void get_iptc_error_string(err_t *err_Info, int err_num);

#endif
