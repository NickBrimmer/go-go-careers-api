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
	if len(o.ID) > 20 {
		return fmt.Errorf("id exceeds maximum length of 20 characters")
	}

	if o.SocID == "" {
		return fmt.Errorf("missing required field: soc_id")
	}
	if len(o.SocID) > 20 {
		return fmt.Errorf("soc_id exceeds maximum length of 20 characters")
	}

	if o.SocTitle == "" {
		return fmt.Errorf("missing required field: soc_title")
	}
	if len(o.SocTitle) > 255 {
		return fmt.Errorf("soc_title exceeds maximum length of 255 characters")
	}

	if o.Title == "" {
		return fmt.Errorf("missing required field: title")
	}
	if len(o.Title) > 255 {
		return fmt.Errorf("title exceeds maximum length of 255 characters")
	}

	if o.SingularTitle == "" {
		return fmt.Errorf("missing required field: singular_title")
	}
	if len(o.SingularTitle) > 255 {
		return fmt.Errorf("singular_title exceeds maximum length of 255 characters")
	}

	if o.Description == "" {
		return fmt.Errorf("missing required field: description")
	}
	if len(o.Description) > 10000 {
		return fmt.Errorf("description exceeds maximum length of 10000 characters")
	}

	if len(o.TypicalEdLevel) > 100 {
		return fmt.Errorf("typical_ed_level exceeds maximum length of 100 characters")
	}

	return nil
}
