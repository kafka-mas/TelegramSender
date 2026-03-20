#ifndef USER_MGMT_H
#define USER_MGMT_H

#include <telebot.h>

#include "type_definition.h"

#define MESSAGES_LIMIT 20 /**< Max get messages per request update*/
#define VERIFICATION_TOKEN_LENGTH_BYTES 24 /**< Max length getting from /dev/urandom for generate token */
#define VERIFICATION_TOKEN_LENGTH VERIFICATION_TOKEN_LENGTH_BYTES * 4 / 3 + 1 /**< Token length in symbols + `'\0'` */
#define POLLING_TIMEOUT 30  /**< Timeout in seconds for long polling*/
#define MAX_WAIT_TIME 150 /**< Timeout in seconds for auth */
#define UPDATES_COUNT 1 /**< Count of updates from `update_types`*/


/**
 * @brief This function is used for verify user and get Telegram first name and user ID
 * 
 * @param handle [in] `telebot_handler_t` from telebot library.
 * @param user [out] `struct User_struct` returns filled struct with user data.
 * @return true on success 
 * @return false if something going wrong
 */
bool add_user(telebot_handler_t *handle, User *user);

/**
 * @brief Not created yet
 * 
 */
int delete_user(UserSearch user);

#ifdef DEBUG
/**
 * @brief Send word `something` in specified chat
 * 
 * @param handle [in] `telebot_handler_t` from telebot library.
 * @param id [in] Telegram user ID.
 */
void send_something(telebot_handler_t *handle, long long int user_id);
#endif //DEBUG

#endif //USER_MGMT_H
