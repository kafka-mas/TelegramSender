#ifndef DATABASE_CONNECTOR_H
#define DATABASE_CONNECTOR_H

#include <sqlite3.h>

#ifdef DEBUG
    #define DATABASE_PATH PROJ_DIR "/debug/data/database.db"
#else
    #define DATABASE_PATH "/var/lib/telegram_sender/database.db"
#endif // DEBUG

int create_table();
int add_user_data(long long int id, char *name);
int delete_user_data();
int get_all_users();

#endif // DATABASE_CONNECTOR_H
