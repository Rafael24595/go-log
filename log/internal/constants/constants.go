package constants

// DefaultPath defines the default directory where log files will be stored 
// if no specific path is provided in the FileProvider.
const DefaultPath = ".logs"

// DefaultBufferSize defines the default capacity of the internal 
// channels used by the logging engines. This helps balance memory 
// usage and throughput.
const DefaultBufferSize = 100
