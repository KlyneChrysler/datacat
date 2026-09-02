package challenge

// Wire shapes only (file taxonomy, standards rule 2).

type verifyRequest struct {
	Token string `json:"token"`
	Nonce string `json:"nonce"`
}

// pageData feeds the challenge template.
type pageData struct {
	Token      string
	Difficulty int
	ReturnTo   string
}
