package response

type Response []byte

var (
	YES = Response{1}
	NO  = Response{0}
)

var (
	SUCCESS               = Response{100}
	INVALID_REQUEST       = Response{101}
	INTERNAL_SERVER_ERROR = Response{102}
	AUTHORIZATION_FAILED  = Response{103}
	ALREADY_EXISTS        = Response{104}
	NOT_FOUND             = Response{105}
)
