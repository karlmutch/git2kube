## git2kube watch folder

Runs watcher that periodically check the provided repository and updates target folder accordingly

### Synopsis

Runs watcher that periodically check the provided repository and updates target folder accordingly

```
git2kube watch folder [flags]
```

### Options

```
  -h, --help                   help for folder
  -t, --target-folder string   path to target folder
```

### Options inherited from parent commands

```
  -b, --branch string         branch name to pull (default "master")
  -c, --cache-folder string   destination on filesystem where cache of repository will be stored (default "/tmp/git2kube/data/")
      --exclude strings       regex that if is a match excludes the file from the upload, example: '*.yaml' or 'folder/*' if you want to match a folder (default [^\..*])
  -g, --git string            git repository address, either http(s) or ssh protocol has to be specified
      --include strings       regex that if is a match includes the file in the upload, example: '*.yaml' or 'folder/*' if you want to match a folder (default [.*])
  -i, --interval int          interval in seconds in which to try refreshing ConfigMap from git (default 10)
  -l, --log-level string      command log level (options: [panic fatal error warning info debug]) (default "info")
```

### SEE ALSO

* [git2kube watch](git2kube_watch.md)	 - Runs watcher that periodically check the provided repository
