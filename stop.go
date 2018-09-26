package main

import (
	"time"
)

const stopDescription = "Stop the current timetracking interval"

type stopCmd struct {
}

func (cmd *stopCmd) Name() string        { return "stop" }
func (cmd *stopCmd) Description() string { return stopDescription }
func (cmd *stopCmd) Default() bool       { return false }

func (cmd *stopCmd) Run(c *config) error {
	ts, err := loadTimeSheet(c.path)
	if err != nil {
		return err
	}

	if err := ts.End(time.Now()); err != nil {
		return err
	}

	if err := ts.Save(c.path); err != nil {
		return err
	}

	ts.Print(c.out, c.roundTo)

	return nil
}
