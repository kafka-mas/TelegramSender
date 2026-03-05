#include <config_read.h>

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>

char *read_token(){

    #ifdef DEBUG
        #define CONFIG_FILE_PATH PROJ_DIR "/debug/configs/sender.conf"
    #else
        #define CONFIG_FILE_PATH "/etc/telegram_sender/sender.conf"
    #endif

    FILE * fp = fopen(CONFIG_FILE_PATH, "r");

    if(fp==NULL)
    {
        perror("Error occured while opening config file");
        return NULL;
    }
    printf("Found config file at: %s\n", CONFIG_FILE_PATH);

    char buffer[256];
    int line_num = 0;
    while((fgets(buffer, 256, fp))!=NULL)
    {
        line_num++;

        buffer[strcspn(buffer, "\n")] = '\0';

        if (buffer[0] == '\0' || buffer[0] == '#') continue;

        char* eq_pos = strchr(buffer, '=');
        if (!eq_pos) continue;

        *eq_pos = '\0';
        char* name = buffer;
        char* value= eq_pos + 1;

        if(strcmp("token", name) == 0){
            char *result = malloc(strlen(value)+1);
            if(result) strcpy(result, value);
            fclose(fp);
            return result;
        }
    }

    fclose(fp);

    return NULL;
}