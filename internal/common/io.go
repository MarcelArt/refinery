package common

import "io"

func ResetIOCursor(s io.Seeker) error {
	_, err := s.Seek(0, io.SeekStart)
	return err
}
