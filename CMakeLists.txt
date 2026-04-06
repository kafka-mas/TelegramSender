cmake_minimum_required(VERSION 3.30)

set(PROJECT_NAME telegram-sender)
set(PROJ_DIR "${CMAKE_CURRENT_SOURCE_DIR}")
set(TEMPLATE_DIR ${CMAKE_CURRENT_SOURCE_DIR}/templates)
set(INCLUDE_DIR ${CMAKE_CURRENT_SOURCE_DIR}/include)

project(${PROJECT_NAME} VERSION 0.1.0 LANGUAGES C)

set(SOURCES
    src/main.c
    src/config_read.c
    src/user_mgmt.c
    src/database_connector.c
)

add_executable(main ${SOURCES})

target_compile_definitions(main PRIVATE PROJ_DIR="${PROJ_DIR}")

target_include_directories(main PRIVATE ${INCLUDE_DIR})

IF(DEBUG)
    MESSAGE("Build debug version")
    SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -g")
    ADD_DEFINITIONS("-DDEBUG")
    set(OUT_NAME "main")
ELSE(DEBUG)
    SET(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -O2")
    set(OUT_NAME "tgsend")
ENDIF(DEBUG)

configure_file(
    ${TEMPLATE_DIR}/package_info.h.in
    ${INCLUDE_DIR}/package_info.h
    @ONLY
)

find_package(CURL REQUIRED)
find_package(PkgConfig)
pkg_check_modules(JSONC REQUIRED json-c)

find_package(OpenSSL REQUIRED)
target_link_libraries(main OpenSSL::Crypto sqlite3 telebot_static 
    ${CURL_LIBRARIES}
    ${JSONC_LIBRARIES})

set_target_properties(main PROPERTIES
    OUTPUT_NAME ${OUT_NAME}
    C_STANDARD 11
    C_STANDARD_REQUIRED ON
)