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

#define VERIFICATION_TOKEN_LENGTH_BYTES 24
#define VERIFICATION_TOKEN_LENGTH VERIFICATION_TOKEN_LENGTH_BYTES * 4 / 3 + 1

#define SIZE_OF_ARRAY(array) (sizeof(array) / sizeof(array[0]))

static void generate_random_token(char *token, size_t size);

void add_user(telebot_handler_t *handle, User *user){
    char rand_token[VERIFICATION_TOKEN_LENGTH];
    generate_random_token(rand_token, sizeof(rand_token));
    if(rand_token[0] == '\0'){
        perror("Error generate verification token");
        return;
    }
    printf("Verification token: %s\n", rand_token);

    telebot_error_e ret;
    telebot_update_t *updates;
    int count;
    int offset = 0; 
    
    ret = telebot_get_updates(*handle, offset, 100, 0, NULL, 0, &updates, &count);
    if (ret == TELEBOT_ERROR_NONE && count > 0) {
        offset = updates[count-1].update_id + 1;
        telebot_put_updates(updates, count);
    }

    telebot_message_t message;
    telebot_update_type_e update_types[] = {TELEBOT_UPDATE_TYPE_MESSAGE};

    int64_t expected_user_id;
    int index;
    bool verified = false;
    #ifdef DEBUG
    while (!verified)
    #else
    const int MAX_WAIT = 20;
    time_t start_time = time(NULL);
    while (!verified && (time(NULL) - start_time < MAX_WAIT))
    #endif
    {
        telebot_update_t *updates;
        ret = telebot_get_updates(*handle, offset, 20, 30, update_types, 1, &updates, &count);
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
                    expected_user_id = message.from->id;
                    // printf("%s: %s \n", message.from->first_name, message.text);
                    ret = telebot_send_message(*handle, message.from->id, "Send your token.", "Markdown", false, false, 0, "");
                } else if (message.from->id == expected_user_id && strcmp(rand_token, message.text) == 0) {
                    verified = true;
                    printf("%s\n", message.text);
                    ret = telebot_send_message(*handle, message.from->id, "Your account succesfully added!!!", "Markdown", false, false, 0, "");
                    break;
                } else if (expected_user_id != 0 && message.from->id == expected_user_id) {
                    ret = telebot_send_message(*handle, message.from->id, "Bad token; Try again.", "Markdown", false, false, 0, "");
                }
            }
            offset = updates[index].update_id + 1;
        }
        telebot_put_updates(updates, count);
    }
}

void send_something(telebot_handler_t *handle, long long int id){
    telebot_error_e ret;
    char *str = "`something`";
    ret = telebot_send_message(*handle, id, str, "Markdown", false, false, 0, "");
    if(ret != TELEBOT_ERROR_NONE){
        perror("Error while send message");
    }
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
