package main

import (
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExe(t *testing.T) {
	err := filepath.WalkDir("/Users/hxy/Desktop/宝可梦/宝可梦叫声全集", func(path string, d fs.DirEntry, err error) error {
		//t.Log(path)
		if strings.HasSuffix(d.Name(), "mp3") {
			//t.Log(filepath.Base(d.Name()))
			out := exec.Command("ffmpeg", "-i", path, "-lavfi", "showwavespic=s=1024x512", strings.ReplaceAll(path, "mp3", "png"))
			err := out.Run()
			if err != nil {
				return err
			}
			//return out.Wait()
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
