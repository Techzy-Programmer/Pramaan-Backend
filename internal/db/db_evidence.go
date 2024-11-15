package db

func CreateInitialEvidenceRecord(pubKey string, hash string, ext string, pth string) error {
	return db.Create(Evidence{
		Hash:       hash,
		OwnerAddr:  pubKey,
		Extension:  ext,
		BlobPath:   pth,
		CreationTx: "0x0",
		Index:      -1,
	}).Error
}

func DeleteEvidenceRecord(hash string) error {
	return db.Where(&Evidence{Hash: hash}).Delete(&Evidence{}).Error
}

func RetrieveEvidenceRecord(pubAddr string, hash string) (*Evidence, error) {
	evidence := &Evidence{}
	err := db.Where(&Evidence{
		OwnerAddr: pubAddr,
		Hash:      hash,
	}).First(&evidence).Error

	return evidence, err
}

func UpdateEvidenceRecord(evidence *Evidence) error {
	return db.Save(evidence).Error
}
