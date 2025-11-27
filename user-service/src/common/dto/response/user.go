package response

type UserBlockInfo struct {
	IsBlocked bool   `json:"is_blocked"`
	Reason    string `json:"reason"`
}
