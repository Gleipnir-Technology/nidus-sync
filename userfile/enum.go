package userfile

type Collection int

const (
	CollectionAudioRaw Collection = iota
	CollectionAudioNormalized
	CollectionAudioTranscoded
	CollectionCSV
	CollectionImageRaw
	CollectionLogo
	CollectionPublicImage
)
