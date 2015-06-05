package sus

import(
	`os`
	`io/ioutil`
	`golang.org/x/net/context`
)

// Creates and configures a store that stores entities by converting them to and from json []byte data and keeps them in the local file system.
func NewJsonFileStore(storeDir string, idf IdFactory, vf VersionFactory) (Store, error) {
	return NewFileStore(storeDir, `json`, jsonMarshaler, jsonUnmarshaler, idf, vf)
}

// Creates and configures a store that stores entities by converting them to and from []byte and keeps them in the local file system.
func NewFileStore(storeDir string, fileExt string, m Marshaler, un Unmarshaler, idf IdFactory, vf VersionFactory) (Store, error) {
	err := os.MkdirAll(storeDir, os.ModeDir)

	if err == nil {

		getFileName := func(id string) string {
			return storeDir + `/` + id + `.` + fileExt
		}

		get := func(ctx context.Context, id string) ([]byte, error) {
			fn := getFileName(id)
			if _, err := os.Stat(fn); err != nil {
				if os.IsNotExist(err) {
					err = EntityDoesNotExist
				}
				return nil, err
			}
			return ioutil.ReadFile(fn)
		}

		put := func(ctx context.Context, id string, d []byte) error {
			return ioutil.WriteFile(getFileName(id), d, os.ModeAppend)
		}

		del := func(ctx context.Context, id string) error {
			return os.Remove(getFileName(id))
		}

		return NewMutexByteStore(get, put, del, m, un, idf, vf), nil
	}

	return nil, err
}