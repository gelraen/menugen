package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type localFile struct {
	filename string
	content  *dataDump
}

func NewLocalProvider(filename string) (Provider, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data := &dataDump{}
	err = json.Unmarshal(b, data)
	if err != nil {
		return nil, err
	}
	return &localFile{
		filename: filename,
		content:  data,
	}, nil
}

func (d *localFile) Dishes() ([]Dish, error) {
	return d.content.Dishes, nil
}

func (d *localFile) update(newData *dataDump) error {
	// 1) serialize and write newData to a temporary file
	f, err := ioutil.TempFile(filepath.Split(d.filename))
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp) // this will fail in the case of success
	err = json.NewEncoder(f).Encode(newData)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	// 2) rename temporary file to filename
	err = os.Rename(tmp, d.filename)
	if err != nil {
		// POSIX says that the filesystem should be left intact in this case.
		// "Upon successful completion, the rename() function shall return 0.
		// Otherwise, it shall return -1, errno shall be set to indicate the
		// error, and neither the file named by old nor the file named by new
		// shall be changed or created."
		// http://pubs.opengroup.org/onlinepubs/9699919799/functions/rename.html
		return err
	}

	// 3) replace d.content with dataDump
	d.content = newData // Make a deep copy?
	return nil
}

func (d *localFile) Restore(b []byte) error {
	dump := &dataDump{}
	err := json.Unmarshal(b, dump)
	if err != nil {
		return err
	}
	if err := d.update(dump); err != nil {
		return err
	}
	return nil
}
