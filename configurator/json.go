package configurator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Json struct {
}

func (j Json) UnMarshal(path string, dest any) error {
	bt, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		b, err := json.Marshal(dest)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, b, 0644)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return fmt.Errorf("ReadAndBind:%w", err)
	}
	if err = json.Unmarshal(bt, dest); err != nil {
		return fmt.Errorf("ReadAndBind:Unmarshal:%w", err)
	}
	return nil
}
