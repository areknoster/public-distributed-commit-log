package validator

type producer struct {
	id string
	// name      string
	// email     string
	publicKey string
}

var producers = []producer{
	{
		id: "arek-noster-manual-hygu9uhib",
		// name:      "Arek Noster manual testing",
		// email:     "arkadiusz.noster@gmail.com",
		publicKey: "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEA4hYP2u2EGdkHHVY0GdvVwylab9hNLiBH3VUFaiJQsg8=\n-----END PUBLIC KEY-----",
	}, {
		id: "jmichalak",
		// name:      "Arek Noster manual testing",
		// email:     "arkadiusz.noster@gmail.com",
		publicKey: `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAEhurys8bur3/TSHoqr/vkRZVSgtiewzqMnDRLj/EGj8=
-----END PUBLIC KEY-----`,
	},
}
