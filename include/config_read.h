#ifndef CONFIG_READ_H
#define CONFIG_READ_H

#ifdef DEBUG
    #define CONFIG_FILE_PATH PROJ_DIR "/debug/configs/sender.conf"
#else
    #define CONFIG_FILE_PATH "/etc/telegram_sender/sender.conf"
#endif // DEBUG

#define BUFFER_LENGTH 255 /**< File read buffer length */

/**
 * @brief Read config file defined by CONFIG_FILE_PATH
 * 
 * @return char* dynamically allocated telegram bot token
 */
char *read_token();

#endif // CONFIG_READ_H