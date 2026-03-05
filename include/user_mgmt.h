#ifndef USER_MGMT_H
#define USER_MGMT_H

#include <telebot.h>

typedef struct{
    long long int id;
    char name[30];
    bool verified;
}User;


void add_user(telebot_handler_t *handle, User *user);
void send_something(telebot_handler_t *handle, long long int id);
// static unsigned char *generate_random_token();

#endif
