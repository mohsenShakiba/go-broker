package socketserver

type credentialStore struct {
	config SocketServerConfig
}

func (cs *credentialStore) isValid(credential string) bool {

	for _, validCred := range cs.config.Credentials {
		if credential == validCred {
			return true
		}
	}

	return false

}