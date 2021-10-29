# semversioner
This package determines the next semantic version to any git repository based on the `type` argument (`major`,`minor`, `patch`) 
and the most recent applied tag. It can also apply [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) 
versioning by using the `--type=conventional` argument.

## Build
```shell
go build -o semversioner
```

## Usage
Given a git repository with the most recent tag as `v0.1.8`

### Default snapshot tag
```shell
$ ~/code/my-app on chore/patch-envoy-logs ◦ ./versioner -
0.1.8-4-fb067b1-SNAPSHOT
```

### Prerelease tags
```shell
$ ~/code/my-app on feature-1 ◦ ./versioner --prerelease rc1
0.1.8-rc1
```

### Semver overrides
```shell
$ ~/code/my-app on feature-1 ◦ ./versioner --type major
1.0.0
```

### Conventional commit versioning
Given the last couple of log commit messages following the conventional commit guidelines:
```shell
commit fb067b14f2e9d24fee367651603424a8304a0845 (HEAD -> testGit)
Author: John doe <john@mailinator.com>
Date:   Tue Nov 16 21:38:29 2021 -0800

    feat(scope)!: this is a test description that breaks

    This this just testing a breaking change that should be a major.

    here is an article [lalal](google.com)

commit c1003815c1cfe3c0cd718550573031a8bf536188
Author: Jane Doe <jane@mailinator.com>
Date:   Sun Oct 17 15:08:37 2021 -0700

    feat(scope): this is a new feature

commit 574a7e20eeaac9905bbf0fd3486ab13a7a0368c3
Author: Bobby Doe <bobby@mailinator.com>
Date:   Wed Sep 29 16:54:49 2021 -0700

    This is not a conventional commit
```

```shell
$ ~/code/my-app on feature-1 ◦ ./versioner --type conventional
1.0.0
```

## Test
```shell
 go test ./... -test.v
```