package execute

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	v1 "github.com/priyawadhwa/prbot/pkg/config/v1"
)

type Config struct {
	PRNumber string
}

func NewConfig(pr int) Config {
	return Config{
		PRNumber: fmt.Sprintf("%v", pr),
	}
}

// Execute executes whatever commands we need and
func Execute(cfg *v1.Config, ecfg Config) ([]byte, error) {
	if err := runCommands(cfg.Execute.Setup, ecfg); err != nil {
		return nil, errors.Wrap(err, "setup")
	}

	output, err := runCommandsCombinedOutput(cfg.Execute.Track, ecfg)
	if err != nil {
		return nil, errors.Wrap(err, "track")
	}

	if err := runCommands(cfg.Execute.Cleanup, ecfg); err != nil {
		return nil, errors.Wrap(err, "cleanup")
	}
	return output, nil
}

func runCommandsCombinedOutput(cmds []v1.Command, ecfg Config) ([]byte, error) {
	var output []byte
	for _, c := range cmds {
		templated, err := applyTemplate(c.Cmd, ecfg)
		if err != nil {
			return nil, errors.Wrap(err, "applying template")
		}
		split := strings.Split(templated, " ")
		cmd := exec.Command(split[0], split[1:]...)
		cmd.Dir = os.ExpandEnv(c.Dir)
		log.Printf("Running %v: %v", c.Name, cmd.Args)
		o, err := cmd.Output()
		if err != nil {
			log.Printf("[%v] failed: %v\n%v", c.Name, err, string(o))
			return nil, errors.Wrapf(err, "running command [%v]", c.Name)
		}
		output = append(output, o...)
	}
	return output, nil
}

func runCommands(cmds []v1.Command, ecfg Config) error {
	for _, c := range cmds {
		templated, err := applyTemplate(c.Cmd, ecfg)
		if err != nil {
			return errors.Wrap(err, "applying template")
		}
		split := strings.Split(templated, " ")
		cmd := exec.Command(split[0], split[1:]...)
		cmd.Dir = os.ExpandEnv(c.Dir)
		log.Printf("Running %v: %v", c.Name, cmd.Args)
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Printf("[%v] failed: %v\n%v", c.Name, err, string(output))
			return errors.Wrapf(err, "running command [%v]", c.Name)
		}
	}
	return nil
}

func applyTemplate(s string, ecfg Config) (string, error) {
	t, err := template.New("parse").Parse(s)
	if err != nil {
		return "", errors.Wrap(err, "parsing template")
	}
	buf := bytes.NewBuffer([]byte{})
	err = t.Execute(buf, ecfg)
	if err != nil {
		return "", errors.Wrap(err, "executing")
	}
	return buf.String(), nil
}
