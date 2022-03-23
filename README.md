# Rustymon Backend

The backend for Rustymon. It provides the API for the android app.

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

