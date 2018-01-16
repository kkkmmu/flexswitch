#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <syslog.h>
#include <errno.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <sys/errno.h>
// Ip Table/Net Filter Library includes
#include "libiptc/libiptc.h"
#include "linux/netfilter/xt_limit.h"
#include "linux/netfilter/xt_physdev.h"
#include "iptables.h"
#include "ipTable.h"

/*--------------------------------------*/
/* Compute netmask address given prefix */
/*--------------------------------------*/
static in_addr_t netmask( int prefix )
{
    if ( prefix == 0 )
        return((in_addr_t)- 1);//~((in_addr_t) -1) );
    else
        return( ~((1 << (32 - prefix)) - 1) );
} /* netmask() */

/*
Operation: ADD
retVal:
case -1: rule already exists
case -2: append failed
case -3: commit failed
case 1: new rule append and commit success

Operation: DELETE
retVal:
case -1: rule doesn't exists
case -2: append failed
case -3: commit failed
case 1: rule delete and commit success
*/

static int check_rule_and_operate(struct ipt_entry *ipEntry_p,
        struct xtc_handle *handle,
        const char* chain, 
        rule_operation_t operation,
        int *err_num)
{
    unsigned char *matchmask = NULL;
    int retVal = -1; 

    matchmask = malloc(ipEntry_p->next_offset);
    if (matchmask == NULL) {
        syslog(LOG_ERR, "failed to create match mask for delete");
        goto early_exit;
    }
    if (!iptc_check_entry(chain, ipEntry_p, matchmask, handle)) {
        switch (operation) {
            case ADD_RULE:
                if (!iptc_append_entry(chain, ipEntry_p, handle)) {
                    syslog(LOG_ERR, "append entry failed: %s", iptc_strerror(errno));
                    *err_num = errno;
                    retVal = -2;
                    goto early_exit;
                }
                retVal = 1;
                break;
            case DELETE_RULE:
                syslog(LOG_INFO, "entry doesn't exist, delete cannot proceed"
                        "but still sysd need to delete the entry from its database");
                retVal = -1; // No need for commit
                break;
            default:
                break;
        }
    } else { // Entry does exists
        switch (operation) {
            case ADD_RULE:
                syslog(LOG_INFO, "entry exist no need to append or commit");
                retVal = -1; // no need for commit
                break;
            case DELETE_RULE:
                if (!iptc_delete_entry(chain, ipEntry_p, matchmask, handle)) {
                    syslog(LOG_ERR, "delete entry failed, %s", iptc_strerror(errno));
                    retVal = -2;
                    *err_num = errno;
                    goto early_exit;
                }
                retVal = 1;
                break;
            default:
                break;
        }
    }
early_exit:
    if (retVal && !iptc_commit(handle)) {
        syslog(LOG_ERR, "commit failed, %s", iptc_strerror(errno));
        *err_num = errno;
        retVal = -3;
    }
    if (matchmask) {
        free(matchmask);
    }
    return retVal;
}

static int insert_rule(struct ipt_entry *ipEntry_p, const char* chain, bool restart,
        int *err_num)
{
    struct xtc_handle  *handle = NULL;
    int retVal = -1;

    // Create Table
    handle = iptc_init("filter");
    if (handle == NULL) {
        syslog(LOG_ERR, "cannot allocate memory to for iptc: %s",
                iptc_strerror(errno));
        goto early_exit;
    }
    if (!iptc_is_chain(chain, handle)) {
        syslog(LOG_ERR,"no such chain %s, error: %s", chain, iptc_strerror(errno));
        *err_num = errno;
        goto early_exit;
    }

    // Before Appending rule... check whether the entry already exists or not
    retVal = check_rule_and_operate(ipEntry_p, handle, chain, ADD_RULE, err_num);
    if (retVal == -1) {
        // suggests that linux has the entry but sysd restarted or system was
        // with saved iptables
        // As a result retVal should be updated informing sysd what to do with
        // the request
        syslog(LOG_INFO, "Rule already existis");
    } else  if ((retVal == -2) || (retVal == -3)) {
        // Adding new rule to chain failed
        syslog(LOG_ERR, "new rule addition failed");
        goto early_exit;
    } else {
        // Adding new rule to chain is success
        syslog(LOG_INFO, "successful commit check iptables");
    }

early_exit:

    if ((retVal == -2) || (retVal == -3)) {
        if (ipEntry_p) {
            free(ipEntry_p);
        }
    } else if (retVal == -1) { // duplicate entry
        // If non restart then inform sysd about this duplicate entry
        if (ipEntry_p && !restart) {
            syslog(LOG_WARNING, "duplicate rule create during non-restart");
            free(ipEntry_p);
        } else {
            syslog(LOG_INFO, "new rule create during restart case, return success");
            retVal = 1;
        }
    }
    if (handle) {
        iptc_free(handle);
    }
    return retVal;
}

