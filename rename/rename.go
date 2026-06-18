package rename

type Rename struct {
	ID       uint   `storm:"id,increment" json:"id"`
	DirPath  string `storm:"index" json:"dirPath"`
	NewName  string `storm:"index" json:"newName"`
	OrigName string `json:"origName"`
}

type StorageBackend interface {
	Save(r *Rename) error
	GetByDirPath(dirPath string) ([]*Rename, error)
	GetByDirAndNewName(dirPath, newName string) (*Rename, error)
	DeleteByDirPath(dirPath string) error
	DeleteByDirAndNewName(dirPath, newName string) error
}

type Storage struct {
	back StorageBackend
}

func NewStorage(back StorageBackend) *Storage {
	return &Storage{back: back}
}

func (s *Storage) AddRename(dirPath, newName, origName string) error {
	return s.back.Save(&Rename{
		DirPath:  dirPath,
		NewName:  newName,
		OrigName: origName,
	})
}

func (s *Storage) GetOriginalNames(dirPath string) map[string]string {
	records, err := s.back.GetByDirPath(dirPath)
	if err != nil {
		return map[string]string{}
	}
	result := make(map[string]string, len(records))
	for _, r := range records {
		result[r.NewName] = r.OrigName
	}
	return result
}

func (s *Storage) GetOriginalName(dirPath, newName string) string {
	r, err := s.back.GetByDirAndNewName(dirPath, newName)
	if err != nil {
		return ""
	}
	return r.OrigName
}

func (s *Storage) DeleteByDirPath(dirPath string) error {
	return s.back.DeleteByDirPath(dirPath)
}
