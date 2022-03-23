```
 ______
|  ___ \            _
| |___) |_   _  ___| |_ _   _ ____   ___  ____  
|  __  /| | | |/___)  _) | | |    \ / _ \|  _ \
| |  \ \\ |_| |___ | |_| |_| | | | | |_| | | | |
|_|   \_|\____(___/ \___)__  |_|_|_|\___/|_| |_|
                       (____/   & a bunch of other languages
```

The backend for Rustymon. It provides the API for the android app as well as a CLI tool. 

## Installation

### Requirements
As of now, only `go` in the version 1.18 and `make` are requirements.

### Build
To build the project, simply run:

```bash
make && sudo make install
```
This will compile the project, move the binaries to `/usr/bin/`, create `/etc/rustymon-server/` with its 
corresponding example configuration files as well as install the systemd unit files.

### Uninstall
To uninstall, run:
```bash
make uninstall
```

## Configuration
To start the configuration, copy `/etc/rustymon-server/example.config.toml` to `/etc/rustymon-server/config.toml`.
Edit the file to match your desired configuration. At least the database section must be configured.

## CLI Usage
This project also comes with a CLI tool named `rustymon`.

## API Usage

### Account registration
- Method: `POST`
- Endpoint: `/register`

**Body**:
```json
{
  "username": "",
  "password": "",
  "email": "",
  "trainer_name": ""
}
```

### Login
- Method: `POST`
- Endpoint: `/login`

This endpoint is used to log a user in. If the request succeeds, a session cookie will be attached.
Attach it on all endpoints that require authentication.

**Body**:
```json
{
  "username": "",
  "password": ""
}
```

### Logout
- Method: `GET` or `POST`
- Endpoint: `/logout`

This endpoint is used to log out a user. It also invalidates any session cookies sent to the server.

### Serverinfo
- Method: `GET`
- Endpoint: `/serverinfo`

**Body**:
```json
{
  "version": 0,
  "registration_disabled": false
}
```
