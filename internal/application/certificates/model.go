package certificates

type Certificate struct {
	ID            string   `json:"id" bson:"_id,omitempty"`
	ProfileID     string   `json:"profileId" bson:"profileId"`
	Name          string   `json:"name" bson:"name"`
	Issuer        string   `json:"issuer" bson:"issuer"`
	CredentialID  *string  `json:"credentialId,omitempty" bson:"credentialId,omitempty"`
	CredentialURL *string  `json:"credentialUrl,omitempty" bson:"credentialUrl,omitempty"`
	Skills        []string `json:"skills" bson:"skills"`
}

type Response struct {
	Certificates []Certificate `json:"certificates"`
}
