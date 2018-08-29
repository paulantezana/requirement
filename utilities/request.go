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

// RequestRequire use only in requirements
type RequestRequire struct {
	RequirementID uint `json:"requirement_id"`
	Type          uint `json:"query"`
}

// RequestQuotation use only in quotations
type RequestQuotation struct {
	RequirementID uint `json:"requirement_id"`
	ID            uint `json:"id"`
	Type          uint `json:"query"`
}
