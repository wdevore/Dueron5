**Usage**



**Install**

1) Create 2GB Ramdisk for compilation. 
    ```
    sudo mkdir -p /media/RAMDisk
    sudo mount -t tmpfs -o size=2048M tmpfs /media/RAMDisk/
    ```
2) Install VSCode and Go, Go outliner extensions
    ```
    "workbench.colorTheme": "Solarized Light",
    "editor.fontSize": 16,
    "terminal.integrated.fontSize": 16,
    "window.zoomLevel": 1,
    "editor.fontFamily": "'monospace', 'Courier New', 'Droid Sans Mono', 'Droid Sans Fallback'"
    ```
3) Download **Go** and extract
4) mv *go* folder to **/usr/local**
5) add to *.bashrc*:
   ```
   export PATH="$PATH:/usr/local/go/bin"
   export GOTMPDIR="/media/RAMDisk"
   export GOPATH="$HOME/Documents/Development/Go"
   ```
7) source *.bashrc*

For Linux:

*Source Code*
1) Install google cloud SDK
2) mkdir directory:
    ```
    $HOME/Documents/Development/Go/src/github.com/wdevore
    ```
if you change the Go code to reference somewhere other than *github.com* then make sure you change the directory to. At the moment the code expects *github.com* to be in the import path.

*Dependencies*
1) sudo apt-get install libsdl2-dev
2) sudo apt-get install libsdl2-ttf-dev
3) add to .bashrc (*optional if compile fails*):
    ```
    export CGO_CFLAGS="-g -O2 -I/usr/include/SDL"
    export CGO_LDFLAGS="-g -O2 -L/usr/lib/x86_64-linux-gnu"
    ```
4) go get -u "github.com/veandco/go-sdl2"
5) go get -u "github.com/emirpasic/gods"
6) go get -u "github.com/fogleman/gg"
7) go get -u "github.com/wcharczuk/go-chart"
8) go get -u github.com/derekparker/delve/cmd/dlv

*Notes*

Originally I had an issue with "**RenderUTF8BlendedWrapped**" because I had installed *libsdl-tff2.0-dev* instead of *libsdl2-tff-dev* so I had commented out the method in **sdl_ttf.go**.