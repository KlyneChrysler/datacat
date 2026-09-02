package challenge

// verifyRequest is the wire shape of a solved proof.
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