static void fill_ip_entry(struct ipt_entry *ipEntry_p, const rule_entry_t *config,
        const unsigned int size_ipt_entry, __u16 proto)
{
    /************************* IP ENTRY ***************************/
    ipEntry_p->target_offset = size_ipt_entry;
    ipEntry_p->ip.proto = proto;

    // Ip Info.... work on netmask
    // If no ip then it will allow anywhere
    if (strlen(config->IpAddr) >= IP_ADDR_MIN_LENGTH) {
        ipEntry_p->ip.src.s_addr = inet_addr(config->IpAddr);
        ipEntry_p->ip.smsk.s_addr = htonl(netmask(config->PrefixLength));
    }

    // Physical Port info
    if (strlen(config->PhysicalPort) > 0) {
        strncpy(ipEntry_p->ip.iniface, config->PhysicalPort, 
                sizeof(ipEntry_p->ip.iniface));
    }

}

int add_iptable_tcp_rule(rule_entry_t *config, ipt_config_t *return_config_p)
{
    struct ipt_entry 		*ipEntry_p = NULL;
    struct ipt_entry_match *match_proto_p;
    struct ipt_standard_target *target_p;
    struct ipt_tcp * tcpinfo;
    unsigned int size_ipt_entry =0, size_ipt_entry_match =0, size_ipt_entry_target =0 ;
    unsigned int size_ipt_tcp=0, entry_size=0;
    int retVal = 0, err_num = 0;

    // Calculate structure length
    size_ipt_entry = XT_ALIGN(sizeof(struct ipt_entry));
    size_ipt_entry_match = XT_ALIGN(sizeof(struct ipt_entry_match));
    size_ipt_tcp = XT_ALIGN(sizeof(struct ipt_tcp));
    size_ipt_entry_target = XT_ALIGN(sizeof(struct ipt_standard_target));
    entry_size =  size_ipt_entry + size_ipt_entry_match + 
        size_ipt_entry_target + size_ipt_tcp;

    // Allocate memory
    ipEntry_p = (struct ipt_entry *) malloc(entry_size);
    if (ipEntry_p == NULL) {
        syslog(LOG_ERR, "No Memory for ip entry");
        return_config_p->entry = NULL;
        return -1;
    }
    bzero(ipEntry_p, entry_size);


    /************************* IP ENTRY ***************************/
    fill_ip_entry(ipEntry_p, config, size_ipt_entry, IPPROTO_TCP);


    /************************* MATCH ENTRY ***************************/
    ipEntry_p->target_offset = XT_ALIGN(ipEntry_p->target_offset + size_ipt_entry_match);
    match_proto_p = (struct ipt_entry_match *)ipEntry_p->elems;

    match_proto_p->u.user.match_size = XT_ALIGN(size_ipt_entry_match + size_ipt_tcp);
    strncpy(match_proto_p->u.user.name, config->Protocol,
            sizeof(match_proto_p->u.user.name)-2);


    /************************* TCP ENTRY ***************************/
    ipEntry_p->target_offset = XT_ALIGN(ipEntry_p->target_offset + size_ipt_tcp);
    tcpinfo = (struct ipt_tcp*)match_proto_p->data;
    // We don't care for src port and hence set it to 0 - 65535
    tcpinfo->spts[0] = 0; tcpinfo->spts[1] = MAX_PORT_NUM;
    if (config->Port != 0) {
        tcpinfo->dpts[0] = tcpinfo->dpts[1] = config->Port;
    } else {
        tcpinfo->dpts[0] = 0; tcpinfo->dpts[1] = MAX_PORT_NUM;
    }


    /************************* TARGET ENTRY ***************************/
    // Action Info
    target_p = (struct ipt_standard_target*)(((void *)ipEntry_p) + ipEntry_p->target_offset);
    strncpy(target_p->target.u.user.name, config->Action, 
            sizeof(target_p->target.u.user.name)-2);
    target_p->target.u.user.target_size = size_ipt_entry_target;
    ipEntry_p->next_offset = XT_ALIGN(ipEntry_p->target_offset + size_ipt_entry_target);

    //retVal = insert_rule(ipEntry_p, INPUT_CHAIN, config->Restart); 
    retVal = insert_rule(ipEntry_p, INPUT_CHAIN, config->Restart, &err_num); 
    if (retVal <= 0) {
        return_config_p->err_num = err_num;
        return_config_p->entry = NULL;
    } else {
        return_config_p->entry = ipEntry_p;
        strncpy(return_config_p->name, config->Name, 
                sizeof(return_config_p->name));
    }
    return retVal;
}

