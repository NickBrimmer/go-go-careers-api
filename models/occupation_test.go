package models

import "testing"

func TestOccupation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		occ     Occupation
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid occupation",
			occ: Occupation{
				ID:             "11-1031.00",
				SocID:          "11-1031",
				SocTitle:       "Chief Executives",
				Title:          "Chief Executive",
				SingularTitle:  "Chief Executive",
				Description:    "Determine and formulate policies and provide overall direction",
				TypicalEdLevel: "Bachelor's degree",
			},
			wantErr: false,
		},
		{
			name: "missing id",
			occ: Occupation{
				SocID:         "11-1031",
				SocTitle:      "Chief Executives",
				Title:         "Chief Executive",
				SingularTitle: "Chief Executive",
				Description:   "Test",
			},
			wantErr: true,
			errMsg:  "missing required field: id",
		},
		{
			name: "missing soc_id",
			occ: Occupation{
				ID:            "11-1031.00",
				SocTitle:      "Chief Executives",
				Title:         "Chief Executive",
				SingularTitle: "Chief Executive",
				Description:   "Test",
			},
			wantErr: true,
			errMsg:  "missing required field: soc_id",
		},
		{
			name: "id too long",
			occ: Occupation{
				ID:            "11-1031.00-extra-chars-here",
				SocID:         "11-1031",
				SocTitle:      "Chief Executives",
				Title:         "Chief Executive",
				SingularTitle: "Chief Executive",
				Description:   "Test",
			},
			wantErr: true,
			errMsg:  "id exceeds maximum length of 20 characters",
		},
		{
			name: "title too long",
			occ: Occupation{
				ID:            "11-1031.00",
				SocID:         "11-1031",
				SocTitle:      "Chief Executives",
				Title:         string(make([]byte, 300)),
				SingularTitle: "Chief Executive",
				Description:   "Test",
			},
			wantErr: true,
			errMsg:  "title exceeds maximum length of 255 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.occ.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
