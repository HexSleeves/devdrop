package devspace

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestResolveServeAddrRejectsExplicitPublicWithoutOptIn(t *testing.T) {
	_, err := resolveServeAddr("0.0.0.0:8787", true, "", false)
	if err == nil || !strings.Contains(err.Error(), "refusing to bind public address") {
		t.Fatalf("expected public-bind guard error, got %v", err)
	}
}

func TestResolveServeAddrAllowsPublicWithFlag(t *testing.T) {
	got, err := resolveServeAddr("0.0.0.0:8787", true, "", true)
	if err != nil {
		t.Fatalf("unexpected error with --allow-public-http: %v", err)
	}
	if got.addr != "0.0.0.0:8787" {
		t.Fatalf("addr = %q", got.addr)
	}
	if !got.public {
		t.Fatal("expected public=true for 0.0.0.0 bind")
	}
}

func TestResolveServeAddrAllowsPublicWithPortEnv(t *testing.T) {
	// PORT env implies a PaaS/proxy that terminates TLS upstream, so a public
	// bind is sanctioned without --allow-public-http.
	got, err := resolveServeAddr("127.0.0.1:8787", false, "8080", false)
	if err != nil {
		t.Fatalf("unexpected error with PORT env: %v", err)
	}
	if got.addr != "0.0.0.0:8080" {
		t.Fatalf("addr = %q, want 0.0.0.0:8080", got.addr)
	}
	if !got.public {
		t.Fatal("expected public=true for PORT-driven 0.0.0.0 bind")
	}
}

func TestResolveServeAddrAllowsLoopback(t *testing.T) {
	got, err := resolveServeAddr("127.0.0.1:8787", true, "", false)
	if err != nil {
		t.Fatalf("loopback bind should not require opt-in: %v", err)
	}
	if got.public {
		t.Fatal("expected public=false for loopback bind")
	}
	if got.addr != "127.0.0.1:8787" {
		t.Fatalf("addr = %q", got.addr)
	}
}

func TestResolveServeAddrRejectsInvalidAddr(t *testing.T) {
	if _, err := resolveServeAddr("not-a-valid-addr", true, "", false); err == nil {
		t.Fatal("expected error for invalid listen address")
	}
}

// executeCommand wires a fresh root command to buffer-backed stdout/stderr, sets
// the given args, executes it, and returns the captured output and the error.
// It does NOT touch the real DEVSPACE_HOME; callers set that up first.
func executeCommand(t *testing.T, version string, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	cmd := NewRootCommand(version)
	var out, errBuf bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return out.String(), errBuf.String(), err
}

func executePrivateCommand(t *testing.T, command *cobra.Command, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	root := &cobra.Command{
		Use:          "devspace",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configureStyles(cmd.OutOrStdout(), false)
		},
	}
	root.AddCommand(command)
	var out, errBuf bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errBuf)
	root.SetArgs(append([]string{command.Name()}, args...))
	err = root.Execute()
	return out.String(), errBuf.String(), err
}

