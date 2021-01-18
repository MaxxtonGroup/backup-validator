package backup

type ResticBackupProvider struct {
}

func (p ResticBackupProvider) Restore() error {
	return nil
}

func NewResticBackupProvider() ResticBackupProvider {
	resticBackupProvider := ResticBackupProvider{}
	return resticBackupProvider
}
