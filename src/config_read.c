#include <config_read.h>

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <ctype.h>

char *read_token(){
    FILE * fp = fopen(CONFIG_FILE_PATH, "r");

    if(fp==NULL)
    {
        perror("Error occured while opening config file");
        return NULL;
    }
    printf("Found config file at: %s\n", CONFIG_FILE_PATH);

    char buffer[BUFFER_LENGTH];
    while((fgets(buffer, BUFFER_LENGTH, fp))!=NULL)
    {
        buffer[strcspn(buffer, "\n")] = '\0';

        if (buffer[0] == '\0' || buffer[0] == '#') continue;

        char* eq_pos = strchr(buffer, '=');
        if (!eq_pos) continue;

        *eq_pos = '\0';
        char* name = buffer;
        char* value= eq_pos + 1;

        if(strcmp("token", name) == 0 && value[0] != '\0'){
            char *result = strdup(value);
            fclose(fp);
            return result;
        }
    }

    fclose(fp);

    return NULL;
}