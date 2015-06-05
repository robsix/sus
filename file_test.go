package sus

import(
	`os`
	`fmt`
	`testing`
	`github.com/stretchr/testify/assert`
)

const(
	_TEST_DIR = `./testData`
)

func Test_NewFileStore_failure(t *testing.T){
	ffs, err := newFooFileStore(`F:\sdf.*$>?/\/\!"£$%^&)(_`, ``, nil, nil)

	assert.Nil(t, ffs, `ffs should be nil`)
	assert.NotNil(t, err, `err should not be nil`)
}

func Test_FileStore_Create(t *testing.T){
	ffs, _ := newFooFileStore(_TEST_DIR, ``, nil, nil)

	id1, f1, err1 := ffs.Create()

	assert.NotEqual(t, ``, id1, `id1 should be a non empty string`)
	assert.NotNil(t, f1, `f1 should not be nil`)
	assert.Equal(t, 0, f1.getVersion(), `f1's version should be 0`)
	assert.Nil(t, err1, `err1 should be nil`)

	id2, f2, err2 := ffs.Create()

	assert.NotEqual(t, ``, id2, `id2 should be a non empty string`)
	assert.NotEqual(t, id1, id2, `id2 should not be id1`)
	assert.NotNil(t, f2, `f2 should not be nil`)
	assert.Equal(t, 0, f2.getVersion(), `f2's version should be 0`)
	assert.True(t, f2 != f1, `f2 should not be f1`)
	assert.Nil(t, err2, `err2 should be nil`)
	os.RemoveAll(_TEST_DIR)
}

func Test_FileStore_CreateMulti_with_zero_count(t *testing.T){
	ffs, err := newFooFileStore(_TEST_DIR, ``, nil, nil)

	ids, fs, err := ffs.CreateMulti(0)

	assert.Nil(t, ids, `ids should be nil`)
	assert.Nil(t, fs, `fs should be nil`)
	assert.Nil(t, err, `err should be nil`)
}

func Test_FileStore_Read_success(t *testing.T){
	ffs, _ := newFooFileStore(_TEST_DIR, ``, nil, nil)

	id, f1, err1 := ffs.Create()

	assert.NotEqual(t, ``, id, `id should be a non empty string`)
	assert.NotNil(t, f1, `f1 should not be nil`)
	assert.Equal(t, 0, f1.getVersion(), `f1's version should be 0`)
	assert.Nil(t, err1, `err1 should be nil`)

	f2, err2 := ffs.Read(id)

	assert.NotNil(t, f2, `f2 should not be nil`)
	assert.Equal(t, f1, f2, `f2 should be f1`)
	assert.Nil(t, err2, `err2 should be nil`)
	os.RemoveAll(_TEST_DIR)
}

func Test_FileStore_ReadMulti_with_zero_count(t *testing.T){
	ffs, err := newFooFileStore(_TEST_DIR, ``, nil, nil)

	f, err := ffs.ReadMulti([]string{})

	assert.Nil(t, f, `f should be nil`)
	assert.Nil(t, err, `err1 should be nil`)
}

func Test_FileStore_Read_EntityDoesNotExist_failure(t *testing.T){
	ffs, _ := newFooFileStore(_TEST_DIR, ``, nil, nil)

	f, err := ffs.Read(``)

	assert.Nil(t, f, `f should be nil`)
	assert.Equal(t, EntityDoesNotExist, err, `err should be EntityDoesNotExist`)
	os.RemoveAll(_TEST_DIR)
}

func Test_FileStore_Update_success(t *testing.T){
	ffs, _ := newFooFileStore(_TEST_DIR, ``, nil, nil)
	id, f, err := ffs.Create()

	err = ffs.Update(id, f)

	assert.Equal(t, 1, f.getVersion(), `f's version should be 1`)
	assert.Nil(t, err, `err should be nil`)
	os.RemoveAll(_TEST_DIR)
}

func Test_FileStore_Update_EntityDoesNotExist_failure(t *testing.T){
	ffs, _ := newFooFileStore(_TEST_DIR, ``, nil, nil)
	_, f, _ := ffs.Create()

	err := ffs.Update(``, f)

	assert.Equal(t, EntityDoesNotExist, err, `err should be EntityDoesNotExist`)
	os.RemoveAll(_TEST_DIR)
}

func Test_FileStore_Update_NonsequentialUpdate_failure(t *testing.T){
	ffs, _ := newFooFileStore(_TEST_DIR, ``, nil, nil)
	id, f, _ := ffs.Create()
	f.incrementVersion()

	err := ffs.Update(id, f)

	assert.Equal(t, NonsequentialUpdate, err, `err should be NonsequentialUpdate`)
	os.RemoveAll(_TEST_DIR)
}

