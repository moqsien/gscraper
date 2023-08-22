## gscraper 是什么?

gscraper 是用来自动在github下载[gvc](https://github.com/moqsien/gvc)所需要的软件安装包。
然后将这些软件包上传到[gitlab](https://gitlab.com/moqsien/gvc_resources)进行缓存。
默认下载软件当前的最新版本。

主要是为了解决中国大陆网络长城对GitHub的阻断，导致正常的下载非常慢或者根本无法进行。

## gscraper当前自动下载的软件列表

- [protobuf](https://github.com/protocolbuffers/protobuf)
- [vlang](https://github.com/vlang/v)
- [v-analyzer](https://github.com/v-analyzer/v-analyzer)
- [typst](https://github.com/typst/typst)
- [neovim](https://github.com/neovim/neovim)
- [vcpkg](https://github.com/microsoft/vcpkg)
- [vcpkg-tool](https://github.com/microsoft/vcpkg-tool)
- [pyenv](https://github.com/pyenv/pyenv)
- [pyenv-win](https://github.com/pyenv-win/pyenv-win)
- [sing-box-rules](https://github.com/lyc8503/sing-box-rules)
- [v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat)
- [gvc](https://github.com/moqsien/gvc)

## gscraper 命令

```bash
>>> gscraper help

NAME:
   gscraper.exe - gscraper <Command> <SubCommand>...

USAGE:
   gscraper.exe [global options] command [command options] [arguments...]

DESCRIPTION:
   gscraper, download files from github for gvc.

COMMANDS:
   show, sh           Show url list to download.
   add, a             Add new download urls[ gscraper add xxx yyy zzz ss].
   remove, rm, r      Remove download url.[ gscraper rm xxx ]
   reset, rs          Reset config file.
   download, down, d  Download files.
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

- 可以通过add命令添加下载项目
- 可以通过download命令下载所有项目或者下载指定项目
- 配置文件存放在当前用户的家目录(os.UserHomeDir()获取)下，文件名为gscraper_conf.json。

添加下载项目url举例：
```text
https://github.com/vlang/v/releases/latest/download/v_linux.zip
https://github.com/neovim/neovim/releases/download/stable/nvim-linux64.tar.gz
https://github.com/typst/typst/releases/latest/download/typst-aarch64-unknown-linux-musl.tar.xz
https://github.com/microsoft/vcpkg-tool/archive/refs/heads/main.zip
https://github.com/pyenv-win/pyenv-win/archive/refs/heads/master.zip
```

## 如何安装？
```bash
go install github.com/moqsien/gscraper@latest
```
