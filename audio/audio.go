package audio

type Audio interface {
	Play(string) error
	Stop()
}