func Test_FileStore_UpdateMulti_LenIdsNotEqualToLenVs_failure(t *testing.T){
	ffs, err := newFooFileStore(_TEST_DIR, ``, nil, nil)

	err = ffs.UpdateMulti([]string{``}, []*foo{})

	assert.Equal(t, LenIdsNotEqualToLenVs, err, `err should be LenIdsNotEqualToLenVs`)
}

func Test_FileStore_UpdateMulti_with_zero_count(t *testing.T){
	ffs, err := newFooFileStore(_TEST_DIR, ``, nil, nil)

	err = ffs.UpdateMulti([]string{}, []*foo{})

	assert.Nil(t, err, `err should be nil`)
}

func Test_FileStore_Delete_success(t *testing.T){
	ffs, _ := newFooFileStore(_TEST_DIR, ``, nil, nil)
	id, f, err := ffs.Create()

	err = ffs.Delete(id)

	assert.Nil(t, err, `err should be nil`)

	f, err = ffs.Read(id)

	assert.Nil(t, f, `f should be nil`)
	assert.IsType(t, EntityDoesNotExist, err, `err should be EntityDoesNotExist`)
	os.RemoveAll(_TEST_DIR)
}

func Test_FileStore_DeleteMulti_with_zero_ids(t *testing.T){
	ffs, err := newFooFileStore(_TEST_DIR, ``, nil, nil)
	ids := []string{}

	err = ffs.DeleteMulti(ids)

	assert.Nil(t, err, `err should be nil`)
}

func Test_FileStore_Read_with_marshaler_error(t *testing.T){
	ffs, err := newFooFileStore(_TEST_DIR, ``, errorMarshaler, nil)
	_, _, err = ffs.Create()

	assert.Equal(t, marshalerErr, err, `err should be marshalerErr`)
	os.RemoveAll(_TEST_DIR)
}

func Test_FileStore_Read_with_unmarshaler_error(t *testing.T){
	marshaler := func(src Version)([]byte,error){return []byte{}, nil}
	ffs, err := newFooFileStore(_TEST_DIR, ``, marshaler, errorUnmarshaler)
	id, _, err := ffs.Create()

	_, err = ffs.Read(id)

	assert.Equal(t, unmarshalerErr, err, `err should be unmarshalerErr`)
	os.RemoveAll(_TEST_DIR)
}

func newFooFileStore(dir string, fileExt string, m Marshaler, un Unmarshaler) (*fooFileStore, error) {
	idSrc := 0
	var err error
	var inner VersionStore
	idf := func() string {
		idSrc++
		return fmt.Sprintf(`%d`, idSrc)
	}
	vf := func() Version {
		return &foo{NewVersion()}
	}
	if(m == nil) {
		inner, err = NewJsonFileStore(dir, idf, vf)
	} else{
		inner, err = NewFileStore(dir, fileExt, m, un, idf, vf)
	}
	if err != nil {
		return nil, err
	}
	return &fooFileStore{
		inner: inner,
	}, nil
}

type fooFileStore struct {
	inner VersionStore
}



func (ffs *fooFileStore) Create() (id string, f *foo, err error) {
	ids, fs, err := ffs.CreateMulti(1)
	if err == nil {
		id = ids[0]
		f = fs[0]
	}
	return
}

func (ffs *fooFileStore) CreateMulti(count uint) (ids []string, fs []*foo, err error) {
	ids, vs, err := ffs.inner.CreateMulti(nil, count)
	if vs != nil {
		count := len(vs)
		fs = make([]*foo, count, count)
		for i := 0; i < count; i++ {
			fs[i] = vs[i].(*foo)
		}
	}
	return
}

func (ffs *fooFileStore) Read(id string) (f *foo, err error) {
	fs, err := ffs.ReadMulti([]string{id})
	if err == nil {
		f = fs[0]
	}
	return
}

func (ffs *fooFileStore) ReadMulti(ids []string) (fs []*foo, err error) {
	vs, err := ffs.inner.ReadMulti(nil, ids)
	if vs != nil {
		count := len(vs)
		fs = make([]*foo, count, count)
		for i := 0; i < count; i++ {
			fs[i] = vs[i].(*foo)
		}
	}
	return
}

func (ffs *fooFileStore) Update(id string, f *foo) (err error) {
	return ffs.UpdateMulti([]string{id}, []*foo{f})
}

func (ffs *fooFileStore) UpdateMulti(ids []string, fs []*foo) (err error) {
	if fs != nil {
		count := len(fs)
		vs := make([]Version, count, count)
		for i := 0; i < count; i++ {
			vs[i] = Version(fs[i])
		}
		err = ffs.inner.UpdateMulti(nil, ids, vs)
	}
	return
}

func (ffs *fooFileStore) Delete(id string) (err error) {
	return ffs.DeleteMulti([]string{id})
}

func (ffs *fooFileStore) DeleteMulti(ids []string) (err error) {
	return ffs.inner.DeleteMulti(nil, ids)
}