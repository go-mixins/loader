package loader

// Loader updates fields in the target object
type Loader interface {
	Load(dest interface{}) error
}

//go:generate moq -out mock/loader.go -pkg mock . Loader
