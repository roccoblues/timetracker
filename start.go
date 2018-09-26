package main

import (
	"time"
)

const startDescription = "Start a new timetracking interval"

type startCmd struct {
}

func (cmd *startCmd) Name() string        { return "start" }
func (cmd *startCmd) Description() string { return startDescription }
func (cmd *startCmd) Default() bool       { return false }

func (cmd *startCmd) Run(c *config) error {
	ts, err := loadTimeSheet(c.path)
	if err != nil {
		return err
	}

	if err := ts.Start(time.Now()); err != nil {
		return err
	}

	if err := ts.Save(c.path); err != nil {
		return err
	}

	ts.Print(c.out, c.roundTo)

	return nil
}
