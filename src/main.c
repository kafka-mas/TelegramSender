#include <stdio.h>
#include <stdlib.h>
#include <telebot.h>

#include "user_mgmt.h"
#include "config_read.h"

int main(int argc, char* argv[]){
    char *token = read_token();
    if (token == NULL){
        perror("Token not found");
        return -1;
    }
    telebot_handler_t handle;
    if (telebot_create(&handle, token) != TELEBOT_ERROR_NONE)
    {
        perror("Telebot create failed");
        return -1;
    }
    free(token);

    User user;
    if(!add_user(&handle, &user)){
        return -1;
    }

    printf("|     ID     |       Name       | Ver |\n");
    printf("|%11lli | %16s |  %i  |\n", user.id, user.name, user.verified);

    delete_user((UserSearch){.type=USER_SEARCH_BY_ID, .value.id=user.id});
    delete_user((UserSearch){.type=USER_SEARCH_BY_NAME, .value.name=user.name});

    telebot_destroy(handle);
    return 0;
}
