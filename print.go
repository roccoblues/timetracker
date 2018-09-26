package main

const printDescription = "Print the current timesheet"

type printCmd struct {
}

func (cmd *printCmd) Name() string        { return "print" }
func (cmd *printCmd) Description() string { return printDescription }
func (cmd *printCmd) Default() bool       { return true }

func (cmd *printCmd) Run(c *config) error {
	ts, err := loadTimeSheet(c.path)
	if err != nil {
		return err
	}

	ts.Print(c.out, c.roundTo)

	return nil
}
