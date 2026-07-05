package devspace

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestReconcileMergeCases(t *testing.T) {
	user := reconcileUser()
	baseProject := testMergeProject("project_app", "apps/app")
	localProject := testMergeProject("project_local", "apps/local")
	remoteProject := testMergeProject("project_remote", "apps/remote")
	baseAccess := ProjectAccess{ProjectID: baseProject.ID, UserID: user.ID, Role: AccessRoleDeveloper, AddedAt: "2026-01-01T00:00:00Z"}

	base := testMergeManifest(user, []Project{baseProject}, []ProjectAccess{baseAccess})
	emptyBase := testMergeManifest(user, nil, nil)
	local := base
	remote := base

	cases := []struct {
		name          string
		base          *Manifest
		local         Manifest
		remote        Manifest
		wantProjects  []Project
		wantOps       []ReconcileOp
		wantConflicts int
		wantTwoWay    bool
	}{
		{
			name:         "add-one-side",
			base:         &base,
			local:        local,
			remote:       testMergeManifest(user, []Project{baseProject, remoteProject}, []ProjectAccess{baseAccess}),
			wantProjects: []Project{baseProject, remoteProject},
			wantOps:      []ReconcileOp{{Action: "added", Kind: "project", Key: remoteProject.ID}},
		},
		{
			name:         "remove-one-side",
			base:         &base,
			local:        local,
			remote:       testMergeManifest(user, nil, nil),
			wantProjects: nil,
			wantOps:      []ReconcileOp{{Action: "removed", Kind: "project", Key: baseProject.ID}, {Action: "removed", Kind: "access", Key: accessKey(baseAccess)}},
		},
		{
			name:         "change-one-side",
			base:         &base,
			local:        local,
			remote:       testMergeManifest(user, []Project{withMergeRemote(baseProject, "git@example.com:team/app.git")}, []ProjectAccess{baseAccess}),
			wantProjects: []Project{withMergeRemote(baseProject, "git@example.com:team/app.git")},
			wantOps:      []ReconcileOp{{Action: "changed", Kind: "project", Key: baseProject.ID}},
		},
		{
			name:          "change-both-different",
			base:          &base,
			local:         testMergeManifest(user, []Project{withMergeName(baseProject, "ours")}, []ProjectAccess{baseAccess}),
			remote:        testMergeManifest(user, []Project{withMergeName(baseProject, "theirs")}, []ProjectAccess{baseAccess}),
			wantProjects:  []Project{withMergeName(baseProject, "ours")},
			wantConflicts: 1,
		},
		{
			name:          "add-add-different",
			base:          &emptyBase,
			local:         testMergeManifest(user, []Project{withMergeName(baseProject, "ours")}, nil),
			remote:        testMergeManifest(user, []Project{withMergeName(baseProject, "theirs")}, nil),
			wantProjects:  []Project{withMergeName(baseProject, "ours")},
			wantConflicts: 1,
		},
		{
			name:          "change-vs-remove",
			base:          &base,
			local:         testMergeManifest(user, []Project{withMergeName(baseProject, "ours")}, []ProjectAccess{baseAccess}),
			remote:        testMergeManifest(user, nil, nil),
			wantProjects:  []Project{withMergeName(baseProject, "ours")},
			wantOps:       []ReconcileOp{{Action: "removed", Kind: "access", Key: accessKey(baseAccess)}},
			wantConflicts: 1,
		},
		{
			name:         "nil-base-two-way-additions",
			local:        testMergeManifest(user, []Project{localProject}, nil),
			remote:       testMergeManifest(user, []Project{remoteProject}, nil),
			wantProjects: []Project{localProject, remoteProject},
			wantOps:      []ReconcileOp{{Action: "added", Kind: "project", Key: remoteProject.ID}},
			wantTwoWay:   true,
		},
		{
			name:          "nil-base-two-way-differences-conflict",
			local:         testMergeManifest(user, []Project{withMergeName(baseProject, "ours")}, nil),
			remote:        testMergeManifest(user, []Project{withMergeName(baseProject, "theirs")}, nil),
			wantProjects:  []Project{withMergeName(baseProject, "ours")},
			wantConflicts: 1,
			wantTwoWay:    true,
		},
		{
			name:         "idempotent",
			base:         &base,
			local:        local,
			remote:       remote,
			wantProjects: []Project{baseProject},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := reconcileManifests(tc.base, tc.local, tc.remote)
			if err != nil {
				t.Fatal(err)
			}
			if got.TwoWay != tc.wantTwoWay {
				t.Fatalf("TwoWay = %t, want %t", got.TwoWay, tc.wantTwoWay)
			}
			if len(got.Conflicts) != tc.wantConflicts {
				t.Fatalf("conflicts = %d, want %d: %+v", len(got.Conflicts), tc.wantConflicts, got.Conflicts)
			}
			if !reflect.DeepEqual(got.Merged.Projects, tc.wantProjects) {
				t.Fatalf("projects = %+v, want %+v", got.Merged.Projects, tc.wantProjects)
			}
			if !reflect.DeepEqual(got.Ops, tc.wantOps) {
				t.Fatalf("ops = %+v, want %+v", got.Ops, tc.wantOps)
			}
		})
	}
}

