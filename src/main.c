#include <stdio.h>
#include <stdlib.h>
#include <telebot.h>

#include "user_mgmt.h"
#include "config_read.h"

int main(int argc, char* argv[]){
    char *token = read_token();
    if (token == NULL){
        perror("Token not found");
        free(token);
        return -1;
    }
    telebot_handler_t handle;
    if (telebot_create(&handle, token) != TELEBOT_ERROR_NONE)
    {
        perror("Telebot create failed");
        return -1;
    }
    free(token);

    User_s user;
    add_user(&handle, &user);

    printf("|     ID     |   Name   | Ver |\n");
    printf("| %-11lli| %-9s|  %i  |\n", user.id, user.name, user.verified);

    telebot_destroy(handle);
    return 0;
}
