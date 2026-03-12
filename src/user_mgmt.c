#include <user_mgmt.h>

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
#include <telebot.h>
#include <fcntl.h>

#include <openssl/evp.h>
#include <openssl/buffer.h>

static void generate_random_token(char *token, size_t size);

static int send_message(telebot_handler_t *handle, long long int id, char *message_s);

bool add_user(telebot_handler_t *handle, User *user){
    #ifdef DEBUG
        char rand_token[] = "123";
    #else
        char rand_token[VERIFICATION_TOKEN_LENGTH];
        generate_random_token(rand_token, sizeof(rand_token));
        if(rand_token[0] == '\0'){
            perror("Error generate verification token");
            return false;
        }
    #endif // DEBUG
    printf("Verification token: %s\n", rand_token);

    telebot_error_e ret;
    telebot_update_t *updates;
    int count;
    int offset = 0; 
    
    const int limit = 100;
    ret = telebot_get_updates(*handle, offset, limit, 0, NULL, 0, &updates, &count);
    if (ret == TELEBOT_ERROR_NONE && count > 0) {
        offset = updates[count-1].update_id + 1;
        telebot_put_updates(updates, count);
    }

    telebot_message_t message;
    telebot_update_type_e update_types[] = {TELEBOT_UPDATE_TYPE_MESSAGE};

/** @todo
 * 10. Отсутствие обработки сигналов
 * Проблема: При ожидании верификации программа может быть прервана сигналом, и пользователь не узнает о неудаче.
 * Решение: Можно добавить обработку SIGINT, чтобы корректно завершить работу и вернуть ошибку.
 */

    int index;
    bool verified = false;
    #ifdef DEBUG
        while (!verified)
    #else
        time_t start_time = time(NULL);
        while (!verified)
    #endif
    {
        telebot_update_t *updates;
        ret = telebot_get_updates(*handle, offset, MESSAGES_LIMIT, POLLING_TIMEOUT, update_types, UPDATES_COUNT, &updates, &count);
        if (ret != TELEBOT_ERROR_NONE){
            sleep(1);
            continue;
        }
        for (index = 0; index < count; index++)
        {
            message = updates[index].message;
            if (message.text)
            {
                if (strstr(message.text, "/start") || strstr(message.text, "/verify")){
                    send_message(handle, message.from->id, "Send your token.");

                } else if (strcmp(rand_token, message.text) == 0) {
                    verified = true;
                    printf("%s\n", message.text);

                    user->id = message.from->id;
                    snprintf(user->name, sizeof(user->name), "%s", message.from->first_name);
                    user->verified = true;

                    send_message(handle, message.from->id, "Your account succesfully added!!!");
                    break;
                } else {
                    send_message(handle, message.from->id, "Bad token; Try again.");
                }
            }
            offset = updates[index].update_id + 1;
        }
        telebot_put_updates(updates, count);
        #ifndef DEBUG
            if(time(NULL) - start_time < MAX_WAIT_TIME){
                perror("Too long auth");
                return false;
            }
        #endif
    }
    return true;
}

#ifdef DEBUG
void send_something(telebot_handler_t *handle, long long int id){
    telebot_error_e ret;
    char *str = "`something`";
    ret = telebot_send_message(*handle, id, str, "Markdown", false, false, 0, "");
    if(ret != TELEBOT_ERROR_NONE){
        perror("Error while send message");
    }
}
#endif

int delete_user(UserSearch user){
    printf("---------------------\n");
    if (user.type == USER_SEARCH_BY_ID) {
        printf("Id: %lli\n", user.value.id);
    } else if (user.type == USER_SEARCH_BY_NAME) {
        printf("Name: %s\n", *user.value.name);
    } else return 1;
    printf("---------------------\n");


    return 0;
}

static void generate_random_token(char *token, size_t size){
    unsigned char random_bytes[VERIFICATION_TOKEN_LENGTH_BYTES];

    int fd = open("/dev/urandom", O_RDONLY);
    if (fd < 0){
        perror("Error open /dev/urandom");
        token[0] = '\0';
        return;
    }

    if (read(fd, random_bytes, sizeof(random_bytes)) < 0){
        perror("Error reading /dev/urandom");
        close(fd);
        token[0] = '\0';
        return;
    }

    close(fd);

    BIO *bio, *b64;
    BUF_MEM *bufferPtr;

    b64 = BIO_new(BIO_f_base64());
    bio = BIO_new(BIO_s_mem());

    if (!b64 || !bio) {
        token[0] = '\0';
        BIO_free_all(bio);
        BIO_free_all(b64);
        return;
    }

    bio = BIO_push(b64, bio);
    BIO_set_flags(bio, BIO_FLAGS_BASE64_NO_NL);
    BIO_write(bio, random_bytes, sizeof(random_bytes));
    BIO_flush(bio);
    BIO_get_mem_ptr(bio, &bufferPtr);

    if (bufferPtr->length < size) {
        memcpy(token, bufferPtr->data, bufferPtr->length);
        token[bufferPtr->length] = '\0';
    } else {
        token[0] = '\0';
    }

    BIO_free_all(bio);
}

static int send_message(telebot_handler_t *handle, long long int id, char *message_s){
    telebot_error_e ret;
    ret = telebot_send_message(*handle, id, message_s, "Markdown", false, false, 0, "");
    if(ret != TELEBOT_ERROR_NONE){
        perror("Error while send message");
        return 1;
    }
    return 0;
}
