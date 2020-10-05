package execute

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	v1 "github.com/priyawadhwa/prbot/pkg/config/v1"
)

// Execute executes whatever commands we need and
func Execute(cfg *v1.Config) ([]byte, error) {
	if err := runCommands(cfg.Execute.Setup); err != nil {
		return nil, errors.Wrap(err, "setup")
	}

	output, err := runCommandsCombinedOutput(cfg.Execute.Track)
	if err != nil {
		return nil, errors.Wrap(err, "track")
	}

	if err := runCommands(cfg.Execute.Cleanup); err != nil {
		return nil, errors.Wrap(err, "cleanup")
	}
	return output, nil
}

func runCommandsCombinedOutput(cmds []v1.Command) ([]byte, error) {
	var output []byte
	for _, c := range cmds {
		split := strings.Split(c.Cmd, " ")
		cmd := exec.Command(split[0], split[1:]...)
		cmd.Dir = os.ExpandEnv(c.Dir)
		log.Printf("Running %v: %v", c.Name, cmd.Args)
		o, err := cmd.Output()
		if err != nil {
			log.Printf("[%v] failed: %v\n%v", c.Name, err, string(o))
			return nil, errors.Wrapf(err, "running %v", c.Name)
		}
		output = append(output, o...)
	}
	return output, nil
}

func runCommands(cmds []v1.Command) error {
	for _, c := range cmds {
		split := strings.Split(c.Cmd, " ")
		cmd := exec.Command(split[0], split[1:]...)
		cmd.Dir = os.ExpandEnv(c.Dir)
		log.Printf("Running %v: %v", c.Name, cmd.Args)
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Printf("[%v] failed: %v\n%v", c.Name, err, string(output))
			return errors.Wrapf(err, "running %v", c.Name)
		}
	}
	return nil
}
