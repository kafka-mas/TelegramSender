#ifndef CONFIG_READ_H
#define CONFIG_READ_H

#include <arpa/inet.h>

#ifdef DEBUG
    #define CONFIG_FILE_PATH PROJ_DIR "/debug/configs/sender.conf"
#else
    #define CONFIG_FILE_PATH "/etc/telegram_sender/sender.conf"
#endif // DEBUG

#define BUFFER_LENGTH 255 /**< File read buffer length */
#define MAX_PORT_STRLEN 6 /**< Conn. port buffer */

/**
 * @brief data from config file
 * 
 */
typedef struct
{
    char token[BUFFER_LENGTH];
    char proxy_address[INET6_ADDRSTRLEN];
    char proxy_port[MAX_PORT_STRLEN];
    char proxy_user[BUFFER_LENGTH];
    char proxy_password[BUFFER_LENGTH];
} Config;

/**
 * @brief Read config file defined by CONFIG_FILE_PATH
 * 
 * @param config [out] Config struct
 * @return int 1 if error, 0 otherwise
 */
int read_conf(Config *config);

#endif // CONFIG_READ_H