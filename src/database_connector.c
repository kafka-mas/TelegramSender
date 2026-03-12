#include <database_connector.h>

#include <stdio.h>

#include <stdio.h>
#include <sqlite3.h>
 
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
                "CREATE TABLE users(id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, f_name TEXT);";
 
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

int add_user_data(long long int id, char *name){

    sqlite3 *db; 
    sqlite3_stmt *res;
    
    int rc = sqlite3_open(DATABASE_PATH, &db);
    if (rc != SQLITE_OK){
        sqlite3_close(db);
        return 1;
    }

    char *sql = "INSERT INTO users (user_id, f_name) VALUES (?, ?)";
     
    rc = sqlite3_prepare_v2(db, sql, -1, &res, 0);

    if (rc == SQLITE_OK) 
    {
        sqlite3_bind_int(res, 1, id);
        sqlite3_bind_text(res, 2, name, -1, SQLITE_STATIC);
        
        int step = sqlite3_step(res);
        if (step == SQLITE_DONE) 
        {
            printf("data inserted\n");
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

int get_all_users(){
    sqlite3 *db;
    sqlite3_stmt *res;
    
    int rc = sqlite3_open(DATABASE_PATH, &db);
    if (rc != SQLITE_OK)
    {
        fprintf(stderr, "Cannot open database: %s\n", sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }
    
    const char *sql = "SELECT * FROM users";  // определяем запрос
    
    rc = sqlite3_prepare_v2(db, sql, -1, &res, 0);
    
    if (rc == SQLITE_OK) 
    {
        // Здесь не требуется привязка параметров, так как запрос без параметров
        // Перебираем все строки результата
        while (sqlite3_step(res) == SQLITE_ROW) 
        {
            int id = sqlite3_column_int(res, 0);
            long long int user_id= sqlite3_column_int(res, 1);
            const unsigned char *f_name = sqlite3_column_text(res, 2);
            printf("Id: %d\tUser ID: %lli\tFirst name: %s\n", id, user_id, f_name);
            // id, user_id, f_name;
        }
    }
    else
    {
        fprintf(stderr, "Failed to prepare statement: %s\n", sqlite3_errmsg(db));
    }
    
    sqlite3_finalize(res);
    sqlite3_close(db);
    return 0;
}

int delete_user_data(){
    return 0;
}