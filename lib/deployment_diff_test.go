package sous

import (
	"log"
	"testing"

	"github.com/nyarly/testify/assert"
	"github.com/samsalisbury/semv"
)

func TestEmptyDiff(t *testing.T) {

	intended := NewDeployments()
	existing := NewDeployments()

	dc := intended.Diff(existing)
	ds := dc.collect()

	if ds.New.Len() != 0 {
		t.Errorf("got %d new; want 0", ds.New.Len())
	}
	if ds.Gone.Len() != 0 {
		t.Errorf("got %d gone; want 0", ds.Gone.Len())
	}
	if ds.Same.Len() != 0 {
		t.Errorf("got %d same; want 0", ds.Same.Len())
	}
	if len(ds.Changed) != 0 {
		t.Errorf("got %d changed; want 0", len(ds.Changed))
	}
}

func makeDepl(repo string, num int) *Deployment {
	version, _ := semv.Parse("1.1.1-latest")
	owners := OwnerSet{}
	owners.Add("judson")
	return &Deployment{
		SourceID: SourceID{
			Location: SourceLocation{
				Repo: repo,
			},
			Version: version,
		},
		DeployConfig: DeployConfig{
			NumInstances: num,
			Env:          map[string]string{},
			Resources: map[string]string{
				"cpu":    ".1",
				"memory": "100",
				"ports":  "1",
			},
		},
		Owners: owners,
	}
}

func TestRealDiff(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lshortfile)
	assert := assert.New(t)

	intended := NewDeployments()
	existing := NewDeployments()

	repoOne := "https://github.com/opentable/one"
	repoTwo := "https://github.com/opentable/two"
	repoThree := "https://github.com/opentable/three"
	repoFour := "https://github.com/opentable/four"

	intended.MustAdd(makeDepl(repoOne, 1)) //remove

	existing.MustAdd(makeDepl(repoTwo, 1)) //same
	intended.MustAdd(makeDepl(repoTwo, 1)) //same

	existing.MustAdd(makeDepl(repoThree, 1)) //changed
	intended.MustAdd(makeDepl(repoThree, 2)) //changed

	existing.MustAdd(makeDepl(repoFour, 1)) //create

	dc := intended.Diff(existing)
	ds := dc.collect()

	if assert.Len(ds.Gone.Snapshot(), 1, "Should have one deleted item.") {
		it, _ := ds.Gone.Any(func(*Deployment) bool { return true })
		assert.Equal(string(it.SourceID.Location.Repo), repoOne)
	}

	if assert.Len(ds.Same.Snapshot(), 1, "Should have one unchanged item.") {
		it, _ := ds.Same.Any(func(*Deployment) bool { return true })
		assert.Equal(string(it.SourceID.Location.Repo), repoTwo)
	}

	if assert.Len(ds.Changed, 1, "Should have one modified item.") {
		assert.Equal(repoThree, string(ds.Changed[0].name.ManifestID.Source.Repo))
		assert.Equal(repoThree, string(ds.Changed[0].Prior.SourceID.Location.Repo))
		assert.Equal(repoThree, string(ds.Changed[0].Post.SourceID.Location.Repo))
		assert.Equal(ds.Changed[0].Post.NumInstances, 1)
		assert.Equal(ds.Changed[0].Prior.NumInstances, 2)
	}

	if assert.Equal(ds.New.Len(), 1, "Should have one added item.") {
		it, _ := ds.New.Any(func(*Deployment) bool { return true })
		assert.Equal(string(it.SourceID.Location.Repo), repoFour)
	}
}
