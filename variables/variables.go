package variables

var (
	// ID - This processor's id.
	ID int

	// N - Number of processors
	N int

	// F - Number of faulty processors
	F int

	// Clients - Size of Clients Set
	Clients int

	// Remote - If we are running locally or remotely
	Remote bool

	// DEFAULT - The default value that is used in the algorithms
	DEFAULT []byte
)

// Initialize - Variables initializer method
func Initialize(id int, n int, c int, rem int) {
	ID = id
	N = n
	F = (N - 1) / 3
	Clients = c
	if rem == 1 {
		Remote = true
	} else {
		Remote = false
	}
	DEFAULT = []byte("")
}