int add_iptable_udp_rule(rule_entry_t *config, ipt_config_t *return_config_p)
{
    struct ipt_entry 		*ipEntry_p = NULL;
    struct ipt_entry_match *match_proto_p;
    struct ipt_standard_target *target_p;
    struct ipt_udp * udpinfo;
    unsigned int size_ipt_entry =0, size_ipt_entry_match =0, size_ipt_entry_target =0 ;
    unsigned int size_ipt_udp=0, entry_size=0;
    int retVal = 0, err_num = 0;

    // Calculate structure length
    size_ipt_entry = XT_ALIGN(sizeof(struct ipt_entry));
    size_ipt_entry_match = XT_ALIGN(sizeof(struct ipt_entry_match));
    size_ipt_udp = XT_ALIGN(sizeof(struct ipt_udp));
    size_ipt_entry_target = XT_ALIGN(sizeof(struct ipt_standard_target));
    entry_size =  size_ipt_entry + size_ipt_entry_match + 
        size_ipt_entry_target + size_ipt_udp;

    // Allocate memory
    ipEntry_p = (struct ipt_entry *) malloc(entry_size);
    if (ipEntry_p == NULL) {
        syslog(LOG_ERR, "No Memory for ip entry");
        return_config_p->entry = NULL;
        return -1;
    }
    bzero(ipEntry_p, entry_size);


    /************************* IP ENTRY ***************************/
    fill_ip_entry(ipEntry_p, config, size_ipt_entry, IPPROTO_UDP);


    /************************* MATCH ENTRY ***************************/
    ipEntry_p->target_offset = XT_ALIGN(ipEntry_p->target_offset + size_ipt_entry_match);
    match_proto_p = (struct ipt_entry_match *)ipEntry_p->elems;

    match_proto_p->u.user.match_size = XT_ALIGN(size_ipt_entry_match + size_ipt_udp);
    strncpy(match_proto_p->u.user.name, config->Protocol,
            sizeof(match_proto_p->u.user.name)-2);


    /************************* UDP ENTRY ***************************/
    ipEntry_p->target_offset = XT_ALIGN(ipEntry_p->target_offset + size_ipt_udp);
    udpinfo = (struct ipt_udp*)match_proto_p->data;
    // We don't care for src port and hence set it to 0 - 65535
    udpinfo->spts[0] = 0; udpinfo->spts[1] = MAX_PORT_NUM;
    if (config->Port != 0) {
        udpinfo->dpts[0] = udpinfo->dpts[1] = config->Port;
    } else {
        udpinfo->dpts[0] = 0; udpinfo->dpts[1] = MAX_PORT_NUM;
    }


    /************************* TARGET ENTRY ***************************/
    // Action Info
    target_p = (struct ipt_standard_target*)(((void *)ipEntry_p) + ipEntry_p->target_offset);
    strncpy(target_p->target.u.user.name, config->Action, 
            sizeof(target_p->target.u.user.name)-2);
    target_p->target.u.user.target_size = size_ipt_entry_target;
    ipEntry_p->next_offset = XT_ALIGN(ipEntry_p->target_offset + size_ipt_entry_target);

    retVal = insert_rule(ipEntry_p, INPUT_CHAIN, config->Restart, &err_num); 
    if (retVal <= 0) {
        return_config_p->err_num = err_num;
        return_config_p->entry = NULL;
    } else {
        return_config_p->entry = ipEntry_p;
        strncpy(return_config_p->name, config->Name, 
                sizeof(return_config_p->name));
    }
    return retVal;
}