func TestReleaseCommandTreeContract(t *testing.T) {
	root := NewRootCommand("v1.2.3-test")
	wantCommands := []string{
		"apply", "doctor", "env", "experimental", "hosted", "init", "plan",
		"project", "scan", "setup", "status", "sync", "ui", "watch",
	}
	var gotCommands []string
	for _, cmd := range root.Commands() {
		if !cmd.Hidden && cmd.Name() != "completion" && cmd.Name() != "help" {
			gotCommands = append(gotCommands, cmd.Name())
		}
	}
	slices.Sort(gotCommands)
	if !slices.Equal(gotCommands, wantCommands) {
		t.Fatalf("visible root commands = %v, want %v", gotCommands, wantCommands)
	}
	if len(gotCommands) > 14 {
		t.Fatalf("visible root command count = %d, want at most 14", len(gotCommands))
	}

	wantGroups := map[string]string{
		"core":         "Core Workflow",
		"management":   "Workspace Management",
		"diagnostics":  "Diagnostics and Automation",
		"experimental": "Experimental",
	}
	for _, group := range root.Groups() {
		want, ok := wantGroups[group.ID]
		if !ok {
			t.Errorf("unexpected root command group %q", group.ID)
			continue
		}
		if strings.TrimSuffix(group.Title, ":") != want {
			t.Errorf("root command group %q title = %q, want %q", group.ID, group.Title, want)
		}
		delete(wantGroups, group.ID)
	}
	if len(wantGroups) != 0 {
		t.Fatalf("missing root command groups: %v", wantGroups)
	}

	stdout, _, err := executeCommand(t, "v1.2.3-test", "--help")
	if err != nil {
		t.Fatalf("root --help error: %v", err)
	}
	for _, want := range []string{"Core Workflow", "Workspace Management", "Diagnostics and Automation", "Experimental", "devspace scan", "devspace status --verbose"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("root --help missing %q:\n%s", want, stdout)
		}
	}
	for _, removed := range []string{"workspace", "tui", "mount", "version"} {
		if slices.Contains(gotCommands, removed) {
			t.Errorf("root still exposes removed command %q", removed)
		}
	}

	stdout, _, err = executeCommand(t, "v1.2.3-test", "--version")
	if err != nil {
		t.Fatalf("root --version error: %v", err)
	}
	if !strings.Contains(stdout, "v1.2.3-test") {
		t.Fatalf("root --version output = %q", stdout)
	}
	for _, name := range []string{"sync", "experimental"} {
		stdout, _, err := executeCommand(t, "test", name)
		if err != nil {
			t.Errorf("devspace %s error: %v", name, err)
			continue
		}
		for _, want := range []string{"Usage:", "devspace " + name + " --help"} {
			if !strings.Contains(stdout, want) {
				t.Errorf("devspace %s output missing %q:\n%s", name, want, stdout)
			}
		}
	}

	for _, args := range [][]string{{"workspace"}, {"tui"}, {"mount"}, {"version"}, {"project", "status"}, {"sync", "workspace"}} {
		if _, _, err := executeCommand(t, "test", args...); err == nil || !strings.Contains(err.Error(), "unknown command") {
			t.Errorf("devspace %s error = %v, want unknown command", strings.Join(args, " "), err)
		}
	}
}

func TestVersionFlagPrintsVersion(t *testing.T) {
	stdout, _, err := executeCommand(t, "v1.2.3-test", "--version")
	if err != nil {
		t.Fatalf("--version error: %v", err)
	}
	if !strings.Contains(stdout, "v1.2.3-test") {
		t.Fatalf("--version output = %q, want configured version", stdout)
	}
}

func TestRootHelpListsSubcommands(t *testing.T) {
	// --help on root: SilenceUsage is true, but help still prints to stdout.
	stdout, _, err := executeCommand(t, "test", "--help")
	if err != nil {
		t.Fatalf("root --help error: %v", err)
	}
	for _, want := range []string{"init", "scan", "watch", "plan", "apply", "hosted", "env", "doctor"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("root --help missing subcommand %q in output:\n%s", want, stdout)
		}
	}
}

func TestHostedServeHelpDocumentsFlags(t *testing.T) {
	stdout, _, err := executeCommand(t, "test", "hosted", "serve", "--help")
	if err != nil {
		t.Fatalf("hosted serve --help error: %v", err)
	}
	for _, want := range []string{"--addr", "--token", "--trusted-proxy", "--allow-public-http"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("hosted serve --help missing flag %q in output:\n%s", want, stdout)
		}
	}
}

func TestWorkspaceRemoteSetHelpDocumentsCommitFlags(t *testing.T) {
	stdout, _, err := executePrivateCommand(t, newWorkspaceCommand(), "remote", "set", "--help")
	if err != nil {
		t.Fatalf("workspace remote set --help error: %v", err)
	}
	for _, want := range []string{"--commit-email", "--commit-name"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("workspace remote set --help missing flag %q in output:\n%s", want, stdout)
		}
	}
}

// initCommandWorkspace prepares an isolated DEVSPACE_HOME + initialized
// workspace and returns the workspace root. Used by the command tests below.
func initCommandWorkspace(t *testing.T) string {
	t.Helper()
	t.Setenv(envHome, t.TempDir())
	workspace := filepath.Join(t.TempDir(), "code")
	if _, err := InitWorkspace(workspace); err != nil {
		t.Fatalf("InitWorkspace: %v", err)
	}
	return workspace
}

func TestScanCommandReportsEmptyWorkspace(t *testing.T) {
	initCommandWorkspace(t)
	stdout, _, err := executeCommand(t, "test", "scan")
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}
	if !strings.Contains(stdout, "Found 0 projects") {
		t.Fatalf("scan output = %q, want 0 projects", stdout)
	}
}

