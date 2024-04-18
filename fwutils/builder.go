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

	_defaultBinaryName      = "main"
	_defaultBinaryExtension = ".bin"
	_buildFolder            = "build"
	_loggerName             = "builder"

	_platformStm32 = "stm32f767"
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

	b.l.Info("building project")

	err = b.buildSequence(platform, project, commitHash)
	if err != nil {
		return "", errors.Wrap(err, "build project")
	}

	b.l.Info("getting binary path")

	var binaryPath string
	if platform == _platformStm32 {
		binaryPath = b.firmwarePath + b.directoryBuilder(_buildFolder, project, platform) + project + _defaultBinaryExtension
	} else {
		binaryPath = b.firmwarePath + b.directoryBuilder(_buildFolder, project, platform) + _defaultBinaryName
	}

	b.l.Info("change directory to original")

	err = os.Chdir(currentDir)
	if err != nil {
		return "", errors.Wrap(err, "change directory")
	}

	return binaryPath, nil
}

func (b *Builder) buildSequence(platform, project, commitHash string) error {
	b.l.Info("checking out commit hash")

	err := b.executeCommand(_checkoutCmd, commitHash)
	if err != nil {
		return errors.Wrap(err, "checkout commit hash")
	}

	b.l.Info("updating submodules")

	err = b.executeCommand(_submoduleUpdateCmd, _submoduleInitArg, _submoduleRecursiveArg)
	if err != nil {
		return errors.Wrap(err, "update submodules")
	}

	b.l.Info("making binary")

	err = b.executeCommand(_makeCmd, _makePlatformArg+platform, _makeProjectArg+project)
	if err != nil {
		return errors.Wrap(err, "make binary")
	}

	return nil
}

func (b *Builder) executeCommand(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)

	_, err := command.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "execute command")
	}

	return nil
}

func (b *Builder) directoryBuilder(component ...string) string {
	var outputPath string

	for _, c := range component {
		outputPath += c + "/"
	}

	return outputPath
}