func TestWorkspaceReconcileNonConflicting(t *testing.T) {
	root := t.TempDir()
	remote := workspaceSyncBareRepo(t)
	workspaceA := filepath.Join(root, "machine-a", "code")
	workspaceB := filepath.Join(root, "machine-b", "code")
	homeA := filepath.Join(root, "home-a")
	homeB := filepath.Join(root, "home-b")

	t.Setenv(envHome, homeA)
	if _, err := InitWorkspace(workspaceA); err != nil {
		t.Fatal(err)
	}
	if _, err := SetManifestRemote(remote); err != nil {
		t.Fatal(err)
	}
	base := Manifest{
		Version:       ManifestVersion,
		WorkspaceRoot: workspaceA,
		Projects:      []Project{hardeningProject("apps/base", ProjectTypeLocal, "")},
	}
	if err := SaveManifest(workspaceA, base); err != nil {
		t.Fatal(err)
	}
	if _, err := PushWorkspaceManifest(); err != nil {
		t.Fatal(err)
	}

	t.Setenv(envHome, homeB)
	if _, err := InitWorkspace(workspaceB); err != nil {
		t.Fatal(err)
	}
	if _, err := SetManifestRemote(remote); err != nil {
		t.Fatal(err)
	}
	if _, err := PullWorkspaceManifest(); err != nil {
		t.Fatal(err)
	}
	localB, err := LoadManifest(workspaceB)
	if err != nil {
		t.Fatal(err)
	}
	localB.Projects = append(localB.Projects, hardeningProject("apps/remote", ProjectTypeLocal, ""))
	if err := SaveManifest(workspaceB, localB); err != nil {
		t.Fatal(err)
	}
	if _, err := PushWorkspaceManifest(); err != nil {
		t.Fatal(err)
	}

	t.Setenv(envHome, homeA)
	localA, err := LoadManifest(workspaceA)
	if err != nil {
		t.Fatal(err)
	}
	localA.Projects = append(localA.Projects, hardeningProject("apps/local", ProjectTypeLocal, ""))
	if err := SaveManifest(workspaceA, localA); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(manifestPath(workspaceA))
	if err != nil {
		t.Fatal(err)
	}

	plan, err := ReconcileWorkspaceManifest("", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Conflicts) != 0 || len(plan.Ops) != 1 {
		t.Fatalf("plan conflicts=%+v ops=%+v", plan.Conflicts, plan.Ops)
	}
	merged, err := LoadManifest(workspaceA)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"apps/base", "apps/local", "apps/remote"} {
		if _, ok := findProject(merged, path); !ok {
			t.Fatalf("missing project %s after reconcile: %+v", path, merged.Projects)
		}
	}
	backup, err := os.ReadFile(filepath.Join(homeA, "manifest-backup.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(backup, before) {
		t.Fatal("backup does not match pre-reconcile manifest")
	}
	baseSnapshot, ok, err := loadBaseManifest()
	if err != nil {
		t.Fatal(err)
	}
	if !ok || !reflect.DeepEqual(baseSnapshot.Projects, merged.Projects) {
		t.Fatalf("base snapshot = %+v ok=%t, want merged %+v", baseSnapshot.Projects, ok, merged.Projects)
	}
	second, err := ReconcileWorkspaceManifest("", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(second.Ops) != 0 || len(second.Conflicts) != 0 {
		t.Fatalf("second reconcile ops=%+v conflicts=%+v", second.Ops, second.Conflicts)
	}
}

func TestWorkspaceReconcileConflictBlocksApply(t *testing.T) {
	root := t.TempDir()
	remote := workspaceSyncBareRepo(t)
	workspaceA := filepath.Join(root, "machine-a", "code")
	workspaceB := filepath.Join(root, "machine-b", "code")
	homeA := filepath.Join(root, "home-a")
	homeB := filepath.Join(root, "home-b")

	t.Setenv(envHome, homeA)
	if _, err := InitWorkspace(workspaceA); err != nil {
		t.Fatal(err)
	}
	if _, err := SetManifestRemote(remote); err != nil {
		t.Fatal(err)
	}
	project := hardeningProject("apps/app", ProjectTypeLocal, "")
	if err := SaveManifest(workspaceA, Manifest{
		Version:       ManifestVersion,
		WorkspaceRoot: workspaceA,
		Projects:      []Project{project},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := PushWorkspaceManifest(); err != nil {
		t.Fatal(err)
	}

	t.Setenv(envHome, homeB)
	if _, err := InitWorkspace(workspaceB); err != nil {
		t.Fatal(err)
	}
	if _, err := SetManifestRemote(remote); err != nil {
		t.Fatal(err)
	}
	if _, err := PullWorkspaceManifest(); err != nil {
		t.Fatal(err)
	}
	remoteManifest, err := LoadManifest(workspaceB)
	if err != nil {
		t.Fatal(err)
	}
	remoteManifest.Projects[0].Name = "remote-app"
	if err := SaveManifest(workspaceB, remoteManifest); err != nil {
		t.Fatal(err)
	}
	if _, err := PushWorkspaceManifest(); err != nil {
		t.Fatal(err)
	}

	t.Setenv(envHome, homeA)
	localManifest, err := LoadManifest(workspaceA)
	if err != nil {
		t.Fatal(err)
	}
	localManifest.Projects[0].Name = "local-app"
	if err := SaveManifest(workspaceA, localManifest); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(manifestPath(workspaceA))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ReconcileWorkspaceManifest("", true); err == nil || !strings.Contains(err.Error(), "unresolved reconcile conflicts") {
		t.Fatalf("apply conflict error = %v", err)
	}
	after, err := os.ReadFile(manifestPath(workspaceA))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, before) {
		t.Fatal("conflicted apply changed local manifest")
	}

	if _, err := ReconcileWorkspaceManifest("remote", true); err != nil {
		t.Fatal(err)
	}
	appliedRemote, err := LoadManifest(workspaceA)
	if err != nil {
		t.Fatal(err)
	}
	if appliedRemote.Projects[0].Name != "remote-app" {
		t.Fatalf("force remote project name = %q", appliedRemote.Projects[0].Name)
	}

	appliedRemote.Projects[0].Name = "local-app"
	if err := SaveManifest(workspaceA, appliedRemote); err != nil {
		t.Fatal(err)
	}
	if err := recordBaseManifest(Manifest{
		Version:       ManifestVersion,
		WorkspaceRoot: workspaceA,
		Projects:      []Project{project},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := ReconcileWorkspaceManifest("local", true); err != nil {
		t.Fatal(err)
	}
	appliedLocal, err := LoadManifest(workspaceA)
	if err != nil {
		t.Fatal(err)
	}
	if appliedLocal.Projects[0].Name != "local-app" {
		t.Fatalf("force local project name = %q", appliedLocal.Projects[0].Name)
	}
}

func reconcileUser() User {
	return User{
		ID:           "user_1",
		Name:         "Test User",
		AgeRecipient: "age1lydx38xc73yjmwfvqfpd2peulfwftx7tv7x4lw6p2gh594h2wyrqx70a4q",
		CreatedAt:    "2026-01-01T00:00:00Z",
	}
}
