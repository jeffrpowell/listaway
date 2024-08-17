package random

const (
	// DefaultN is the default value of n.
	DefaultN = 6
)

const (
	// CharSetAlphabeticLower are literally just valid alphabetic lowercase printable ASCII chars.
	CharSetAlphabeticLower = "abcdefghijklmnopqrstuvwxyz"

	// CharSetAlphabeticUpper are literally just valid alphabetic uppercase printable ASCII chars.
	CharSetAlphabeticUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// CharSetAlphabetic are literally just valid alphabetic printable ASCII chars.
	CharSetAlphabetic = CharSetAlphabeticLower + CharSetAlphabeticUpper

	// CharSetNumeric are literally just valid numeric chars.
	CharSetNumeric = "0123456789"

	// CharSetAlphaNumeric are literally just valid alphanumeric printable ASCII chars.
	CharSetAlphaNumeric = CharSetAlphabetic + CharSetNumeric

	// CharSetUnambiguousUpper are a set of unambiguous uppercase characters.
	CharSetUnambiguousUpper = "ABCDEFGHJKLMNPQRTUVWYXZ2346789"
)
