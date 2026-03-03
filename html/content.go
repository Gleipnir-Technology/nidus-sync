package html

type Content[T any] struct {
	C      T
	Config ContentConfig
	URL    ContentURL
}
