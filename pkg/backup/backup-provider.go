
package backup

type BackupProvider interface {

	Restore() error

}