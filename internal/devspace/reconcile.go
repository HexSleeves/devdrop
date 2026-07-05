package devspace

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
)

type ReconcileOp struct {
	Action string `json:"action"`
	Kind   string `json:"kind"`
	Key    string `json:"key"`
}

type ReconcileResult struct {
	Merged    Manifest        `json:"merged"`
	Ops       []ReconcileOp   `json:"ops"`
	Conflicts []MergeConflict `json:"conflicts"`
	TwoWay    bool            `json:"twoWay"`
}

type ReconcileRemoteSource struct {
	GitRemote string `json:"gitRemote,omitempty"`
	Remote    string `json:"remote,omitempty"`
	Commit    string `json:"commit,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	Version   int    `json:"version,omitempty"`
}

type ReconcilePlan struct {
	Version       int                   `json:"version"`
	CreatedAt     string                `json:"createdAt"`
	Backend       string                `json:"backend"`
	ManifestHash  string                `json:"manifestHash"`
	RemoteSource  ReconcileRemoteSource `json:"remoteSource"`
	Ops           []ReconcileOp         `json:"ops"`
	Conflicts     []MergeConflict       `json:"conflicts"`
	TwoWay        bool                  `json:"twoWay"`
	Merged        Manifest              `json:"merged"`
	WorkspaceRoot string                `json:"workspaceRoot"`
}

func reconcileManifests(base *Manifest, local, remote Manifest) (ReconcileResult, error) {
	var result ReconcileResult
	var err error
	if base != nil {
		result.Merged, result.Conflicts, err = mergeManifests(*base, local, remote)
		if err != nil {
			return ReconcileResult{}, err
		}
	} else {
		result, err = reconcileTwoWay(local, remote)
		if err != nil {
			return ReconcileResult{}, err
		}
	}
	result.Ops = diffReconcileOps(local, result.Merged)
	return result, nil
}

func reconcileTwoWay(local, remote Manifest) (ReconcileResult, error) {
	if err := ValidateManifest(local); err != nil {
		return ReconcileResult{}, fmt.Errorf("local manifest failed validation: %w", err)
	}
	if err := ValidateManifest(remote); err != nil {
		return ReconcileResult{}, fmt.Errorf("remote manifest failed validation: %w", err)
	}
	merged := local
	var conflicts []MergeConflict

	projects, projectConflicts := twoWayProjectRecords(local.Projects, remote.Projects)
	conflicts = append(conflicts, projectConflicts...)
	merged.Projects = projects

	access, accessConflicts := twoWayAccessRecords(local.Access, remote.Access)
	conflicts = append(conflicts, accessConflicts...)
	merged.Access = access

	if len(conflicts) == 0 {
		if err := ValidateManifest(merged); err != nil {
			return ReconcileResult{}, fmt.Errorf("merged manifest failed validation: %w", err)
		}
	}
	return ReconcileResult{Merged: merged, Conflicts: conflicts, TwoWay: true}, nil
}

func twoWayProjectRecords(local, remote []Project) ([]Project, []MergeConflict) {
	localByID := projectByID(local)
	remoteByID := projectByID(remote)
	keys := map[string]bool{}
	for key := range localByID {
		keys[key] = true
	}
	for key := range remoteByID {
		keys[key] = true
	}
	merged := make([]Project, 0, len(keys))
	var conflicts []MergeConflict
	for _, key := range sortedKeys(keys) {
		localProject, inLocal := localByID[key]
		remoteProject, inRemote := remoteByID[key]
		switch {
		case inLocal && inRemote:
			merged = append(merged, localProject)
			if !reflect.DeepEqual(localProject, remoteProject) {
				conflicts = append(conflicts, MergeConflict{Entity: "project", Key: key, Field: "*", Ours: fmt.Sprintf("%+v", localProject), Theirs: fmt.Sprintf("%+v", remoteProject)})
			}
		case inLocal:
			merged = append(merged, localProject)
		case inRemote:
			merged = append(merged, remoteProject)
		}
	}
	sortProjects(merged)
	return merged, conflicts
}

func twoWayAccessRecords(local, remote []ProjectAccess) ([]ProjectAccess, []MergeConflict) {
	localByKey := accessByKey(local)
	remoteByKey := accessByKey(remote)
	keys := map[string]bool{}
	for key := range localByKey {
		keys[key] = true
	}
	for key := range remoteByKey {
		keys[key] = true
	}
	merged := make([]ProjectAccess, 0, len(keys))
	var conflicts []MergeConflict
	for _, key := range sortedKeys(keys) {
		localAccess, inLocal := localByKey[key]
		remoteAccess, inRemote := remoteByKey[key]
		switch {
		case inLocal && inRemote:
			merged = append(merged, localAccess)
			if !reflect.DeepEqual(localAccess, remoteAccess) {
				conflicts = append(conflicts, MergeConflict{Entity: "access", Key: key, Field: "*", Ours: fmt.Sprintf("%+v", localAccess), Theirs: fmt.Sprintf("%+v", remoteAccess)})
			}
		case inLocal:
			merged = append(merged, localAccess)
		case inRemote:
			merged = append(merged, remoteAccess)
		}
	}
	sortAccess(merged)
	return merged, conflicts
}

func diffReconcileOps(local, merged Manifest) []ReconcileOp {
	var ops []ReconcileOp
	localProjects := projectByID(local.Projects)
	mergedProjects := projectByID(merged.Projects)
	projectKeys := map[string]bool{}
	for key := range localProjects {
		projectKeys[key] = true
	}
	for key := range mergedProjects {
		projectKeys[key] = true
	}
	for _, key := range sortedKeys(projectKeys) {
		localProject, inLocal := localProjects[key]
		mergedProject, inMerged := mergedProjects[key]
		switch {
		case !inLocal && inMerged:
			ops = append(ops, ReconcileOp{Action: "added", Kind: "project", Key: key})
		case inLocal && !inMerged:
			ops = append(ops, ReconcileOp{Action: "removed", Kind: "project", Key: key})
		case inLocal && inMerged && !reflect.DeepEqual(localProject, mergedProject):
			ops = append(ops, ReconcileOp{Action: "changed", Kind: "project", Key: key})
		}
	}

	localAccess := accessByKey(local.Access)
	mergedAccess := accessByKey(merged.Access)
	accessKeys := map[string]bool{}
	for key := range localAccess {
		accessKeys[key] = true
	}
	for key := range mergedAccess {
		accessKeys[key] = true
	}
	for _, key := range sortedKeys(accessKeys) {
		localGrant, inLocal := localAccess[key]
		mergedGrant, inMerged := mergedAccess[key]
		switch {
		case !inLocal && inMerged:
			ops = append(ops, ReconcileOp{Action: "added", Kind: "access", Key: key})
		case inLocal && !inMerged:
			ops = append(ops, ReconcileOp{Action: "removed", Kind: "access", Key: key})
		case inLocal && inMerged && !reflect.DeepEqual(localGrant, mergedGrant):
			ops = append(ops, ReconcileOp{Action: "changed", Kind: "access", Key: key})
		}
	}
	return ops
}

func ReconcileWorkspaceManifest(force string, apply bool) (ReconcilePlan, error) {
	if force != "" && force != "local" && force != "remote" {
		return ReconcilePlan{}, fmt.Errorf("force must be one of: local, remote")
	}
	cfg, err := syncConfig()
	if err != nil {
		return ReconcilePlan{}, err
	}
	local, err := LoadManifest(cfg.WorkspaceRoot)
	if err != nil {
		return ReconcilePlan{}, err
	}
	hash, err := ManifestHash(cfg.WorkspaceRoot)
	if err != nil {
		return ReconcilePlan{}, err
	}
	remote, err := fetchLocalizedWorkspaceRemoteManifest(cfg)
	if err != nil {
		return ReconcilePlan{}, err
	}
	source := reconcileGitRemoteSource(cfg)
	if previous, ok := reusableReconcilePlan(local, source); ok {
		plan := ReconcilePlan{
			Version:       1,
			CreatedAt:     nowRFC3339(),
			Backend:       "git",
			ManifestHash:  hash,
			RemoteSource:  source,
			TwoWay:        previous.TwoWay,
			Merged:        local,
			WorkspaceRoot: cfg.WorkspaceRoot,
		}
		if err := SaveReconcilePlan(plan); err != nil {
			return ReconcilePlan{}, err
		}
		return plan, nil
	}
	baseManifest, hasBase, err := loadBaseManifest()
	if err != nil {
		return ReconcilePlan{}, err
	}
	var base *Manifest
	if hasBase {
		base = &baseManifest
	}
	result, err := reconcileManifests(base, local, remote)
	if err != nil {
		return ReconcilePlan{}, err
	}
	if force != "" && len(result.Conflicts) > 0 {
		result.Merged = forceReconcileConflicts(result.Merged, result.Conflicts, local, remote, force)
		result.Conflicts = nil
		if err := ValidateManifest(result.Merged); err != nil {
			return ReconcilePlan{}, fmt.Errorf("forced merged manifest failed validation: %w", err)
		}
		result.Ops = diffReconcileOps(local, result.Merged)
	}
	plan := ReconcilePlan{
		Version:       1,
		CreatedAt:     nowRFC3339(),
		Backend:       "git",
		ManifestHash:  hash,
		RemoteSource:  source,
		Ops:           result.Ops,
		Conflicts:     result.Conflicts,
		TwoWay:        result.TwoWay,
		Merged:        result.Merged,
		WorkspaceRoot: cfg.WorkspaceRoot,
	}
	if err := SaveReconcilePlan(plan); err != nil {
		return ReconcilePlan{}, err
	}
	if !apply {
		return plan, nil
	}
	if len(plan.Conflicts) > 0 {
		return plan, fmt.Errorf("unresolved reconcile conflicts:\n%s", formatReconcileConflictErrors(plan.Conflicts))
	}
	current, err := LoadManifest(cfg.WorkspaceRoot)
	if err != nil {
		return ReconcilePlan{}, err
	}
	currentHash, err := ManifestHash(cfg.WorkspaceRoot)
	if err != nil {
		return ReconcilePlan{}, err
	}
	if currentHash != plan.ManifestHash {
		return ReconcilePlan{}, fmt.Errorf("manifest changed since reconcile was generated; run `devspace workspace reconcile` again before apply")
	}
	backup, err := manifestBackupPath()
	if err != nil {
		return ReconcilePlan{}, err
	}
	if err := writeJSON(backup, current, 0o600); err != nil {
		return ReconcilePlan{}, err
	}
	if err := SaveManifest(cfg.WorkspaceRoot, plan.Merged); err != nil {
		return ReconcilePlan{}, err
	}
	if err := recordBaseManifest(plan.Merged); err != nil {
		return ReconcilePlan{}, err
	}
	return plan, nil
}

func reusableReconcilePlan(local Manifest, source ReconcileRemoteSource) (ReconcilePlan, bool) {
	previous, err := LoadReconcilePlan()
	if err != nil {
		return ReconcilePlan{}, false
	}
	if previous.Version == 0 || previous.Backend != "git" || len(previous.Conflicts) > 0 {
		return ReconcilePlan{}, false
	}
	if previous.RemoteSource.Commit == "" || previous.RemoteSource.Commit != source.Commit {
		return ReconcilePlan{}, false
	}
	if !reflect.DeepEqual(local, previous.Merged) {
		return ReconcilePlan{}, false
	}
	return previous, true
}

func reconcileGitRemoteSource(cfg Config) ReconcileRemoteSource {
	source := ReconcileRemoteSource{
		GitRemote: "origin",
		Remote:    redactRemote(cfg.ManifestRemote),
	}
	repo, err := expandPath(cfg.ManifestRepoPath)
	if err != nil {
		return source
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	source.Commit = mustGit(ctx, repo, "rev-parse", "HEAD")
	return source
}

func forceReconcileConflicts(merged Manifest, conflicts []MergeConflict, local, remote Manifest, force string) Manifest {
	projects := projectByID(merged.Projects)
	access := accessByKey(merged.Access)
	localProjects := projectByID(local.Projects)
	remoteProjects := projectByID(remote.Projects)
	localAccess := accessByKey(local.Access)
	remoteAccess := accessByKey(remote.Access)
	for _, conflict := range conflicts {
		switch conflict.Entity {
		case "project":
			source := localProjects
			if force == "remote" {
				source = remoteProjects
			}
			if project, ok := source[conflict.Key]; ok {
				projects[conflict.Key] = project
			} else {
				delete(projects, conflict.Key)
			}
		case "access":
			source := localAccess
			if force == "remote" {
				source = remoteAccess
			}
			if grant, ok := source[conflict.Key]; ok {
				access[conflict.Key] = grant
			} else {
				delete(access, conflict.Key)
			}
		}
	}
	merged.Projects = make([]Project, 0, len(projects))
	for _, project := range projects {
		merged.Projects = append(merged.Projects, project)
	}
	sortProjects(merged.Projects)
	merged.Access = make([]ProjectAccess, 0, len(access))
	for _, grant := range access {
		merged.Access = append(merged.Access, grant)
	}
	sortAccess(merged.Access)
	return merged
}

func formatReconcileConflictErrors(conflicts []MergeConflict) string {
	lines := make([]string, 0, len(conflicts))
	for _, conflict := range conflicts {
		lines = append(lines, fmt.Sprintf("- %s %s %s: local=%q remote=%q", conflict.Entity, conflict.Key, conflict.Field, conflict.Ours, conflict.Theirs))
	}
	return strings.Join(lines, "\n")
}

func SaveReconcilePlan(plan ReconcilePlan) error {
	path, err := reconcilePlanPath()
	if err != nil {
		return err
	}
	return writeJSON(path, plan, 0o600)
}

func LoadReconcilePlan() (ReconcilePlan, error) {
	var plan ReconcilePlan
	path, err := reconcilePlanPath()
	if err != nil {
		return plan, err
	}
	if err := readJSON(path, &plan); err != nil {
		return plan, err
	}
	return plan, nil
}

func sortProjects(projects []Project) {
	slices.SortFunc(projects, func(a, b Project) int {
		return strings.Compare(a.Path, b.Path)
	})
}

func sortAccess(access []ProjectAccess) {
	slices.SortFunc(access, func(a, b ProjectAccess) int {
		return strings.Compare(accessKey(a), accessKey(b))
	})
}
