# githoard

Quickly backup collections of git repositories with the help of [gitea]!

## Why?

Anything that can go wrong, will go wrong. Protect valuable repositories (including other's 
repositories) from accidental deletion by mirroring them!

## Quickstart

To quickly get you up and running, follow these easy steps:

### Login

githoard needs to know where to mirror archives and how to authenticate with [gitea] first. You will
need a personal access token to get started:

```
$ githoard login https://gitea.example.org abcdef01234

authenticated with https://gitea.example.org as dotStart#1
```

If you wish to mirror private GitHub repositories, you may also specify a personal access token
for GitHub:

```
$ githoard login -github-token=abcdef01234 https://gitea.example.org abcdef01234
```

Don't worry if you missed this when first setting up. You can invoke `githoard login` again at any
point in time!

### Backup!

`githoard` provides you with two primary commands to create mirrors:

1. `githoard repo` which creates a mirror for a single repository
2. `githoard profile` which creates mirrors for all repositories of a given user

Simply pass the URL of the desired repository or profile as seen in your browser:

```
$ githoard repo https://github.com/dotStart/Beacon

created repository dotStart/Beacon#1 - https://gitea.example.org/dotStart/Beacon
```

or alternatively:

```
$ githoard profile https://github.com/dotStart

created repository dotStart/Beacon#1 - https://gitea.example.org/dotStart/Beacon
created repository dotStart/overlord#2 - https://gitea.example.org/dotStart/overlord
created repository dotStart/githoard#3 - https://gitea.example.org/dotStart/githoard
```

That's it!

## License

```
Copyright [yyyy] [name of copyright owner]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

[gitea]: https://gitea.io/en-us/
