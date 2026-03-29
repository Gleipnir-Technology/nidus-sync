package file

type Collection int

const (
	CollectionAudioRaw Collection = iota
	CollectionAudioNormalized
	CollectionAudioTranscoded
	CollectionAvatar
	CollectionCSV
	CollectionImageRaw
	CollectionLogo
	CollectionPublicImage
)
