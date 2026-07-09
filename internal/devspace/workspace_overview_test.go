package devspace

import (
	"strings"
	"testing"
)

func TestWorkspaceOverviewBuildsSavedWorkspaceFacts(t *testing.T) {
	seedWorkspaceOverview(t)

	overview, err := buildWorkspaceOverview()
	if err != nil {
		t.Fatalf("buildWorkspaceOverview: %v", err)
	}
	if overview.WorkspaceRoot == "" || overview.ManifestVersion != ManifestVersion {
		t.Fatalf("overview root/version = %+v", overview)
	}
	if overview.ThisMachine != "laptop (machine_one)" {
		t.Fatalf("this machine = %q", overview.ThisMachine)
	}
	if len(overview.Machines) != 2 || len(overview.Users) != 1 || len(overview.Teams) != 1 {
		t.Fatalf("overview saved facts = %+v", overview)
	}
	if overview.Summary.ProjectsTracked != 2 || overview.Summary.Hydrated != 1 || overview.Summary.Placeholders != 1 {
		t.Fatalf("summary = %+v", overview.Summary)
	}
}

func TestWorkspaceOverviewRedactsManifestRemote(t *testing.T) {
	seedWorkspaceOverview(t)

	overview, err := buildWorkspaceOverview()
	if err != nil {
		t.Fatalf("buildWorkspaceOverview: %v", err)
	}
	if overview.Sync.ManifestRemote != "https://redacted@example.invalid/org/manifest.git" {
		t.Fatalf("manifest remote = %q", overview.Sync.ManifestRemote)
	}
	if strings.Contains(overview.Sync.ManifestRemote, "secret") {
		t.Fatalf("manifest remote leaked credentials: %q", overview.Sync.ManifestRemote)
	}

	var out strings.Builder
	printWorkspaceOverview(&out, overview)
	if strings.Contains(out.String(), "secret") {
		t.Fatalf("table leaked credentials:\n%s", out.String())
	}
}

func TestStatusVerboseRendersWorkspaceOverview(t *testing.T) {
	seedWorkspaceOverview(t)

	stdout, _, err := executeCommand(t, "test", "status", "--verbose")
	if err != nil {
		t.Fatalf("status --verbose error: %v", err)
	}
	for _, want := range []string{"Workspace", "Machines", "NAME", "ID", "LAST SEEN", "Users:", "Teams:", "Projects tracked: 2"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("status --verbose output missing %q:\n%s", want, stdout)
		}
	}
	if strings.Contains(stdout, "secret") {
		t.Fatalf("status --verbose output leaked credentials:\n%s", stdout)
	}
}

func seedWorkspaceOverview(t *testing.T) string {
	t.Helper()
	workspace := initCommandWorkspace(t)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	cfg.MachineID = "machine_one"
	cfg.MachineName = "laptop"
	cfg.ManifestRemote = "https://user:secret@example.invalid/org/manifest.git"
	cfg.HostedSyncEndpoint = "https://hosted.example.invalid"
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	if err := SaveManifest(workspace, Manifest{
		Version:       ManifestVersion,
		WorkspaceRoot: workspace,
		Machines: []Machine{
			{ID: "machine_one", Name: "laptop", OS: "darwin", Arch: "arm64", WorkspaceRoot: workspace, LastSeenAt: "2026-07-07T12:00:00Z"},
			{ID: "machine_two", Name: "tower", OS: "linux", Arch: "amd64", WorkspaceRoot: "/home/me/code", LastSeenAt: "2026-07-06T12:00:00Z"},
		},
		Users: []User{
			{ID: "user_one", Name: "Ada", AgeRecipient: "age1lydx38xc73yjmwfvqfpd2peulfwftx7tv7x4lw6p2gh594h2wyrqx70a4q", Status: "active", CreatedAt: "2026-07-01T00:00:00Z"},
		},
		Teams: []Team{
			{ID: "team_core", Name: "Core", Members: []TeamMember{{UserID: "user_one", Role: AccessRoleOwner, AddedAt: "2026-07-01T00:00:00Z"}}, CreatedAt: "2026-07-01T00:00:00Z"},
		},
		Projects: []Project{
			{ID: "project_api", Name: "api", Path: "apps/api", Type: ProjectTypeGit, HydrateMode: HydrateOnDemand},
			{ID: "project_docs", Name: "docs", Path: "docs", Type: ProjectTypeLocal, HydrateMode: HydrateManual},
		},
	}); err != nil {
		t.Fatalf("SaveManifest: %v", err)
	}
	if err := SaveState(State{
		MachineID:     "machine_one",
		WorkspaceRoot: workspace,
		LastSyncAt:    "2026-07-07T12:30:00Z",
		LastScanAt:    "2026-07-07T12:45:00Z",
		Projects:      map[string]ProjectState{"project_api": {Exists: true, Hydrated: true, EnvFilePresent: true}, "project_docs": {Exists: true, Placeholder: true}},
	}); err != nil {
		t.Fatalf("SaveState: %v", err)
	}
	return workspace
}
