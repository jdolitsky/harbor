// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package action

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/goharbor/harbor/src/pkg/retention/res"
	"github.com/stretchr/testify/suite"
)

// IndexTestSuite tests the rule index
type IndexTestSuite struct {
	suite.Suite

	candidates []*res.Candidate
}

// TestIndexEntry is entry of IndexTestSuite
func TestIndexEntry(t *testing.T) {
	suite.Run(t, new(IndexTestSuite))
}

// SetupSuite ...
func (suite *IndexTestSuite) SetupSuite() {
	Register("fakeAction", newFakePerformer)

	suite.candidates = []*res.Candidate{{
		Namespace:  "library",
		Repository: "harbor",
		Kind:       "image",
		Tag:        "latest",
		PushedTime: time.Now().Unix(),
		Labels:     []string{"L1", "L2"},
	}}
}

// TestRegister tests register
func (suite *IndexTestSuite) TestGet() {
	p, err := Get("fakeAction", nil)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), p)

	results, err := p.Perform(suite.candidates)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(results))
	assert.Condition(suite.T(), func() (success bool) {
		r := results[0]
		success = r.Target != nil &&
			r.Error == nil &&
			r.Target.Repository == "harbor" &&
			r.Target.Tag == "latest"

		return
	})
}

type fakePerformer struct{}

// Perform the artifacts
func (p *fakePerformer) Perform(candidates []*res.Candidate) (results []*res.Result, err error) {
	for _, c := range candidates {
		results = append(results, &res.Result{
			Target: c,
		})
	}

	return
}

func newFakePerformer(params interface{}) Performer {
	return &fakePerformer{}
}
