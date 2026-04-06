# Build Libs
FROM gcc:trixie AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    libcurl4-openssl-dev libjson-c-dev cmake binutils make\
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src
RUN git clone https://github.com/smartnode/telebot.git -b v9.4

WORKDIR /usr/src/telebot/build

RUN ["sh", "-c", "cmake .. && make"]

# Build workspace
FROM gcc:trixie

# Set timezone
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

RUN apt-get update && apt-get install -y --no-install-recommends \
    git curl locales iputils-ping
RUN locale-gen en_US.UTF-8 ru_RU.UTF-8

# Install all the toolchain dependencies for container
RUN apt-get install -y --no-install-recommends \
    gdb file git curl cmake libsqlite3-dev libjson-c-dev\
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/telebot/build/libtelebot_static.a /usr/lib/
COPY --from=builder /usr/src/telebot/include/* /usr/include/

# Add user
ENV USER_NAME=user
RUN adduser $USER_NAME; \
usermod -aG sudo $USER_NAME; \
echo "$USER_NAME:password" | chpasswd
USER $USER_NAME
ENV HOME=/home/$USER_NAME
WORKDIR $HOME

RUN echo "alias run='cmake -DDEBUG=ON .. > /dev/null&&make > /dev/null && ./main'" >> ~/.bashrc
