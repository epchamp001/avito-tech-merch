package password

import (
	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Log("Given the need to test HashPassword function.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen hashing a valid password.", testID)
		{
			password := "securepassword123"

			hashedPassword, err := HashPassword(password)
			if err != nil {
				t.Fatalf("\t\t%s\tTest %d:\tShould be able to hash the password: %v", "❌", testID, err)
			}
			t.Logf("\t\t%s\tTest %d:\tShould be able to hash the password.", "✅", testID)
			if hashedPassword == password {
				t.Fatalf("\t\t%s\tTest %d:\tHashed password should not be equal to the original password.", "❌", testID)
			}
			t.Logf("\t\t%s\tTest %d:\tHashed password should not be equal to the original password.", "✅", testID)

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
			if err != nil {
				t.Fatalf("\t\t%s\tTest %d:\tHashed password should match the original password: %v", "❌", testID, err)
			}
			t.Logf("\t\t%s\tTest %d:\tHashed password should match the original password.", "✅", testID)
		}

		testID++
		t.Logf("\tTest %d:\tWhen hashing an empty password.", testID)
		{
			password := ""

			hashedPassword, err := HashPassword(password)
			if err != nil {
				t.Fatalf("\t\t%s\tTest %d:\tShould be able to hash an empty password: %v", "❌", testID, err)
			}
			t.Logf("\t\t%s\tTest %d:\tShould be able to hash an empty password.", "✅", testID)

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
			if err != nil {
				t.Fatalf("\t\t%s\tTest %d:\tHashed password should match the original password: %v", "❌", testID, err)
			}
			t.Logf("\t\t%s\tTest %d:\tHashed password should match the original password.", "✅", testID)
		}
	}
}

func TestCheckPassword(t *testing.T) {
	t.Log("Given the need to test CheckPassword function.")
	{
		testCases := []struct {
			name           string
			hashedPassword string
			password       string
			expected       bool
		}{
			{
				name:           "Valid password",
				hashedPassword: "$2a$10$X8KqUkGJFIOZR3lkvDYRAOMRnsCoEP/hIiA1ymWsBzkdcji9u4/LC",
				password:       "pass123",
				expected:       true,
			},
			{
				name:           "Invalid password",
				hashedPassword: "$2a$10$abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456",
				password:       "wrongpassword",
				expected:       false,
			},
			{
				name:           "Empty password",
				hashedPassword: "$2a$10$abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456",
				password:       "",
				expected:       false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Logf("\tWhen checking password: %s.", tc.name)
				{
					result := CheckPassword(tc.hashedPassword, tc.password)

					if diff := cmp.Diff(tc.expected, result); diff != "" {
						t.Fatalf("\t\t%s\tCheckPassword result mismatch (-expected +got):\n%s", "❌", diff)
					}
					t.Logf("\t\t%s\tCheckPassword result should match the expected value.", "✅")
				}
			})
		}
	}
}
