# trace
A minimalist and efficient logging package for Go.

## Features
* Concurrent safe
* Logging groups. You can produce, for example, an audit.log file that is separated from other logs
* Minimalist design. Two logging levels:
 	* Trace, for developers writing code
 	* Info, for operators running code
* Logging groups and levels can be enabled and disabled during runtime. Nice for simulators.
* Configurable output location per logging group