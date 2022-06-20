# mcsvradmin
Minecraft server management and administration tool

## Building

To build, you can just run make in the local directory.

    make

On Windows, `make` is not available. Instead you can just run the following manually:

    cd static
    go run -tags=dev static_generate.go
    cd ..
    go build -o mcsvradmin.exe

## Running in dev mode

In dev mode, you have to run with the static/files directory in place (the code will automatically
pull the latest values out of there).

    go run main.go <path-to-minecraft-directory>
