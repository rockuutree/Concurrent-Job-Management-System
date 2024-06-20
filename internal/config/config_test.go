package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testdataDirectory represents the relative path to the
// test configuration definition directory.
//
// Note: the path is relative to the testdata from _this_ directory.
const testdataDirectory = "../testdata/config"

// TestConfig defines a suite of configuration tests, all with
// their own description.
func TestConfig(t *testing.T) {
	t.Run("no jobs", func(t *testing.T) {
		_, err := ReadConfig(filepath.Join(testdataDirectory, "empty.yaml"))
		require.EqualError(
			t,
			err,
			"scheduler configuration must specify at least one job",
		)
	})

	t.Run("job doesn't specify a name", func(t *testing.T) {
		_, err := ReadConfig(filepath.Join(testdataDirectory, "nameless.yaml"))
		require.EqualError(
			t,
			err,
			"all jobs must specify a name",
		)
	})

	t.Run("duplicate jobs", func(t *testing.T) {
		_, err := ReadConfig(filepath.Join(testdataDirectory, "duplicate.yaml"))
		require.EqualError(
			t,
			err,
			`scheduler configuration contains duplicate job: "echo"`,
		)
	})

	t.Run("single job", func(t *testing.T) {
		config, err := ReadConfig(filepath.Join(testdataDirectory, "single.yaml"))
		require.NoError(t, err)

		assert.NotNil(t, config)
		assert.Len(t, config.Jobs, 1)
		assert.Equal(t, "echo", config.Jobs[0].Name)
	})

	// TODO: Write more tests here!
	//
	// Hint: Check out the additional testdata files, and include a test for each
	//       one. These can be found in the internal/testdata/config directory.
	//
	// Hint: You're encouraged to use a similar testing strategy as the other tests
	//       already written here.

	t.Run("dependencies", func(t *testing.T) {
		config, err := ReadConfig(filepath.Join(testdataDirectory, "dependencies.yaml"))
		require.NoError(t, err)

		assert.NotNil(t, config)
		assert.Len(t, config.Jobs, 7)

		// Check the order and dependencies of jobs
		assert.Equal(t, "hello", config.Jobs[0].Name)
		assert.Equal(t, "world", config.Jobs[1].Name)
		assert.Equal(t, []string{"hello"}, config.Jobs[1].DependsOn)
		assert.Equal(t, "what", config.Jobs[2].Name)
		assert.Equal(t, "is", config.Jobs[3].Name)
		assert.Equal(t, []string{"what"}, config.Jobs[3].DependsOn)
		assert.Equal(t, "your", config.Jobs[4].Name)
		assert.Equal(t, []string{"is"}, config.Jobs[4].DependsOn)
		assert.Equal(t, "name?", config.Jobs[5].Name)
		assert.Equal(t, []string{"your"}, config.Jobs[5].DependsOn)
		assert.Equal(t, "done", config.Jobs[6].Name)
		assert.Empty(t, config.Jobs[6].DependsOn)
	})

	t.Run("multiple", func(t *testing.T) {
		config, err := ReadConfig(filepath.Join(testdataDirectory, "multiple.yaml"))
		require.NoError(t, err)

		assert.NotNil(t, config)
		assert.Len(t, config.Jobs, 3)
		assert.Equal(t, "echo", config.Jobs[0].Name)
		assert.Equal(t, "ls", config.Jobs[1].Name)
		assert.Equal(t, "cd", config.Jobs[2].Name)
	})

	t.Run("delay", func(t *testing.T) {
		config, err := ReadConfig(filepath.Join(testdataDirectory, "delay.yaml"))
		require.NoError(t, err)

		assert.NotNil(t, config)
		assert.Len(t, config.Jobs, 3)
		assert.Equal(t, "ls", config.Jobs[0].Name)
		assert.Equal(t, 0, config.Jobs[0].Delay)
		assert.Equal(t, "cd", config.Jobs[1].Name)
		assert.Equal(t, 2, config.Jobs[1].Delay)
		assert.Equal(t, "echo", config.Jobs[2].Name)
		assert.Equal(t, []string{"ls"}, config.Jobs[2].DependsOn)
		assert.Equal(t, 3, config.Jobs[2].Delay)
	})

	t.Run("graph", func(t *testing.T) {
		config, err := ReadConfig(filepath.Join(testdataDirectory, "graphExample.yaml"))
		require.NoError(t, err)

		assert.NotNil(t, config)
		assert.Len(t, config.Jobs, 3)
		assert.Equal(t, "foo", config.Jobs[0].Name)
		assert.Equal(t, "bar", config.Jobs[1].Name)
		assert.Equal(t, "baz", config.Jobs[2].Name)
		assert.Equal(t, []string{"foo", "bar"}, config.Jobs[2].DependsOn)
	})
}