func TestStatusCommandReportsWorkspaceHealth(t *testing.T) {
	initCommandWorkspace(t)
	stdout, _, err := executeCommand(t, "test", "status")
	if err != nil {
		t.Fatalf("status error: %v", err)
	}
	for _, want := range []string{"Projects tracked: 0", "Machine:", "Workspace:"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("status output missing %q:\n%s", want, stdout)
		}
	}
}

func TestStatusCommandVerboseShowsRedactedWorkspaceOverview(t *testing.T) {
	seedWorkspaceOverview(t)
	stdout, _, err := executeCommand(t, "test", "status", "--verbose")
	if err != nil {
		t.Fatalf("status --verbose error: %v", err)
	}
	for _, want := range []string{"Workspace", "Machines", "Manifest remote: https://redacted@example.invalid/org/manifest.git", "Projects tracked: 2"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("status --verbose missing %q:\n%s", want, stdout)
		}
	}
	if strings.Contains(stdout, "secret") {
		t.Fatalf("status --verbose leaked credentials:\n%s", stdout)
	}
}

func TestStatusCommandSelectsProject(t *testing.T) {
	seedWorkspaceOverview(t)
	stdout, _, err := executeCommand(t, "test", "status", "api")
	if err != nil {
		t.Fatalf("status api error: %v", err)
	}
	for _, want := range []string{"Project: api", "Path: apps/api", "Hydrated: true"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("status api missing %q:\n%s", want, stdout)
		}
	}
}

func TestStatusCommandProjectJSONIsClean(t *testing.T) {
	seedWorkspaceOverview(t)
	stdout, _, err := executeCommand(t, "test", "status", "api", "--json")
	if err != nil {
		t.Fatalf("status api --json error: %v", err)
	}
	var row ProjectListRow
	if err := json.Unmarshal([]byte(stdout), &row); err != nil {
		t.Fatalf("status api --json did not parse: %v\n%s", err, stdout)
	}
	if row.Project.Name != "api" || row.Project.Path != "apps/api" || !row.State.Hydrated {
		t.Fatalf("status api --json row = %+v", row)
	}
	if strings.Contains(stdout, "\x1b[") {
		t.Fatalf("status api --json contains ANSI escapes: %q", stdout)
	}
}

func TestStatusCommandRejectsInvalidCombinations(t *testing.T) {
	for _, tc := range []struct {
		args    []string
		wantErr string
	}{
		{args: []string{"status", "api", "extra"}, wantErr: "accepts at most 1 arg"},
		{args: []string{"status", "api", "--verbose"}, wantErr: "--verbose cannot be combined with a project"},
		{args: []string{"status", "--verbose", "--json"}, wantErr: "--verbose and --json are mutually exclusive"},
	} {
		if _, _, err := executeCommand(t, "test", tc.args...); err == nil || !strings.Contains(err.Error(), tc.wantErr) {
			t.Errorf("devspace %s error = %v, want %q", strings.Join(tc.args, " "), err, tc.wantErr)
		}
	}
}

func TestStatusCommandHelpDocumentsConsolidatedSurface(t *testing.T) {
	stdout, _, err := executeCommand(t, "test", "status", "--help")
	if err != nil {
		t.Fatalf("status --help error: %v", err)
	}
	for _, want := range []string{"status [project]", "--verbose", "--json", "devspace status api --json"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("status --help missing %q:\n%s", want, stdout)
		}
	}
}

func TestDoctorCommandRunsAndReports(t *testing.T) {
	initCommandWorkspace(t)
	stdout, _, err := executeCommand(t, "test", "doctor")
	if err != nil {
		t.Fatalf("doctor error: %v", err)
	}
	// doctor always prints something; just confirm it produced output.
	if strings.TrimSpace(stdout) == "" {
		t.Fatal("doctor produced no output")
	}
}

func TestProjectAddAndStatusCommands(t *testing.T) {
	initCommandWorkspace(t)

	// Track a project via the command surface.
	stdout, _, err := executeCommand(t, "test", "project", "add", "apps/app")
	if err != nil {
		t.Fatalf("project add error: %v", err)
	}
	if !strings.Contains(stdout, "Added project app at apps/app") {
		t.Fatalf("project add output = %q", stdout)
	}

	// Status for a single tracked project.
	stdout, _, err = executeCommand(t, "test", "status", "app")
	if err != nil {
		t.Fatalf("status app error: %v", err)
	}
	for _, want := range []string{"Project: app", "Path: apps/app", "Hydrated:"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("project status output missing %q:\n%s", want, stdout)
		}
	}
}

