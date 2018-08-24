package utilities

// Request struct
// 2 = primordial
// 1 = minimal
// 0 = all
type Request struct {
	Search      string `json:"search"`
	CurrentPage uint   `json:"current_page"`
	Limit       uint   `json:"limit"`
	Type        uint   `json:"query"`
}

type RequestRequire struct {
	RequirementID uint `json:"requirement_id"`
	Type          uint `json:"query"`
}

type RequestQuotation struct {
	RequirementID uint `json:"requirement_id"`
	ID            uint `json:"id"`
	Type          uint `json:"query"`
}
