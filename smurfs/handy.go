package smurfs

//Should helpers  &  logger for errors


func MakePointer[T any ] (t T) *T {
	return &t
}