func TestProjectCommandListsTrackedProjects(t *testing.T) {
	initCommandWorkspace(t)

	stdout, _, err := executeCommand(t, "test", "project")
	if err != nil {
		t.Fatalf("project list empty error: %v", err)
	}
	if !strings.Contains(stdout, "No tracked projects.") {
		t.Fatalf("project list empty output = %q", stdout)
	}

	if _, _, err := executeCommand(t, "test", "project", "add", "apps/api"); err != nil {
		t.Fatalf("project add error: %v", err)
	}
	stdout, _, err = executeCommand(t, "test", "project")
	if err != nil {
		t.Fatalf("project list error: %v", err)
	}
	for _, want := range []string{"api", "apps/api", ProjectTypeLocal, "hydrated"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("project list output missing %q:\n%s", want, stdout)
		}
	}

	stdout, _, err = executeCommand(t, "test", "status", "api")
	if err != nil {
		t.Fatalf("status api after list wiring error: %v", err)
	}
	if !strings.Contains(stdout, "Project: api") {
		t.Fatalf("project status output = %q", stdout)
	}
}

func TestProjectCommandShowsSavedProjectStateWithoutMutating(t *testing.T) {
	workspace := initCommandWorkspace(t)
	manifest := Manifest{
		Version:       ManifestVersion,
		WorkspaceRoot: workspace,
		Projects: []Project{
			{ID: "project_api", Name: "api", Path: "apps/api", Type: ProjectTypeGit, HydrateMode: HydrateOnDemand},
			{ID: "project_docs", Name: "docs", Path: "docs/site", Type: ProjectTypeGit, HydrateMode: HydrateOnDemand},
			{ID: "project_worker", Name: "worker", Path: "services/worker", Type: ProjectTypeLocal, HydrateMode: HydrateManual},
		},
	}
	if err := SaveManifest(workspace, manifest); err != nil {
		t.Fatalf("SaveManifest: %v", err)
	}
	state := State{
		Projects: map[string]ProjectState{
			"project_api": {
				Exists: true, Placeholder: true, CurrentBranch: "main", LastCheckedAt: nowRFC3339(),
			},
			"project_docs": {
				Missing: true, LastCheckedAt: nowRFC3339(),
			},
			"project_worker": {
				Exists: true, Hydrated: true, Dirty: true, EnvFilePresent: true, LastCheckedAt: nowRFC3339(),
			},
		},
	}
	if err := SaveState(state); err != nil {
		t.Fatalf("SaveState: %v", err)
	}
	beforeManifest, beforeState := readProjectListFiles(t, workspace)

	stdout, _, err := executeCommand(t, "test", "project")
	if err != nil {
		t.Fatalf("project list error: %v", err)
	}
	for _, want := range []string{"apps/api", "placeholder", "main", "docs/site", "missing", "services/worker", "hydrated", "yes"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("project list output missing %q:\n%s", want, stdout)
		}
	}
	afterManifest, afterState := readProjectListFiles(t, workspace)
	if beforeManifest != afterManifest {
		t.Fatal("project list mutated manifest")
	}
	if beforeState != afterState {
		t.Fatal("project list mutated state")
	}
}

func TestProjectCommandJSONHasStableFieldNames(t *testing.T) {
	workspace := initCommandWorkspace(t)
	manifest := Manifest{
		Version:       ManifestVersion,
		WorkspaceRoot: workspace,
		Projects: []Project{
			{ID: "project_api", Name: "api", Path: "apps/api", Type: ProjectTypeGit, Remote: "https://example.invalid/api.git", HydrateMode: HydrateOnDemand},
		},
	}
	if err := SaveManifest(workspace, manifest); err != nil {
		t.Fatalf("SaveManifest: %v", err)
	}
	state := State{Projects: map[string]ProjectState{
		"project_api": {Exists: true, Placeholder: true, CurrentBranch: "main", LastCheckedAt: nowRFC3339()},
	}}
	if err := SaveState(state); err != nil {
		t.Fatalf("SaveState: %v", err)
	}

	stdout, _, err := executeCommand(t, "test", "project", "--json")
	if err != nil {
		t.Fatalf("project --json error: %v", err)
	}
	var rows []struct {
		Project Project      `json:"project"`
		State   ProjectState `json:"state"`
	}
	if err := json.Unmarshal([]byte(stdout), &rows); err != nil {
		t.Fatalf("project --json did not parse: %v\n%s", err, stdout)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %+v, want one row", rows)
	}
	if rows[0].Project.Path != "apps/api" || !rows[0].State.Placeholder {
		t.Fatalf("row = %+v, want full project and matching state", rows[0])
	}
	for _, want := range []string{`"project"`, `"state"`, `"hydrateMode"`, `"placeholder"`} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("project --json missing field %q:\n%s", want, stdout)
		}
	}
	if strings.Contains(stdout, "\x1b[") {
		t.Fatalf("project --json contains ANSI escapes:\n%q", stdout)
	}
}

