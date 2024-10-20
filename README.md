# ujxl
Jpeg XL universal batch utility (wrapper for jxl official binaries https://github.com/libjxl/libjxl)

Usage: ujxl.exe [path/]pattern [destination]

* Uses .ini style config file placed next to the executable.
* File path can be a "*" pattern.
* If destination omitted writes next to the source files.
* Alway overwites destination.
* Supports multithreading by specifying maxWorker parameter.
