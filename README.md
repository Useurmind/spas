# SPA Server (spas) - A webserver for single page applications.

The main motivation for this server is to serve single page applications without much trouble.

We focus here on SPAs that are distribute in one or several js files together with some resources.

This server will be tested with a default webpack setup. But I hope it can also serve other setups well enough.

## Usage

```
Usage:

         spas.exe [options]

Available options:

        --configfile    a path to a config file that contains the configuration for the spa server (default: spas.config.json)
        --address       address to listen on (default: )
        --port          port to listen on (default: 8080)
        --servefolder   the folder to serve (default: current working directory, e.g.: /app)
        --htmlindexfile path to the root index file of the spa app (default: index.html)
```

## Example Config File

```json
{
    "address": "127.0.0.1",
    "port": "8089",
    "serveFolder": "/wwwroot",
    "htmlIndexFile": "index.html"
}
```

## Best Practices

### Do not serve root

Do not serve the linux root folder. The spa server will scan its serve folder for all(!) files and it will serve them. Therefore, handing the linux root as a serve folder is a bad idea because the spa server will serve the complete file system. Always create a new folder in which you put only files that may be served by the spa server.