func TestProjectHelpDocumentsJSONFlag(t *testing.T) {
	stdout, _, err := executeCommand(t, "test", "project", "--help")
	if err != nil {
		t.Fatalf("project --help error: %v", err)
	}
	for _, want := range []string{"--json", "add", "hydrate", "remove", "update"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("project --help missing %q:\n%s", want, stdout)
		}
	}
}

func TestProjectUpdateCommandValidation(t *testing.T) {
	if _, _, err := executeCommand(t, "test", "project", "update"); err == nil || !strings.Contains(err.Error(), "<project> or --all") {
		t.Fatalf("project update error = %v, want missing target", err)
	}
	if _, _, err := executeCommand(t, "test", "project", "update", "--all", "app"); err == nil || !strings.Contains(err.Error(), "either --all or <project>") {
		t.Fatalf("project update --all app error = %v, want exclusive target", err)
	}
	stdout, _, err := executeCommand(t, "test", "project", "update", "--help")
	if err != nil {
		t.Fatalf("project update --help error: %v", err)
	}
	if !strings.Contains(stdout, "--all") {
		t.Fatalf("project update --help missing --all:\n%s", stdout)
	}
}

func readProjectListFiles(t *testing.T, workspace string) (string, string) {
	t.Helper()
	manifest, err := os.ReadFile(manifestPath(workspace))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	home, err := appHome()
	if err != nil {
		t.Fatalf("appHome: %v", err)
	}
	state, err := os.ReadFile(filepath.Join(home, "state.json"))
	if err != nil {
		t.Fatalf("read state: %v", err)
	}
	return string(manifest), string(state)
}

func TestPlanAndApplyCommandsRoundTrip(t *testing.T) {
	initCommandWorkspace(t)
	if _, _, err := executeCommand(t, "test", "project", "add", "services/api"); err != nil {
		t.Fatalf("project add error: %v", err)
	}

	planOut, _, err := executeCommand(t, "test", "plan")
	if err != nil {
		t.Fatalf("plan error: %v", err)
	}
	if !strings.Contains(planOut, "Planned changes:") {
		t.Fatalf("plan output = %q", planOut)
	}

	applyOut, _, err := executeCommand(t, "test", "apply")
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}
	if !strings.Contains(applyOut, "Applied safe plan actions.") {
		t.Fatalf("apply output = %q", applyOut)
	}
}

func TestSetupPlanCommandReportsNoSetupCommands(t *testing.T) {
	// An empty workspace has no detected setup commands; the plan prints the
	// "(none)" branch of printSetupPlan and exits cleanly.
	initCommandWorkspace(t)
	stdout, _, err := executeCommand(t, "test", "setup", "plan")
	if err != nil {
		t.Fatalf("setup plan error: %v", err)
	}
	if !strings.Contains(stdout, "Setup commands:") || !strings.Contains(stdout, "(none)") {
		t.Fatalf("setup plan output = %q", stdout)
	}
}

func TestStatusJSONHasStableFieldNames(t *testing.T) {
	initCommandWorkspace(t)
	stdout, _, err := executeCommand(t, "test", "status", "--json")
	if err != nil {
		t.Fatalf("status --json error: %v", err)
	}
	var report WorkspaceStatusReport
	if err := json.Unmarshal([]byte(stdout), &report); err != nil {
		t.Fatalf("status --json did not parse: %v\n%s", err, stdout)
	}
	if report.ProjectsTracked != 0 {
		t.Fatalf("projectsTracked = %d, want 0", report.ProjectsTracked)
	}
	if !strings.Contains(stdout, `"machine"`) || !strings.Contains(stdout, `"projectsTracked"`) {
		t.Fatalf("status --json missing expected field names:\n%s", stdout)
	}
}

