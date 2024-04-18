package fwutils

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_submoduleUpdateCmd = "git submodule update"
	_checkoutCmd        = "git checkout"
	_makeCmd            = "make"

	_submoduleInitArg      = "--init"
	_submoduleRecursiveArg = "--recursive"
	_makePlatformArg       = "PLATFORM="
	_makeProjectArg        = "PROJECT="

	_defaultBinaryName = "main.bin"
	_buildFolder       = "/build"
	_loggerName        = "builder"
)

// Builder takes a firmware directory, builds a binary, and outputs it
type Builder struct {
	firmwarePath string

	l *zap.Logger
}

// NewBuilder returns a builder
func NewBuilder(firmwarePath string, l *zap.Logger) *Builder {
	return &Builder{
		firmwarePath: firmwarePath,
		l:            l.Named(_loggerName),
	}
}

// Open verifies that the firmware directory is present
func (b *Builder) Open() error {
	b.l.Info("searching for directory")

	_, err := os.Stat(b.firmwarePath)
	if err != nil {
		return errors.Wrap(err, "search for directory")
	}

	b.l.Info("directory found")

	return nil
}

// Build will checkout a commit hash, build, and return the binary
func (b *Builder) Build(platform, project, commitHash string) (string, error) {

	b.l.Info("saving current directory")

	currentDir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "get pwd")
	}

	b.l.Info("changing to firmware directory")

	err = os.Chdir(b.firmwarePath)
	if err != nil {
		return "", errors.Wrap(err, "change directory")
	}

	b.l.Info("checking out commit hash")

	checkoutCmd := exec.Command(_checkoutCmd, commitHash)

	_, err = checkoutCmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "checkout commit hash")
	}

	b.l.Info("updating submodules")

	submoduleCmd := exec.Command(_submoduleUpdateCmd, _submoduleInitArg, _submoduleRecursiveArg)

	_, err = submoduleCmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "update submodules")
	}

	buildCmd := exec.Command(_makeCmd, _makePlatformArg+platform, _makeProjectArg+project)

	_, err = buildCmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "update submodules")
	}

	//binaryPath := b.firmwarePath + _buildFolder + project + platform +

	return "", nil
}
