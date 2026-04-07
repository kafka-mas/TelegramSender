cmake_minimum_required(VERSION 3.26)

set(PROJECT_NAME telegram-sender)
set(PROJ_DIR "${CMAKE_CURRENT_SOURCE_DIR}")
set(TEMPLATE_DIR ${CMAKE_CURRENT_SOURCE_DIR}/templates)
set(INCLUDE_DIR ${CMAKE_CURRENT_SOURCE_DIR}/include)

project(${PROJECT_NAME} VERSION 0.2.0 LANGUAGES C)

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

target_link_libraries(
            main
            OpenSSL::Crypto
            sqlite3
            telebot_static
            ${CURL_LIBRARIES}
            ${JSONC_LIBRARIES}
        )

set_target_properties(main PROPERTIES
    OUTPUT_NAME ${OUT_NAME}
    C_STANDARD 11
    C_STANDARD_REQUIRED ON
)

# ========== INSTALL ==========

install(TARGETS main
        RUNTIME DESTINATION bin          # /usr/bin/tgsend
        )

set(PACKAGE_TEMPLATES ${CMAKE_CURRENT_SOURCE_DIR}/package_templates)

# package_templates/sender.conf.template -> /etc/telegram_sender/sender.conf
install(FILES ${PACKAGE_TEMPLATES}/sender.conf.template
        DESTINATION /etc/telegram_sender/
        RENAME sender.conf
        PERMISSIONS OWNER_READ GROUP_READ WORLD_READ   # права 644
        )

install(DIRECTORY DESTINATION /var/lib/telegram_sender/)
install(CODE "execute_process(COMMAND /usr/bin/tgsend --create-db)")

# ========== CPACK .DEB PACKAGE GENERATION ==========

set(CPACK_GENERATOR "DEB")

set(CPACK_PACKAGE_NAME "${PROJECT_NAME}")
set(CPACK_PACKAGE_VERSION "${PROJECT_VERSION}")
set(CPACK_PACKAGE_DESCRIPTION_SUMMARY "Telegram message sender daemon")
set(CPACK_PACKAGE_VENDOR "Kafka_mas")
set(CPACK_PACKAGE_CONTACT "kafka_mas@disroo.org")

set(CPACK_DEBIAN_PACKAGE_DEPENDS "libssl3, libsqlite3-0, libcurl4, libjson-c5")
# set(CPACK_RPM_PACKAGE_REQUIRES "openssl-libs, sqlite, libcurl, json-c")

set(CPACK_SET_DESTDIR ON)
set(CPACK_INSTALL_PREFIX "/usr")

set(CPACK_DEBIAN_PACKAGE_CONTROL_EXTRA
    "${PACKAGE_TEMPLATES}/postinst"
    )

include(CPack)