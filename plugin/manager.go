package plugin

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"

	"github.com/notaryproject/notation-go/dir"
	"github.com/notaryproject/notation-go/plugin/proto"
)

// ErrNotCompliant is returned by plugin methods when the response is not
// compliant.
var ErrNotCompliant = errors.New("plugin not compliant")

// ErrNotRegularFile is returned when the plugin file is not an regular file.
var ErrNotRegularFile = errors.New("not regular file")

// Manager manages plugins installed on the system.
type Manager interface {
	Get(ctx context.Context, name string) (Plugin, error)
	List(ctx context.Context) ([]string, error)
}

// UberManager implements Manager and supports both plugin as library as well as executable
type UberManager struct {
	cliManager *CLIManager
	libPlugins map[string]Plugin
}

// NewUberManager returns Manager which supports both plugin as library as well as executable
func NewUberManager(cxt context.Context, cliMgr *CLIManager, libs ...Plugin) (*UberManager, error) {
	libPlugins := make(map[string]Plugin, len(libs))
	for _, lib := range libs {
		metadataResp, err := lib.GetMetadata(cxt, &proto.GetMetadataRequest{})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize library plugin")
		}
		libPlugins[metadataResp.Name] = lib
	}

	return &UberManager{
		cliManager: cliMgr,
		libPlugins: libPlugins,
	}, nil
}

// Get returns a plugin by its name. If the plugin is not found, the error is of type os.ErrNotExist.
func (m *UberManager) Get(ctx context.Context, name string) (Plugin, error) {
	if p, ok := m.libPlugins[name]; ok {
		return p, nil
	}
	if m.cliManager != nil {
		return m.cliManager.Get(ctx, name)
	}

	return nil, fmt.Errorf("'%s' plugin not found", name)
}

// List produces a list of the plugin names
func (m *UberManager) List(ctx context.Context) ([]string, error) {
	list := make([]string, len(m.libPlugins))

	i := 0
	for name := range m.libPlugins {
		list[i] = name
		i++
	}

	if m.cliManager != nil {
		cliList, err := m.cliManager.List(ctx)
		if err != nil {
			return nil, err
		}
		for _, cliName := range cliList {
			list = append(list, cliName)
		}
	}

	return list, nil
}

// CLIManager implements Manager
type CLIManager struct {
	pluginFS dir.SysFS
}

// NewCLIManager returns CLIManager for named pluginFS.
func NewCLIManager(pluginFS dir.SysFS) *CLIManager {
	return &CLIManager{pluginFS: pluginFS}
}

// Get returns a plugin on the system by its name.
//
// If the plugin is not found, the error is of type os.ErrNotExist.
func (m *CLIManager) Get(ctx context.Context, name string) (Plugin, error) {
	pluginPath := path.Join(name, binName(name))
	path, err := m.pluginFS.SysPath(pluginPath)
	if err != nil {
		return nil, err
	}

	// validate and create plugin
	return NewCLIPlugin(ctx, name, path)
}

// List produces a list of the plugin names on the system.
func (m *CLIManager) List(_ context.Context) ([]string, error) {
	var plugins []string
	fs.WalkDir(m.pluginFS, ".", func(dir string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dir == "." {
			// Ignore root dir.
			return nil
		}
		typ := d.Type()
		if !typ.IsDir() || typ&fs.ModeSymlink != 0 {
			// Ignore non-directories and symlinked directories.
			return nil
		}

		// add plugin name
		plugins = append(plugins, d.Name())
		return fs.SkipDir
	})
	return plugins, nil
}
