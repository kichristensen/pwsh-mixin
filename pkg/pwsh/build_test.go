package pwsh

import (
	"bytes"
	"context"
	"os"
	"testing"

	"get.porter.sh/porter/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestMixin_Build(t *testing.T) {
	testcases := []struct {
		name           string
		inputFile      string
		wantOutputFile string
	}{
		{name: "build with config", inputFile: "testdata/build-input-with-config.yaml", wantOutputFile: "testdata/build-with-config.txt"},
		{name: "build with config wihtout psresource", inputFile: "testdata/build-input-with-config-without-psresource.yaml", wantOutputFile: "testdata/build-with-config-without-psresource.txt"},
		{name: "build without config", inputFile: "testdata/build-input-without-config.yaml", wantOutputFile: "testdata/build-without-config.txt"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := os.ReadFile(tc.inputFile)
			require.NoError(t, err)

			m := NewTestMixin(t)
			m.In = bytes.NewReader(b)

			err = m.Build(context.Background())
			require.NoError(t, err)

			test.CompareGoldenFile(t, tc.wantOutputFile, m.TestContext.GetOutput())
		})
	}
}
