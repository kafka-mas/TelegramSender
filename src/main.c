#include <stdio.h>
#include <stdlib.h>
#include <telebot.h>
#include <locale.h>
#include <getopt.h>
#include <unistd.h>
#include <string.h>

#include "package_info.h"
#include "type_definition.h"
#include "database_connector.h"
#include "user_mgmt.h"
#include "config_read.h"

/**
 * @brief Describe CLI parameters.
 * 
 */
int is_add_user = 0;
int is_set_default = 0;
// int is_send_file = 0;
int is_send_to_default = 0;
int is_specify_user = 0;
int is_show_help = 0;
int is_show_version = 0;
int is_create_db = 0;
int is_list_users = 0;
int is_use_stdin = 0;
// char *filename = NULL;
char *given_user_id = NULL;

/**
 * @brief Generate help message
 * 
 * @return int 1 if error, 0 otherwise
 */
int show_help();

/**
 * @brief Generate version message
 * 
 * @return int 1 if error, 0 otherwise 
 */
int show_ver();

int main(int argc, char* argv[]){
    setlocale(LC_ALL, "");

    struct option long_options[] = {
        {"add-user",    no_argument,        0,              'a'},
        {"set-default", no_argument,        &is_set_default,  1},
        {"create-db",   no_argument,        &is_create_db,    1},
        // {"file",        required_argument,  0,              'f'},
        {"list-users",  no_argument,        0,              'a'},
        {"user",        required_argument,  0,              'u'},
        {"default",     no_argument,        0,              'd'},
        {"version",     no_argument,        0,              'v'},
        {"help",        no_argument,        0,              'h'},
        {0, 0, 0, 0}
    };

    int arg;
    while ((arg = getopt_long(argc, argv, "u:adhvl", long_options, NULL)) != -1) {
        switch (arg)
        {
        case 'u':
            is_specify_user = 1;
            given_user_id = optarg;
            break;
        // case 'f':
        //     is_send_file = 1;
        //     filename = optarg;
        //     if (strcmp(filename, "-") == 0) {
        //         use_stdin = 1;
        //         filename = NULL;
        //     }
        //     break;
        case 'a':
            is_add_user = 1;
            break;
        case 'd':
            is_send_to_default = 1;
            break;
        case 'h':
            is_show_help = 1;
            break;
        case 'l':
            is_list_users = 1;
            break;
        case 'v':
            is_show_version = 1;
            break;
        default:
            break;
        }
    }

    if(is_show_version) {show_ver(); return 0;}
    if(is_show_help) {show_help(); return 0;}
// int is_add_user = 0;
// int is_set_default = 0;
// int is_send_file = 0;
// int is_send_to_default = 0;
// int is_specify_user = 0;
// int is_show_help = 0;
// int is_show_version = 0;
    if(is_specify_user&&is_send_to_default){
        fprintf(stderr, "Error: incompatible launch options\n");
        return -1;
    }

    if(is_add_user + is_set_default + is_create_db + is_list_users /*+ is_send_file*/ > 1){
        fprintf(stderr, "Error: incompatible launch options\n");
        return -1;
    }
    
    if(is_create_db){
        if(create_table() != 0){
            fprintf(stderr, "Error: cant't create DB\n");
            return  -1;
        }
        return 0;
    }

    if(is_set_default){
        if(set_default_user() == 1) return 1;
        return 0;
    }

    if(is_list_users){
        if(print_all_users() == 1) return 1;
        return 0;
    }

    FILE *input = NULL;

    /**
     * @brief Start bot
     * 
     */
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


    if(is_add_user){
        User user;
        if(!add_user(&handle, &user)){
            telebot_destroy(handle);
            return -1;
        }
        if(!user.verified) return -1;
        if(add_user_data(user.user_id, user.name) != 0)
        {
            fprintf(stderr, "Error: cant't add user in DB\n");
            telebot_destroy(handle);
            return  -1;
        }
        print_all_users();
        telebot_destroy(handle);
        return 0;
    }

    if (!isatty(STDIN_FILENO)) {
        is_use_stdin = 1;
    }

    if(is_use_stdin){
        input = stdin;
        int responce = 0;

        char buffer[256];
        char msg[4096];
        size_t pos = 0;
        for(pos; pos < 3; pos++){
            msg[pos] = '`';
        }msg[3] = '\n'; pos++;

        while (fgets(buffer, sizeof(buffer), input) != NULL) {
            size_t len = strlen(buffer);
            if (pos + len < sizeof(msg)) {
                memcpy(msg + pos, buffer, len);
                pos += len;
            } else {
                fprintf(stderr, "Message too long\n");
                responce = 1;
                break;
            }
        }
        if(pos + 4 >= 4096){
            msg[pos-4] = '\n';
            for(int i = 3; i > 0; i--){
                msg[pos-i] = '`';
            }msg[pos] = '\0';
        }
        else{
            msg[pos++] = '\n';
            for(size_t i = pos; i < pos + 3; i++){
                msg[i] = '`';
            }msg[pos+5] = '\0';
        }

        long long int user_id = 0;
        if(is_send_to_default)
        {
            if(get_default_user_id(&user_id) == 1){
                telebot_destroy(handle);
                return 1;
            }
        }

        send_text(&handle, msg, user_id);

        return responce;
    }



    // printf("\nis_add_user: %i\n", is_add_user);
    // printf("is_specify_user: %i. Arg: %s\n", is_specify_user, given_user_id);
    // // printf("is_send_file: %i. Arg: %s\n", is_send_file, filename);
    // printf("is_set_default: %i\n", is_set_default);
    // printf("is_show_help: %i\n", is_show_help);
    // printf("is_send_to_default: %i\n\n", is_send_to_default);
    show_help();
    telebot_destroy(handle);
    return 0;
}

int show_help(){
    const char *help = "Usage: " BIN_NAME " [OPTIONS] ...\n"
                       "  -a        --add-user      Add new user\n"                         //< Done
                       "            --create-db     Create database (drop if exist)\n"      //< Done
                       "  -d        --default       Send to default user\n"                 //< Not created yet
                    //    "  -f        --file          Specify input file\n"                   //< Not created yet
                       "  -h        --help          Show help\n"                            //< Done
                       "  -l        --list-users    Show all users\n"                       //< Done
                       "            --set-default   Select default user in database\n"      //< Done
                       "  -u        --user          Specify user to send\n"                 //< Done
                       "  -v        --version       Show version info\n"                    //< Done
                       ;

    printf("%s\n", help);
    return 0;
}

int show_ver(){
    const char *ver = PROJECT_NAME " " PROJECT_VERSION "\n";

    printf("%s\n", ver);
    return 0;
}