func TestDoctorJSONHasStableFieldNames(t *testing.T) {
	initCommandWorkspace(t)
	stdout, _, err := executeCommand(t, "test", "doctor", "--json")
	if err != nil {
		t.Fatalf("doctor --json error: %v", err)
	}
	var report doctorReport
	if err := json.Unmarshal([]byte(stdout), &report); err != nil {
		t.Fatalf("doctor --json did not parse: %v\n%s", err, stdout)
	}
	if len(report.Checks) == 0 {
		t.Fatal("doctor --json reported zero checks")
	}
	if !strings.Contains(stdout, `"severity"`) || !strings.Contains(stdout, `"hardFailures"`) {
		t.Fatalf("doctor --json missing expected field names:\n%s", stdout)
	}
}

func TestWorkspaceDiffJSONHasStableFieldNames(t *testing.T) {
	initCommandWorkspace(t)
	stdout, _, err := executePrivateCommand(t, newWorkspaceCommand(), "diff", "--json")
	// No manifest remote is configured in this isolated workspace, so the
	// command is expected to fail; the point of this test is only that a
	// configured remote's diff would serialize with these exact field names,
	// verified against the zero-value struct's JSON shape.
	if err == nil {
		var diff ManifestDiff
		if unmarshalErr := json.Unmarshal([]byte(stdout), &diff); unmarshalErr != nil {
			t.Fatalf("workspace diff --json did not parse: %v\n%s", unmarshalErr, stdout)
		}
	}
	data, marshalErr := json.Marshal(ManifestDiff{})
	if marshalErr != nil {
		t.Fatalf("ManifestDiff failed to marshal: %v", marshalErr)
	}
	for _, want := range []string{`"added"`, `"removed"`, `"changed"`} {
		if !strings.Contains(string(data), want) {
			t.Fatalf("ManifestDiff JSON missing field %q:\n%s", want, data)
		}
	}
}

func TestMountPreviewJSONHasStableFieldNames(t *testing.T) {
	workspace := initCommandWorkspace(t)
	if err := SaveManifest(workspace, Manifest{
		Version:       ManifestVersion,
		WorkspaceRoot: workspace,
		Projects:      []Project{hardeningProject("apps/lazy", ProjectTypeGit, "https://example.invalid/lazy.git")},
	}); err != nil {
		t.Fatal(err)
	}
	stdout, _, err := executePrivateCommand(t, newMountCommand(), filepath.Join(t.TempDir(), "mnt"), "--preview", "--json")
	if err != nil {
		t.Fatalf("mount --preview --json error: %v", err)
	}
	var entries []MountEntry
	if err := json.Unmarshal([]byte(stdout), &entries); err != nil {
		t.Fatalf("mount --preview --json did not parse: %v\n%s", err, stdout)
	}
	if len(entries) != 1 || entries[0].Path != "apps/lazy" {
		t.Fatalf("entries = %+v, want one entry for apps/lazy", entries)
	}
	for _, want := range []string{`"hydrateMode"`, `"dirty"`, `"envPresent"`, `"setupHint"`} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("mount --preview --json missing expected field name %s:\n%s", want, stdout)
		}
	}
	if entries[0].Dirty || entries[0].EnvPresent || entries[0].SetupHint != "" {
		t.Fatalf("mount --preview --json diagnostics = %+v", entries[0])
	}
}

func TestProfileOrDefaultAndSortedProjectNames(t *testing.T) {
	// Directly exercise the two trivial helpers that the command scenarios
	// above don't reach.
	if profileOrDefault("") != "dev" {
		t.Fatal("profileOrDefault('') should be 'dev'")
	}
	if profileOrDefault("prod") != "prod" {
		t.Fatal("profileOrDefault('prod') should be 'prod'")
	}
	m := Manifest{Projects: []Project{{Name: "zeta"}, {Name: "alpha"}}}
	if got := sortedProjectNames(m); got != "alpha, zeta" {
		t.Fatalf("sortedProjectNames = %q, want %q", got, "alpha, zeta")
	}
}

// Silence prevents unused-import warnings if cobra is only referenced in
// helpers above; keep the import explicit.
var _ = cobra.Command{}
