#include <stdio.h>
#include <stdlib.h>
#include <telebot.h>
#include <locale.h>

#include "type_definition.h"
#include "database_connector.h"
#include "user_mgmt.h"
#include "config_read.h"

int main(int argc, char* argv[]){
    setlocale(LC_ALL, "");
    // char *token = read_token();
    // if (token == NULL){
    //     perror("Token not found");
    //     return -1;
    // }
    // telebot_handler_t handle;
    // if (telebot_create(&handle, token) != TELEBOT_ERROR_NONE)
    // {
    //     perror("Telebot create failed");
    //     return -1;
    // }
    // free(token);

    User user;
    // if(!add_user(&handle, &user)){
    //     return -1;
    // }

    // // create_table();
    // if(user.verified) add_user_data(user.user_id, user.name);

    // telebot_destroy(handle);
    return 0;
}
