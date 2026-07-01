# FUSE Lazy Workspace Mount Prototype

`devspace mount <mountpoint>` is a prototype for a read-only workspace view backed
by the DevDrop manifest. It is intentionally outside normal sync, plan, apply,
and hydrate workflows so the CLI still works on machines without FUSE.

## Library Selection

The prototype uses `github.com/hanwen/go-fuse/v2/fs`.

- `go-fuse/v2` is the best fit for DevDrop because it is Go-native, current, and
  ships a higher-level `fs` node API plus loopback filesystem support.
- `bazil.org/fuse` is also Go-native, but its published package is older and the
  API would require more custom filesystem plumbing for this spike.
- `github.com/jacobsa/fuse` has useful examples, but it is less aligned with a
  manifest tree that can switch from virtual entries to loopback directories.

## Behavior

The mount exposes tracked project paths from `.devdrop/manifest.json`.

- `ls <mountpoint>` lists top-level manifest path segments without hydrating
  projects.
- Traversing into an on-demand Git project runs the same safety path as
  `devspace project hydrate <project>`.
- Hydration failures are returned to the filesystem caller and logged to stderr;
  the mount does not convert them into empty successful directories.
- Local-only, manual, metadata-only, or missing projects that cannot hydrate
  automatically are represented by a stub directory containing `.devdrop-status`.
- The mountpoint must be empty. DevDrop refuses to mount over non-empty
  directories so local files are not hidden.

## Platform Requirements

FUSE support is optional and platform-specific.

- macOS requires macFUSE or a compatible FUSE implementation installed and
  approved by the OS.
- Linux requires `/dev/fuse` and permission to mount FUSE filesystems, often via
  `fusermount3` or equivalent distribution packaging.
- CI and normal CLI workflows do not require FUSE.

If FUSE is unavailable, use:

```bash
devspace mount /tmp/devspace-mount --preview
```

`--preview` prints the manifest-backed entries and hydration status without
mounting anything.

## Follow-Up Cards

- Add an integration test job that runs only on FUSE-capable hosts and exercises
  actual mount traversal, hydration success, and hydration failure propagation.
- Add a richer project status view in the mount for dirty repositories, missing
  `.env` files, and setup hints.
- Decide whether the long-term mount should expose project paths, project names,
  or both through separate virtual directories.
- Add unmount diagnostics and stale mount cleanup guidance after real-user
  testing on macOS and Linux.
