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
As of now, only `go` in the version 1.18 and `make` are requirements for building.

Additional, the use of `nginx` as reverse proxy is encouraged.

### Build
To build the project, simply run:

```bash
make && sudo make install
```
This will compile the project, move the binaries to `/usr/bin/`, create `/etc/rustymon-server/` with its 
corresponding example configuration files as well as install the systemd unit files.

### Additional steps
As running any binary with root context is generelly a bad idea, the service file uses systemd's 
`DynamicUser` option to set the context to a minimum. 
As only root is able to bind to any port under 1000, you should use a reverse proxy for binding to port
80 and 443. It also should serve the static files.

**Example nginx configuration**:
```nginx
server  {
    listen 80;
    listen [::]:80;

    server_name YOUR_DOMAIN_GOES_HERE;

    // Redirect to port 443 to enforce encryption
    return 302 https://$host$request_uri;
}

server {
    listen 443 ssl;
    listen [::]:443 ssl;

    // TODO: Set certificate + privkey
    ssl_certificate /PATH/TO/FULLCHAIN;
    ssl_certificate_key /PATH/TO/PRIVKEY;

    server_name YOUR_DOMAIN_GOES_HERE;

    // Static files should be served with nginx
    location /static {
        // TODO: Set path to directory above static dir
        root /PATH/TO/STATIC/FILES;
        try_files $uri $uri/ =404;
    }

    // This is the location of rustymon-server
    location / {
        // TODO: Set address to host, port combination you set in config.toml
        proxy_pass http://LOCALBIND; 
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; 
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Set all `TODO` values required in the comments, reload nginx and you're good to go.

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

### Request password reset
- Method `POST`
- Endpoint: `/requestPasswordResetByUsername` and `/requestPasswordResetByEmail`

As both, username and email must be unique, they both can be used to identify
a user. If the user was found, an email with further instructions for resetting
the password will be sent.

**Body**:
```json
{
  "username": ""
}
```

or 

```json
{
  "email": ""
}
```

### Serverinfo
- Method: `GET`
- Endpoint: `/serverinfo`

**Body**:
```json
{
  "version": 0,
  "registration_disabled": false,
  "password_reset_disabled": false
}
```
