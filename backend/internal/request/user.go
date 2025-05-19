package request

type GetWorkerListReq struct {
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`

	UserEmail string `json:"-"`
}
