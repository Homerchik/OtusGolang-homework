package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name  string
	value EnvValue
}

func TestPrepareEnv(t *testing.T) {
	varName := "TEST"
	cases := []testCase{
		{name: "Check normal variable set", value: EnvValue{Value: "5678", NeedRemove: false}},
		{name: "Check empty variable set", value: EnvValue{Value: "", NeedRemove: false}},
		{name: "Check delete variable set", value: EnvValue{Value: "5678", NeedRemove: true}},
	}
	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			require.NoError(t, os.Setenv(varName, "123"))
			require.NoError(t, prepareEnv(Environment{varName: testcase.value}))
			if testcase.value.NeedRemove {
				require.Equal(t, "", os.Getenv(varName))
			} else {
				require.Equal(t, testcase.value.Value, os.Getenv(varName))
			}
		})
	}
}

func TestRunCmd(t *testing.T) {
	t.Run("", func(t *testing.T) {
		code := RunCmd([]string{"ls", "-la"}, Environment{})
		require.Equal(t, 0, code)
	})

	t.Run("", func(t *testing.T) {
		code := RunCmd([]string{"go", "run"}, Environment{})
		require.Equal(t, 1, code)
	})
}
