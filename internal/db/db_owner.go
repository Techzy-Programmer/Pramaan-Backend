package db

func CreateOwner(owner *Owner) error {
	return db.Create(owner).Error
}

func RetrieveOwner(pubAddr string) (*Owner, error) {
	owner := &Owner{}
	err := db.Preload("Master").Where(&Owner{PubAddress: pubAddr}).First(&owner).Error

	return owner, err
}

func BridgeOwner(subOwnerPubAddr string, msg string, accessTx string, master string) error {
	masterOwner := &Owner{}
	err := db.Where(&Owner{PubAddress: master}).First(&masterOwner).Error
	if err != nil {
		return err
	}

	subOwner := &Owner{}
	err = db.Where(&Owner{PubAddress: subOwnerPubAddr}).First(&subOwner).Error
	if err != nil {
		return err
	}

	subOwner.MasterId = &masterOwner.PubAddress
	subOwner.Master = masterOwner
	subOwner.AccessTx = &accessTx
	subOwner.MSG = &msg

	return db.Save(subOwner).Error
}
