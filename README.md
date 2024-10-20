# ujxl
Jpeg XL universal batch utility (wrapper for jxl official binaries https://github.com/libjxl/libjxl)

Usage: ujxl "[path]filename|wildcard.ext" [destination path]

* Uses .ini style config file placed next to the executable.
* In the filename can be used a "*" wildcard.
* If destination omitted writes next to the source files.
* Always overwites destination.
* Supports multithreading by specifying maxWorker parameter.
