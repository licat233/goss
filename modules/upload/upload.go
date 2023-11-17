package _upload

import (
	"errors"

	"github.com/licat233/goss/pkg/bucket"
)

const Name = "upload file processing tools"

var CheckoutFileExts []string = nil

type Upload struct {
	Bucket *bucket.Bucket
}

func New() *Upload {
	return &Upload{
		Bucket: bucket.New(nil),
	}
}

func Run() error {
	// utils.Message("This module has not been developed!")
	// return nil
	return New().Run()
}

func (s *Upload) init() error {
	if !s.Bucket.Status {
		return errors.New("bucket not create")
	}
	return nil
}

func (s *Upload) Run() error {
	var err error
	if err = s.init(); err != nil {
		return err
	}
	if err = s.handler(); err != nil {
		return err
	}
	return nil
}

func (s *Upload) handler() error {
	return s.Bucket.UploadFiles(nil)
}
