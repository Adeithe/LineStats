package bitwise

type Flag uint32

const (
	RESPOND_TO_COMMANDS Flag = 1 << iota
	RECORD_LOGS

	BLACKLISTED
	ADMINISTRATOR

	DONT_RESPOND_WHEN_LIVE
	BLOCK_PYRAMIDS
)

func ShouldJoinChannel(n uint32) bool {
	return Has(n, RESPOND_TO_COMMANDS) || Has(n, RECORD_LOGS)
}
