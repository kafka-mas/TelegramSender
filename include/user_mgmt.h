#ifndef USER_MGMT_H
#define USER_MGMT_H

#include <telebot.h>

typedef struct User_struct{
    long long int id;
    char name[8];
    bool verified;
}User_s;

typedef union User_union
{
    long long int id;
    char name[8];
}User_t;


void add_user(telebot_handler_t *handle, User_s *user);
// Not active
void delete_user();
#ifdef DEBUG
void send_something(telebot_handler_t *handle, long long int id);
#endif //DEBUG

#endif //USER_MGMT_H