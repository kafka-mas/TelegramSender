#include <database_connector.h>

#include <stdio.h>
#include <stddef.h>
#include <string.h>
#include <sqlite3.h>

/**
 * @brief For `generate_table`. Set characters set for table.
 * 
 */
typedef enum {
    TOP_LINE,
    MID_LINE,
    BOT_LINE
} LINE_TYPE;

/**
 * @brief generate dividing lines for table
 * 
 * @param str pointer to string
 * @param line_type dividing line type enumerated in `LINE_TYPE` 
 * @param length line length in symbols
 * @param positions positions where “jumpers” for vertical dividers are inserted
 */
static void generate_line(char *str, const LINE_TYPE line_type, const int length, const int *positions, size_t pos_count);

int create_table() {
    sqlite3 *db;    // указатель на базу данных
    char *err_msg = 0;  // сообщение об ошибке
    // создаем базу данных
    int rc  = sqlite3_open(DATABASE_PATH, &db);
    // если подключение прошло неудачно
    if (rc != SQLITE_OK)
    {
        sqlite3_close(db);
        return 1;
    }

    char *sql = "DROP TABLE IF EXISTS users;"
                "CREATE TABLE users(id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER UNIQUE, f_name TEXT, is_default INTEGER);";
 
    rc = sqlite3_exec(db, sql, 0, 0, &err_msg);
    if (rc != SQLITE_OK )
    {
        printf("SQL error: %s\n", err_msg);
        sqlite3_free(err_msg);      // очищаем ресурсы
        sqlite3_close(db);
        return 1;
    }

    sqlite3_close(db);
    printf("Table created\n");
    return 0;
}

int add_user_data(const long long int user_id, const char *name){

    sqlite3 *db; 
    sqlite3_stmt *res;
    
    int rc = sqlite3_open(DATABASE_PATH, &db);
    if (rc != SQLITE_OK){
        sqlite3_close(db);
        return 1;
    }

    char *sql = "INSERT INTO users (user_id, f_name, is_default) VALUES (?, ?, ?);";
     
    rc = sqlite3_prepare_v2(db, sql, -1, &res, 0);

    if (rc == SQLITE_OK) 
    {
        sqlite3_bind_int64(res, 1, user_id);
        sqlite3_bind_text(res, 2, name, -1, SQLITE_STATIC);
        sqlite3_bind_int(res, 3, 0);
        
        int step = sqlite3_step(res);
        if (step == SQLITE_DONE) 
        {
            printf("data inserted\n");
        }
        else {
            fprintf(stderr, "Error: %s\n", sqlite3_errmsg(db));
            sqlite3_finalize(res);
            sqlite3_close(db);
            return 1;
        }
    } 
    else
    {
        fprintf(stderr, "Error: %s\n", sqlite3_errmsg(db));
        sqlite3_finalize(res);
        sqlite3_close(db);
        return 1;
    }
    sqlite3_finalize(res);
    sqlite3_close(db);

    return 0;
}

