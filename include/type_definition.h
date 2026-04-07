#ifndef TYPE_DEFINITION_H
#define TYPE_DEFINITION_H

#include <stdbool.h>

#define NAME_LENGTH 64 /**< User's first name in telegram */

/**
 * @brief Describe user to put into DB
 */
typedef struct User_struct{
    char name[NAME_LENGTH];
    long long int user_id;
    bool verified;
}User;

/**
 * @brief Search parameter (id or name)
 * 
 */
typedef enum {
    USER_SEARCH_BY_USER_ID,
    USER_SEARCH_BY_NAME,
    USER_SEARCH_BY_ID
} UserSearchType;

/**
 * @brief Search user in DB by name or by ID
 * 
 */
typedef struct {
    UserSearchType type;
    union {
        char name[NAME_LENGTH];
        long long int user_id;
        int id;
    } value;
} UserSearch;

#endif // TYPE_DEFINITION_H