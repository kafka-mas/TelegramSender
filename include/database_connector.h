/**
 * @file database_connector.h
 * @author Kafka-mas (kafka_mas@disroot.org)
 * 
 * @copyright Copyright (c) 2026
 * 
 * @todo Setup library paths
 */

#ifndef DATABASE_CONNECTOR_H
#define DATABASE_CONNECTOR_H

#include <sqlite3.h>

#include "type_definition.h"

#ifdef DEBUG
    #define DATABASE_PATH PROJ_DIR "/debug/data/database.db"
#else
    #define DATABASE_PATH "/var/lib/telegram_sender/database.db"
#endif // DEBUG

#define SIZE_OF_ARRAY(array) (sizeof(array) / sizeof(array[0])) /**< Number of array elements */
#define UNICODE_LINE_LEN(length) ((length * 3) + 1) /**< Unicode-symbols line length */

/**
 * @brief Create a sqlite table object
 * 
 * @return int `0` if all good, `1` if something wrong
 */
int create_table();

/**
 * @brief Add new user to database
 * 
 * @param id telegram user ID
 * @param name username
 * @return int `0` if all good, `1` if something wrong
 */
int add_user_data(const long long int user_id, const char *name);

/**
 * @brief delete user by `id`/`user_id`/`username` 
 * 
 * @param parameter search parameter and value
 * @return int 1 if error otherwise 0
 */
int delete_user_from_db(const UserSearch *parameter);

/**
 * @brief Print all users to console
 * 
 * @return int `0` if all good, `1` if something wrong
 */
int print_all_users();

/**
 * @brief select user by id in database
 * @todo it is necessary to make it possible both by id and username
 * 
 * @param id [in] user database id
 * @param user [out] User type from `user_mgmt.h`
 * @return int 
 */
int select_user(const UserSearch *parameter, User *user);

/**
 * @brief Set the default user for sending
 * 
 * @return int 1 if error, else 0
 */
int set_default_user();

/**
 * @brief Get the default user id object (if setted)
 * 
 * @return int 1 if error, otherwise 0
 */
int get_default_user_id(long long int *user_id);

#endif // DATABASE_CONNECTOR_H
