package models

import "fmt"

type Occupation struct {
	ID             string `json:"id"`
	SocID          string `json:"soc_id"`
	SocTitle       string `json:"soc_title"`
	Title          string `json:"title"`
	SingularTitle  string `json:"singular_title"`
	Description    string `json:"description"`
	TypicalEdLevel string `json:"typical_ed_level"`
}

func (o *Occupation) Validate() error {
	if o.ID == "" {
		return fmt.Errorf("missing required field: id")
	}
	if o.SocID == "" {
		return fmt.Errorf("missing required field: soc_id")
	}
	if o.SocTitle == "" {
		return fmt.Errorf("missing required field: soc_title")
	}
	if o.Title == "" {
		return fmt.Errorf("missing required field: title")
	}
	if o.SingularTitle == "" {
		return fmt.Errorf("missing required field: singular_title")
	}
	if o.Description == "" {
		return fmt.Errorf("missing required field: description")
	}
	return nil
}