int add_iptable_icmp_rule(rule_entry_t *config, ipt_config_t *return_config_p)
{
    struct ipt_entry 		*ipEntry_p = NULL;
    struct ipt_entry_match *match_proto_p;
    struct ipt_standard_target *target_p;
    struct ipt_icmp *icmpinfo;
    unsigned int size_ipt_entry =0, size_ipt_entry_match =0, size_ipt_entry_target =0 ;
    unsigned int size_ipt_icmp=0, entry_size=0;
    int retVal = 0, err_num = 0;

    // Calculate structure length
    size_ipt_entry = XT_ALIGN(sizeof(struct ipt_entry));
    size_ipt_entry_match = XT_ALIGN(sizeof(struct ipt_entry_match));
    size_ipt_icmp = XT_ALIGN(sizeof(struct ipt_icmp));
    size_ipt_entry_target = XT_ALIGN(sizeof(struct ipt_standard_target));
    entry_size =  size_ipt_entry + size_ipt_entry_match + 
        size_ipt_entry_target + size_ipt_icmp;

    // Allocate memory
    ipEntry_p = (struct ipt_entry *) malloc(entry_size);
    if (ipEntry_p == NULL) {
        syslog(LOG_ERR, "No Memory for ip entry");
        return_config_p->entry = NULL;
        return -1;
    }
    bzero(ipEntry_p, entry_size);


    /************************* IP ENTRY ***************************/
    fill_ip_entry(ipEntry_p, config, size_ipt_entry, IPPROTO_ICMP);


    /************************* MATCH ENTRY ***************************/
    ipEntry_p->target_offset = XT_ALIGN(ipEntry_p->target_offset + size_ipt_entry_match);
    match_proto_p = (struct ipt_entry_match *)ipEntry_p->elems;

    match_proto_p->u.user.match_size = XT_ALIGN(size_ipt_entry_match + size_ipt_icmp);
    strncpy(match_proto_p->u.user.name, config->Protocol,
            sizeof(match_proto_p->u.user.name)-2);


    /************************* ICMP ENTRY ***************************/
    ipEntry_p->target_offset = XT_ALIGN(ipEntry_p->target_offset + size_ipt_icmp);
    icmpinfo = (struct ipt_icmp*)match_proto_p->data;
    // We do not support icmp type... so set it to 255 for all icmp type
    icmpinfo->type = ICMP_PORT_MAX;
    // We don't care for code and hence set it to 0 - 255
    icmpinfo->code[0] = 0; icmpinfo->code[1] = ICMP_PORT_MAX;

    /************************* TARGET ENTRY ***************************/
    // Action Info
    target_p = (struct ipt_standard_target*)(((void *)ipEntry_p) + ipEntry_p->target_offset);
    strncpy(target_p->target.u.user.name, config->Action, 
            sizeof(target_p->target.u.user.name)-2);
    target_p->target.u.user.target_size = size_ipt_entry_target;
    ipEntry_p->next_offset = XT_ALIGN(ipEntry_p->target_offset + size_ipt_entry_target);

    //retVal = insert_rule(ipEntry_p, INPUT_CHAIN, config->Restart); 
    retVal = insert_rule(ipEntry_p, INPUT_CHAIN, config->Restart, &err_num); 
    if (retVal <= 0) {
        return_config_p->err_num = err_num;
        return_config_p->entry = NULL;
    } else {
        return_config_p->entry = ipEntry_p;
        strncpy(return_config_p->name, config->Name, 
                sizeof(return_config_p->name));
    }
    return retVal;
}

int del_iptable_rule(ipt_config_t *cfg_entry_p)
{
    struct xtc_handle  *handle = NULL;
    int retVal = -1, err_num = 0;

    handle = iptc_init("filter");
    if (handle == NULL) {
        syslog(LOG_ERR, "cannot allocate memory to for iptc: %s",
                iptc_strerror(errno));
        err_num = errno;
        goto early_exit;
    }

    retVal = check_rule_and_operate(cfg_entry_p->entry, handle, INPUT_CHAIN,
            DELETE_RULE, &err_num);
    if (retVal == -1) {
        syslog(LOG_ERR, "rule doesn't exists, sysd delete your db info");
    } else if ((retVal == -2) || (retVal == -3)) {
        syslog(LOG_ERR, "executing delete entry failed");
        retVal = -1;
        goto early_exit;
    } else {
        syslog(LOG_INFO, "rule delete and commit success"); 
    }
    retVal = 1;

early_exit:

    if (handle) {
        iptc_free(handle);
    }

    if (retVal == 1) {
        if (cfg_entry_p->entry) {
            free(cfg_entry_p->entry);
        }
    }
    cfg_entry_p->err_num = errno;
    return retVal;
}

void get_iptc_error_string(err_t *err_Info, int err_num)
{
    strncpy(err_Info->err_string, iptc_strerror(err_num), 
            sizeof(err_Info->err_string));
}
