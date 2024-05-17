package credentials

func WriteUserCredentials(credentials UserCredentials) error {
	return nil
}

func ReadUserCredentials() (UserCredentials, error) {
	return UserCredentials{}, nil
}

func WriteSpartanTokenCredentials(credentials SpartanTokenCredentials) error {
	return nil
}

func ReadSpartanTokenCredentials() (SpartanTokenCredentials, error) {
	return SpartanTokenCredentials{}, nil
}