int print_all_users(){
    sqlite3 *db;
    sqlite3_stmt *res;
    
    int rc = sqlite3_open(DATABASE_PATH, &db);
    if (rc != SQLITE_OK)
    {
        fprintf(stderr, "Cannot open database: %s\n", sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }
    
    const char *sql = "SELECT * FROM users;";  // определяем запрос
    
    rc = sqlite3_prepare_v2(db, sql, -1, &res, 0);
    
    if (rc == SQLITE_OK) 
    {
        const char *pipe = "│";

        int line_length = 95;
        int separator_pos[] = {8, 28};

        char top_line[UNICODE_LINE_LEN(line_length)];
        top_line[0] = '\0';
        char mid_line[UNICODE_LINE_LEN(line_length)];
        mid_line[0] = '\0';
        char bot_line[UNICODE_LINE_LEN(line_length)];
        bot_line[0] = '\0';

        generate_line(top_line, TOP_LINE, line_length, separator_pos, 2);
        generate_line(mid_line, MID_LINE, line_length, separator_pos, 2);
        generate_line(bot_line, BOT_LINE, line_length, separator_pos, 2);

        printf("%s\n", top_line);
        printf("%s   Id %s      User ID      %s First Name                                                       %s\n", pipe, pipe, pipe, pipe);
        printf("%s\n", mid_line);

        while (sqlite3_step(res) == SQLITE_ROW) 
        {
            const int id = sqlite3_column_int(res, 0);
            const long long int user_id= sqlite3_column_int64(res, 1);
            const unsigned char *f_name = sqlite3_column_text(res, 2);
            const char is_default = (sqlite3_column_int(res, 3) == 1) ? '*' : ' ' ;
            printf("%s %c %2d %s %17lli %s %-64s %s\n", pipe, is_default, id, pipe, user_id, pipe, f_name, pipe);
        }
        printf("%s\n", bot_line);
    }
    else
    {
        fprintf(stderr, "Failed to prepare statement: %s\n", sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }
    
    sqlite3_finalize(res);
    sqlite3_close(db);
    return 0;
}

int delete_user_from_db(const UserSearch *parameter){
    sqlite3 *db; 
    sqlite3_stmt *res;
    
    int rc = sqlite3_open(DATABASE_PATH, &db);
    if (rc != SQLITE_OK){
        sqlite3_close(db);
        return 1;
    }

    const char *sql =   (parameter->type == USER_SEARCH_BY_ID) ? "DELETE FROM users WHERE id = ?;\0" :
                        (parameter->type == USER_SEARCH_BY_USER_ID) ? "DELETE FROM users WHERE user_id = ?;\0" :
                        (parameter->type == USER_SEARCH_BY_NAME) ? "DELETE FROM users WHERE f_name = ?;\0" : NULL;
     
    rc = sqlite3_prepare_v2(db, sql, -1, &res, 0);

    if (rc == SQLITE_OK) 
    {
        if (parameter->type == USER_SEARCH_BY_ID){
            sqlite3_bind_int(res, 1, parameter->value.id);
        } else if (parameter->type == USER_SEARCH_BY_NAME) {
            sqlite3_bind_text(res, 1, parameter->value.name, -1, SQLITE_STATIC);
        } else if (parameter->type == USER_SEARCH_BY_USER_ID) {
            sqlite3_bind_int64(res, 1, (parameter->value.user_id));
        }

        int step = sqlite3_step(res);
        if (step == SQLITE_DONE) 
        {
            printf("User deleted\n");
        }
        else {
            fprintf(stderr, "Error: %s\n", sqlite3_errmsg(db));
            sqlite3_finalize(res);
            sqlite3_close(db);
            return 1;
        }
    } 
    else
    {
        fprintf(stderr, "Error: %s\n", sqlite3_errmsg(db));
        sqlite3_finalize(res);
        sqlite3_close(db);
        return 1;
    }
    sqlite3_finalize(res);
    sqlite3_close(db);

    return 0;
}

int select_user(const UserSearch *parameter, User *user){
    if(parameter == NULL){
        fprintf(stderr, "Error, parameter is NULL");
    }
    sqlite3 *db;
    sqlite3_stmt *res;

    int rc  = sqlite3_open(DATABASE_PATH, &db);
    if (rc != SQLITE_OK)
    {
        sqlite3_close(db);
        return 1;
    }

    const char *sql =   (parameter->type == USER_SEARCH_BY_ID) ? "SELECT * FROM users WHERE id = ?;" :
                        (parameter->type == USER_SEARCH_BY_NAME) ? "SELECT * FROM users WHERE f_name = ?;" :
                        (parameter->type == USER_SEARCH_BY_USER_ID) ? "SELECT * FROM users WHERE user_id = ?;" : NULL;

    if (!sql) {
        sqlite3_close(db);
        return 1;
    }
    
    rc = sqlite3_prepare_v2(db, sql, -1, &res, 0);
     
    if (rc == SQLITE_OK) 
    {
        if (parameter->type == USER_SEARCH_BY_ID){
            sqlite3_bind_int(res, 1, parameter->value.id);
        } else if (parameter->type == USER_SEARCH_BY_NAME) {
            sqlite3_bind_text(res, 1, parameter->value.name, -1, SQLITE_STATIC);
        } else if (parameter->type == USER_SEARCH_BY_USER_ID) {
            sqlite3_bind_int64(res, 1, (parameter->value.user_id));
        }

        if (sqlite3_step(res) == SQLITE_ROW) {
            user->user_id = sqlite3_column_int64(res, 1);
            snprintf(user->name, sizeof(user->name), "%s", sqlite3_column_text(res, 2));
        }
    }
    else
    {
        fprintf(stderr, "Error: %s\n", sqlite3_errmsg(db));
        sqlite3_finalize(res);
        sqlite3_close(db);
        return 1;
    }
    sqlite3_finalize(res);
    sqlite3_close(db);
    return 0;
}


int set_default_user(){
    print_all_users();
    
    return 0;
}


static void generate_line(char *str, const LINE_TYPE line_type, const int length, const int *positions, size_t pos_count){

    const char *top_left_corn       = "╭";
    const char *bottom_left_corn    = "╰";
    const char *top_right_corn      = "╮";
    const char *bottom_right_corn   = "╯";
    const char *top_t               = "┬";
    const char *left_t              = "├";
    const char *cross               = "┼";
    const char *right_t             = "┤";
    const char *bottom_t            = "┴";
    const char *line                = "─";

    for (int i = 0; i < length; i++) {
        strcat(str, line);
    }
    str[UNICODE_LINE_LEN(length) - 1] = '\0';

    const char *left;
    const char *middle;
    const char *right;

    if(line_type==TOP_LINE){
        left   = top_left_corn;
        middle = top_t;
        right  = top_right_corn;
    } else if (line_type == MID_LINE) {
        left   = left_t;
        middle = cross;
        right  = right_t;
    } else if (line_type == BOT_LINE) {
        left   = bottom_left_corn;
        middle = bottom_t;
        right  = bottom_right_corn;
    }

    memcpy(str, left, strlen(left));
    if (length > 0) {
        memcpy(str + (length - 1) * strlen(right), right, strlen(right));
    }

    if (positions != NULL) {
        for (size_t i = 0; i < pos_count; i++) {
            int pos = positions[i] - 1;
            if (pos >= 0 && pos < length  * strlen(middle)) {
                memcpy(str + pos * 3, middle, 3);
            }
        }
    }
}