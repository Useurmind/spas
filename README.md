# SPA Server (spas) - A webserver for single page applications.

The main motivation for this server is to serve single page applications without much trouble.

We focus here on SPAs that are distribute in one or several js files together with some resources.

This server will be tested with a simple spa setup and some assets. But I hope it can also serve other setups well enough.

## Usage

Options can either be passed via command line, environment variable or config file.

If the command line option is named `--xxx`:

    - the environment variable is called `SPAS_XXX`, e.g. `--address` becomes `SPAS_ADDRESS`
    - the config file options is called just `xxx`

### Command line Help

Here is the command line help:

```
Usage:

         app.exe [options]

Available options:

        --configfile    a path to a config file that contains the configuration for the spa server (default: spas.config.json)
        --address       address to listen on (default: )
        --port          port to listen on (default: 8080)
        --servefolder   the folder to serve (default: current working directory, e.g.: C:\Users\Jochen\Projekte\Projekte\spas)
        --htmlindexfile path to the root index file of the spa app (default: index.html)
        --certfilepath  path to the ssl certificate (chain) for the server (default: spas_cert.pem)
        --keyfilepath   path to the private key of the ssl certificate of the server (default: spas_key.pem)
        --forcehttp     Set this flag to force serving http (default: false)
```

## Example Config File

The config file is a json file. By default the file `spas.config.json` is read from the working directory.

```json
{
    "address": "127.0.0.1",
    "port": "8089",
    "serveFolder": "/wwwroot",
    "htmlIndexFile": "index2.html",
    "certFilePath": "spas_cert2.pem",
    "keyFilePath": "spas_key2.pem",
    "forceHTTP": true
}
```

### Docker

You can find the docker image on dockerhub: https://hub.docker.com/r/useurmind/spas

You will have to build your own docker image with the spa files or add a volume that contains the spa to serve.

Run test test app by building a docker image with it:
```
docker build -f ./Dockerfile_test . -t spas:test
docker run -p 8080:8080 spas:test
```

Run the test app by adding a mount to the folder:

```
docker run -p 8080:8080 -v <path_to_spas_git_root>/test_resources/test_app:/www -e SPAS_SERVEFOLDER=/www useurmind/spas
```

(Make sure the slashes are all backslashes for the soure folder on windows -.-)

## How it works

Basically it servers all files in the `servefolder` as is via the golang file server functionality. It a file is not found it will serve the `htmlindexfile` file.

## Use as module in your own app

You can use the spas handler in your own app. First get the module:

```
go get github.com/Useurmind/spas
```

Then import the handler package.

```golang
import (
    "github.com/Useurmind/spas/handler
)
```

And finally create a handler to serve the SPA:

```golang
spasOptions := handler.Options{
    // only these two options are actually used by the handler
    // the rest is only used to start the spa server and not relevant here
    ServeFolder: "www",
    HTMLIndexFile: "index.html",
}
spasHandler := handler.NewSPASHandler(&spasOptions)
```

The spas handler is just a normal [http.Handler](https://golang.org/pkg/net/http/#Handler) you can use in the usual places.

## Best Practices

### Do not serve root

Do not serve the linux root folder. The spa server will scan its serve folder for all(!) files and it will serve them. Therefore, handing the linux root as a serve folder is a bad idea because the spa server will serve the complete file system. Always create a new folder in which you put only files that may be served by the spa server.

### Do not serve from spas folder

Do not put your app files in the same folder as the spa server. Put them into a different folder. Else the spa server could serve itself. Which is strange.

