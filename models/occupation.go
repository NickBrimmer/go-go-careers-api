package models

type Occupation struct {
	ID             string `json:"id"`
	SocID          string `json:"soc_id"`
	SocTitle       string `json:"soc_title"`
	Title          string `json:"title"`
	SingularTitle  string `json:"singular_title"`
	Description    string `json:"description"`
	TypicalEdLevel string `json:"typical_ed_level"`
}
