package db

func CreateInitialEvidenceRecord(pubKey string, evId string, ext string) error {
	return db.Create(Evidence{
		EvId:       evId,
		OwnerAddr:  pubKey,
		Extension:  ext,
		CreationTx: "0x0",
		Index:      -1,
	}).Error
}

func DeleteEvidenceRecord(evId string) error {
	return db.Where(&Evidence{EvId: evId}).Delete(&Evidence{}).Error
}

func RetrieveEvidenceRecord(pubAddr string, evId string) (*Evidence, error) {
	evidence := &Evidence{}
	err := db.Where(&Evidence{
		OwnerAddr: pubAddr,
		EvId:      evId,
	}).First(&evidence).Error

	return evidence, err
}

func UpdateEvidenceRecord(evidence *Evidence) error {
	return db.Save(evidence).Error
}
