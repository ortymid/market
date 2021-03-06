package jwt

import (
	"encoding/json"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
}

// UnmarshalJSON converts Claims.UserID to a string if it is an integer.
func (c *Claims) UnmarshalJSON(data []byte) error {
	// Standard claims first.
	sc := jwt.StandardClaims{}
	err := json.Unmarshal(data, &sc)
	if err != nil {
		return err
	}
	c.StandardClaims = sc

	// Try to unmarshal the id as a string.
	idStr := struct {
		ID string `json:"id"`
	}{}
	err = json.Unmarshal(data, &idStr)
	if err == nil {
		// Successful unmarshaling. Continuing the normal flow.
		c.UserID = idStr.ID
	} else {
		// Try to unmarshal the id as an integer.
		if err, ok := err.(*json.UnmarshalTypeError); ok && err.Field == "id" {
			idInt := struct {
				ID int64 `json:"id"`
			}{}
			err := json.Unmarshal(data, &idInt)
			// Give up trying.
			if err != nil {
				return err
			}

			// Set the id as a string.
			c.UserID = strconv.FormatInt(idInt.ID, 10)
		} else {
			// Error is not caused by id field. Returning.
			return err
		}
	}
	return nil
}
