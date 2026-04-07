# tgsend

`tgsend` is a console utility written in C that reads data from `stdin` or from a file and sends it to a specified Telegram chat via a bot.

[Russian](./README.ru.md)

---

## Usage

### Getting Started

Create a bot using [@BotFather](https://t.me/BotFather)

Specify your bot token in the configuration file (`/etc/telegram_sender/sender.conf`):

```conf
token=1234567890:ABCDEFGHJKLMNOPQRSTUVWXYZ1234567890
```

Create the database for the utility:

> [!NOTE]
> If the utility was installed as a `.deb` package, you can skip this step

```bash
tgsend --create-db
```

#### Adding a new user

To add a user, run the utility with the `-a` or `--add-user` option:

```bash
tgsend -a
```

This command generates a random key that must be sent in a dialog with the bot. Upon success, the user's name and Telegram User ID will be saved for future use.

### Command-line options

|     Option     | Short option |                          Description                                      |
|----------------|--------------|---------------------------------------------------------------------------|
| --add-user     |      -a      | Add a user to the database (database must already exist)                  |
| --create-db    |              | Create (or overwrite) the user database                                   |
| --delete-user  |              | Delete a user from the database                                           |
| --default      |      -d      | Send data to the default user                                             |
| --file         |      -f      | Specify the file to send                                                  |
| --help         |      -h      | Show help                                                                 |
| --user_id      |      -i      | Send data to a user by their Telegram ID (user must exist in the database)|
| --list-users   |      -l      | List all users                                                            |
| --set-default  |              | Set the default user                                                      |
| --user         |      -u      | Send data to a user by name (the First Name field in the table)           |
| --version      |      -v      | Show version                                                              |

### Examples

Show help:

```bash
tgsend -h
```

Send a file:

```bash
tgsend -f file-preview.md
```

Send `stdin` to the default user:

```bash
echo "Hello" | tgsend -d
```

---

## Installation

### Prebuilt packages

#### `.deb` package

Download the `.deb` file from Releases and install it:

`sudo dpkg -i tg-send_1.0.0_amd64.deb`

#### Other Linux distributions

Download the binary from Releases, then make it executable:

```bash
chmod +x tgsend && sudo mv tgsend /usr/local/bin/
```

Create the necessary directories and files:

```bash
sudo mkdir -p /var/lib/telegram_sender/
sudo mkdir -p /etc/telegram_sender/
sudo chmod -R 777 /var/lib/telegram_sender/
sudo touch /etc/telegram_sender/sender.conf
sudo chmod -R 755 /etc/telegram_sender/
```

Specify the bot token:

```bash
echo "token=<token>" | sudo tee /etc/telegram_sender/sender.conf   # Token must be specified without spaces
```

Create the database:

```bash
tgsend --create-db
```

### Building from source (using CMake)

Development and building are done inside a dev-container (see `.devcontainer` and `Dockerfile`).
