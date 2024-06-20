package scheduler

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"uw.cse374.go/scheduler/internal/config"
)

func TestScheduler(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected string
	}{
		{
			name: "single job",
			config: &config.Config{
				Jobs: []config.Job{
					{Name: "echo"},
				},
			},
			expected: "echo\n",
		},
		{
			name: "multiple",
			config: &config.Config{
				Jobs: []config.Job{
					{Name: "ls"},
					{Name: "cd"},
					{Name: "echo"},
				},
			},
			expected: "ls\ncd\necho\n",
		},
		{
			name: "dependencies",
			config: &config.Config{
				Jobs: []config.Job{
					{Name: "ls"},
					{Name: "cd"},
					{Name: "echo", DependsOn: []string{"ls"}},
				},
			},
			expected: "ls\ncd\necho\n",
		},
	}

	for _, cases := range tests {
		t.Run(cases.name, func(t *testing.T) {
			var buffer bytes.Buffer
			scheduler := New(cases.config)
			err := scheduler.Run(&buffer, false) // Pass an io.Writer and debug flag
			require.NoError(t, err)
			assert.Equal(t, cases.expected, buffer.String())
		})
	}

	t.Run("multiple", func(t *testing.T) {
		config, err := config.ReadConfig("../testdata/config/multiple.yaml")
		require.NoError(t, err)

		var buffer bytes.Buffer
		scheduler := New(config)
		err = scheduler.Run(&buffer, false)
		require.NoError(t, err)

		assert.Contains(t, buffer.String(), "echo\n")
		assert.Contains(t, buffer.String(), "ls\n")
		assert.Contains(t, buffer.String(), "cd\n")
	})

	t.Run("dependencies", func(t *testing.T) {
		config, err := config.ReadConfig("../testdata/config/dependencies.yaml")
		require.NoError(t, err)

		var buffer bytes.Buffer
		scheduler := New(config)
		err = scheduler.Run(&buffer, false)
		require.NoError(t, err)

		expected := "hello\nworld\nwhat\nis\nyour\nname?\ndone\n"
		assert.Equal(t, expected, buffer.String())
	})

	t.Run("duplicate", func(t *testing.T) {
		_, err := config.ReadConfig("../testdata/config/duplicate.yaml")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "scheduler configuration contains duplicate job")
	})

	t.Run("graph example", func(t *testing.T) {
		config, err := config.ReadConfig("../testdata/config/graphExample.yaml")
		require.NoError(t, err)

		var buffer bytes.Buffer
		scheduler := New(config)
		err = scheduler.Run(&buffer, false)
		require.NoError(t, err)

		expected := "foo\nbar\nbaz\n"
		assert.Equal(t, expected, buffer.String())
	})
}
