# ca-go
A Go library with multiple packages to be shared by services.

## Decisions
* This library is designed to be a replacement of the old Go library: `cultureamp/glamplify`. 
  * Creating a glamplify V2 is not easier than creating a new repo. 
  * But having a new repo can allow us using both libraries together for a while so we don’t have to do the migration at the beginning.
* It's created as a public repository.
  * We won't put packages that can't be public in at this stage. Although we might do that in the future.
  * Make this repo **internal** is not that easy due to the way we config **BuildKite**.
  * Too much effort for the library users update their config and get the private packages. 
* This repo will only have one go.mod so we can't version control them separately.
  * There is no package registry like npm for Go. Can't really do version control separately.   

## Migrate approach:
* Create `cultureamp/ca-go` as a public repo
* Start to move some packages from the shims repo here. e.g. `/ref`. Keep using `glamplify` for existing packages.
* Start to move packages from `glamplify` to here. Stop using `glamplify` when all packages are migrated.
* Make this repo internal.

## Contributing

To work on `ca-go`, you'll need a working Go installation. The project currently
targets Go 1.17.

### Setting up your environment

You can use [VSCode Remote
Containers](https://code.visualstudio.com/docs/remote/containers) to get
up-and-running quickly. A basic configuration is defined in the `.devcontainer/`
directory. This works locally and via [GitHub
Codespaces](https://github.com/features/codespaces).

**Locally**:

1. Clone `ca-go` and open the directory in VSCode.
2. A prompt should appear on the bottom-right of the editor, offering to start a
   Remote Containers session. Click **Reopen in Container**.
3. If a prompt didn't appear, open the Command Palette (i.e. Cmd + Shift + P)
   and select **Remote-Containers: Open Folder in Container...**

**Codespaces**:

1. Click the **Code** button above the source listing on the repository
   homepage.
2. Click **New codespace**.