cmake_minimum_required(VERSION 3.30)
project(telegram_sender)

set(SOURCES
    src/main.c
    src/config_read.c
    src/user_mgmt.c
    src/database_connector.c
)

add_executable(main ${SOURCES})

set(PROJ_DIR "${CMAKE_CURRENT_SOURCE_DIR}")
target_compile_definitions(main PRIVATE PROJ_DIR="${PROJ_DIR}")

target_include_directories(main PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/include)

IF(DEBUG)
    MESSAGE("Build debug version")
    SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -g")
    ADD_DEFINITIONS("-DDEBUG")
ELSE(DEBUG)
    SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O2")
ENDIF(DEBUG)

find_package(OpenSSL REQUIRED)
<<<<<<< HEAD
target_link_libraries(main OpenSSL::Crypto sqlite3 telebot)
=======
target_link_libraries(main OpenSSL::Crypto telebot)
>>>>>>> 6e7ec1f (Create add_user)

set_target_properties(main PROPERTIES
    C_STANDARD 11
    C_STANDARD_REQUIRED ON
)