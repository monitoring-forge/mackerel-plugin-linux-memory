# mackerel-plugin-linux-memory

[mackerel](https://mackerel.io) 用のカスタムメトリクスプラグインです。Linux の `/proc/meminfo` からメモリ情報を収集し、Mackerel に投稿します。

## Features

- `/proc/meminfo` の各種メモリ統計をバイト単位で収集
- Linux カーネル 3.14 以降の `MemAvailable` に対応
- `MemAvailable` がない古いカーネルでは `MemTotal - MemFree - Buffers - Cached` で `used` を推定
- 未取得のフィールドは `0` を返し、プラグインのクラッシュを回避

## Requirements

- Linux カーネル（`/proc/meminfo` が必要）
- Go 1.25 以上（ビルドする場合）

## Installation

### バイナリをダウンロードする

GitHub Releases から最新版のバイナリをダウンロードしてください。

### mkr コマンドでインストールする

```bash
mkr plugin install monitoring-forge/mackerel-plugin-linux-memory
```

### ソースからビルドする

```bash
git clone https://github.com/monitoring-forge/mackerel-plugin-linux-memory.git
cd mackerel-plugin-linux-memory
go build -o mackerel-plugin-linux-memory
```

Linux 向けにクロスコンパイルする場合:

```bash
make linux
```

## Usage

```bash
$ ./mackerel-plugin-linux-memory
linux-memory.total      5198495744      1627264100
linux-memory.available  2149441536      1627264100
linux-memory.used       3049054208      1627264100
linux-memory.kernelstack        20594688        1627264100
linux-memory.vmallocused        23023616        1627264100
linux-memory.pagetables 25468928        1627264100
linux-memory.mapped     92598272        1627264100
linux-memory.anonpages  1818071040      1627264100
linux-memory.slab       405975040       1627264100
linux-memory.buffers    173051904       1627264100
linux-memory.cached     1653784576      1627264100
linux-memory.free       989589504       1627264100
```

### オプション

```bash
$ ./mackerel-plugin-linux-memory --version
mackerel-plugin-linux-memory 0.0.7
Compiler: gc go1.25.0
```

## Metrics

| メトリクス名            | 説明                                              | ソース (`/proc/meminfo`) |
| ----------------------- | ------------------------------------------------- | ------------------------ |
| `linux-memory.total`    | 搭載メモリ総量                                    | `MemTotal`               |
| `linux-memory.available`| 利用可能なメモリ量                                | `MemAvailable`           |
| `linux-memory.used`     | 使用中メモリ量                                    | `MemTotal - MemAvailable`（`MemAvailable` 非対応時は `MemTotal - MemFree - Buffers - Cached`） |
| `linux-memory.free`     | 完全に未使用のメモリ量                            | `MemFree`                |
| `linux-memory.buffers`  | カーネルバッファキャッシュ                        | `Buffers`                |
| `linux-memory.cached`   | ページキャッシュ                                  | `Cached`                 |
| `linux-memory.slab`     | スラブで使用中のメモリ                            | `Slab`                   |
| `linux-memory.anonpages`| 匿名ページで使用中のメモリ                        | `AnonPages`              |
| `linux-memory.mapped`   | mmap されたメモリ                                 | `Mapped`                 |
| `linux-memory.pagetables` | ページテーブルで使用中のメモリ                  | `PageTables`             |
| `linux-memory.kernelstack` | カーネルスタックで使用中のメモリ               | `KernelStack`            |
| `linux-memory.vmallocused` | `vmalloc` で使用中のメモリ                     | `VmallocUsed`            |

すべてのメトリクスはバイト単位で出力されます。


## License

MIT